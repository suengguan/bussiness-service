package service

import (
	"os"
	"strconv"
	"time"

	"model"
	"utility/fileoperator"
	"utility/httpclient"

	daoApi "api/dao_service"
	systemApi "api/system_service"
	"system-service/kubernetes/KubeRESTfulClient/ResourceModel" // todo

	"encoding/json"

	"github.com/astaxie/beego"
)

type JobService struct {
}

func (this *JobService) Create(job *model.Job) (*model.Job, error) {
	var err error
	var newJob *model.Job

	// get project
	beego.Debug("->get project")
	var project *model.Project
	project, err = daoApi.BussinessDaoApi.GetProjectById(job.Project.Id)
	if err != nil {
		return nil, err
	}

	// create job
	beego.Debug("->create job")
	job.Id = 0
	job.Project = project
	newJob, err = daoApi.BussinessDaoApi.CreateJob(job)
	if err != nil {
		return nil, err
	}

	// get user
	beego.Debug("->get user")
	var user *model.User
	user, err = daoApi.UserDaoApi.GetById(project.User.Id)
	if err != nil {
		return nil, err
	}

	// run job
	beego.Debug("->run job")
	ch := make(chan string, 1)
	go this.runJob(user, newJob, ch)
	go func() {
		v, _ := <-ch
		beego.Debug("value :", v)
		close(ch)
	}()

	beego.Debug("result:", *newJob)

	return newJob, err
}

func (this *JobService) Update(job *model.Job) (*model.Job, error) {
	var err error
	var newJob *model.Job

	// check job status
	beego.Debug("->check job status")

	// get project
	beego.Debug("->get project")
	var project *model.Project
	project, err = daoApi.BussinessDaoApi.GetProjectById(job.Project.Id)
	if err != nil {
		return nil, err
	}

	// update job
	beego.Debug("->update job")
	newJob, err = daoApi.BussinessDaoApi.UpdateJob(job)
	newJob.Project = project
	if err != nil {
		return nil, err
	}

	// get user
	beego.Debug("->get user")
	var user *model.User
	user, err = daoApi.UserDaoApi.GetById(project.User.Id)
	if err != nil {
		return nil, err
	}

	// run job
	beego.Debug("->run job")
	ch := make(chan string, 1)
	go this.runJob(user, newJob, ch)
	go func() {
		v, _ := <-ch
		beego.Debug("value :", v)
		close(ch)
	}()

	beego.Debug("result:", *newJob)

	return newJob, err
}

func (this *JobService) createWorkspace(user *model.User, job *model.Job) error {
	var err error

	for _, module := range job.Modules {
		for _, pod := range module.Pods {
			var podPath string

			var cfg = beego.AppConfig
			podPath = cfg.String("workspace")
			// pod path
			podPath += "/" + user.Name
			podPath += "/" + job.Project.Name + "-" + strconv.FormatInt(job.Project.Id, 36)
			podPath += "/" + job.Name + "-" + strconv.FormatInt(job.Id, 36)
			podPath += "/" + module.Name + "-" + strconv.FormatInt(module.Id, 36)

			beego.Debug("path:", podPath)
			err = os.MkdirAll(podPath, os.ModePerm) //生成多级目录
			if err != nil {
				return err
			}

			var logfn string
			logfn = podPath + "/" + pod.Name + "-" + strconv.FormatInt(pod.Id, 36) + ".log"

			beego.Debug("log file:", logfn)
			var statusfn string
			statusfn = podPath + "/" + pod.Name + "-" + strconv.FormatInt(pod.Id, 36) + ".st"

			fileoperator.Write(logfn, "")
			var status model.JobStatus
			var statusBody []byte
			statusBody, err = json.Marshal(&status)
			if err != nil {
				beego.Debug("json Unmarshal data failed")
			}
			fileoperator.Write(statusfn, string(statusBody))
			beego.Debug("status file:", statusfn)
		}
	}

	return err
}

func (this *JobService) runJob(user *model.User, job *model.Job, ch chan string) {
	beego.Debug("==run job:", job.Name)
	//var err error

	// 1.standalone
	moduleCnt := len(job.Modules)
	chs := make(chan string, moduleCnt)

	for _, m := range job.Modules {
		var svc JobService
		go svc.runModule(user, m, job, chs)
	}

	// waiting for finish
	for i := 0; i < moduleCnt; i++ {
		v, _ := <-chs
		beego.Debug("value :", v)
	}
	close(chs)

	// 2.flow

	// 3. mix

	ch <- "finish job :" + job.Name
}

func (this *JobService) runModule(user *model.User, module *model.Module, job *model.Job, ch chan string) {
	// run a module
	beego.Debug("==run module:", module.Name)

	podCnt := len(module.Pods)

	// channel buffer
	chs := make(chan string, podCnt)

	for _, p := range module.Pods {
		var svc JobService
		go svc.runPod(user, p, module, job, chs)
	}

	for i := 0; i < podCnt; i++ {
		v, _ := <-chs
		beego.Debug("value :", v)
	}

	close(chs)

	ch <- "finish module :" + module.Name
}

func (this *JobService) runPod(user *model.User, pod *model.Pod, module *model.Module, job *model.Job, ch chan string) {
	// 启动Processor
	beego.Debug("run pod :", pod.Name+"-"+module.Name+"-"+job.Name)

	var err error

	// create workspace
	beego.Debug("->create workspace")
	err = this.createWorkspace(user, job)
	if err != nil {
		beego.Debug("->create workspace failed!", err)
		return
	}

	// create rc
	//beego.Debug("== start to create pod")
	namespace := user.Name
	label := pod.RcName
	name := pod.SvcName
	// todo config
	//image := config.ALGORITHM_ENV_IMAGE
	containerPort := 8080
	image := "192.168.0.21:5000/system/algorithm-env:v4"
	limitCpu := strconv.FormatFloat(pod.Cpu, 'f', 3, 32)
	limitMemory := strconv.FormatFloat(pod.Memory*1024.0, 'f', 0, 32) + "Mi"
	envDescriptionStr := module.Description

	// set env
	beego.Debug("== set pod env")
	var envList []*ResourceModel.Env

	var envDescription ResourceModel.Env
	envDescription.Name = "DESCRIPTION"
	envDescription.Value = envDescriptionStr
	envList = append(envList, &envDescription)
	beego.Debug("->  " + envDescription.Name + " = " + envDescription.Value)

	var envLogFileName ResourceModel.Env
	envLogFileName.Name = "LOG_FILE_NAME"
	envLogFileName.Value = pod.Name + "-" + strconv.FormatInt(pod.Id, 36) + ".log"
	envList = append(envList, &envLogFileName)
	beego.Debug("->  " + envLogFileName.Name + " = " + envLogFileName.Value)

	var envStatusFileName ResourceModel.Env
	envStatusFileName.Name = "STATUS_FILE_NAME"
	envStatusFileName.Value = pod.Name + "-" + strconv.FormatInt(pod.Id, 36) + ".st"
	envList = append(envList, &envStatusFileName)
	beego.Debug("->  " + envStatusFileName.Name + " = " + envStatusFileName.Value)

	// set host path
	var workspacePath = "/pme2017/workspace"

	beego.Debug("== set host path")
	var hostPathList []*ResourceModel.ContainerHostVolume
	var h1 ResourceModel.ContainerHostVolume
	h1.Name = "input"
	h1.ContainerPath = "/input"
	h1.HostPath = workspacePath + "/" + user.Name + "/data/input"
	h1.ReadOnly = false
	hostPathList = append(hostPathList, &h1)
	beego.Debug("->  " + h1.Name + " = " + h1.HostPath + ":" + h1.ContainerPath)

	var h2 ResourceModel.ContainerHostVolume
	h2.Name = "output"
	h2.ContainerPath = "/output"
	h2.HostPath = workspacePath + "/" + user.Name + "/data/output"
	h2.ReadOnly = false
	hostPathList = append(hostPathList, &h2)
	beego.Debug("->  " + h2.Name + " = " + h2.HostPath + ":" + h2.ContainerPath)

	var h3 ResourceModel.ContainerHostVolume
	h3.Name = "algorithm"
	h3.ContainerPath = "/algorithm"
	// todo
	//h3.HostPath = config.ALGORITHM_PATH + "/" + module.Algorithm
	h3.HostPath = "/pme2017/algorithm" + "/" + module.Algorithm
	h3.ReadOnly = true
	hostPathList = append(hostPathList, &h3)
	beego.Debug("->  " + h3.Name + " = " + h3.HostPath + ":" + h3.ContainerPath)

	var h4 ResourceModel.ContainerHostVolume
	h4.Name = "workspace"
	h4.ContainerPath = "/workspace"
	var podPath string
	podPath = workspacePath
	// pod path
	podPath += "/" + user.Name
	podPath += "/" + job.Project.Name + "-" + strconv.FormatInt(job.Project.Id, 36)
	podPath += "/" + job.Name + "-" + strconv.FormatInt(job.Id, 36)
	podPath += "/" + module.Name + "-" + strconv.FormatInt(module.Id, 36)
	h4.HostPath = podPath
	h4.ReadOnly = false
	hostPathList = append(hostPathList, &h4)
	beego.Debug("->  " + h4.Name + " = " + h4.HostPath + ":" + h4.ContainerPath)

	beego.Debug("== create replication controller")
	beego.Debug("->namespace", namespace)
	beego.Debug("->label", label)
	beego.Debug("->image", image)
	beego.Debug("->containerPort", containerPort)
	beego.Debug("->limitCpu", limitCpu)
	beego.Debug("->limitMemory", limitMemory)
	beego.Debug("->envList", envList)
	beego.Debug("->hostPathList", hostPathList)

	err = systemApi.ApiCreateRc(namespace, label, image, containerPort, limitCpu, limitMemory, envList, hostPathList)
	if err != nil {
		beego.Debug("create rc failed")
	}

	// get pod name
	var podName string
	podName, err = systemApi.ApiGetPodName(namespace, label)
	if err != nil {
		beego.Debug("get pod name failed")
	}
	beego.Debug("== get pod name:", podName)
	pod.PodName = podName
	pod.Module = module
	pod.Module.Pods = nil
	beego.Debug("== update pod:", *pod)

	var newPod *model.Pod
	newPod, err = daoApi.BussinessDaoApi.UpdatePod(pod)
	if err != nil {
		beego.Debug("update pod name failed")
	}
	newPod.Module = nil

	// create svc
	beego.Debug("== create service")
	err = systemApi.ApiCreateSvc(namespace, newPod.SvcName, label, containerPort)
	if err != nil {
		beego.Debug("create svc failed")
	}

	// update resource used
	var resource *model.Resource
	resource, err = daoApi.ResourceDaoApi.GetByUserId(user.Id)
	if err != nil {
		beego.Debug("get resource failed")
	}
	resource.CpuUsageResource += newPod.Cpu
	resource.MemoryUsageResource += newPod.Memory

	_, err = daoApi.ResourceDaoApi.Update(resource)
	if err != nil {
		beego.Debug("update resource failed")
	}

	// load data
	var data model.Module
	data.Id = module.Id
	data.Name = module.Name
	data.InputFiles = module.InputFiles
	data.OutputFiles = module.OutputFiles
	data.Algorithm = module.Algorithm
	data.Parameters = module.Parameters
	data.Description = module.Description
	data.Pods = append(data.Pods, newPod)
	data.Job = module.Job

	var ip string
	var cfg = beego.AppConfig
	// todo
	// if config.DEBUG_ONLY {
	// 	ip = "http://192.168.0.206"
	// } else {
	// 	// label + namespace
	// 	ip = "http://" + name + "." + user.Name
	// }

	if cfg.String("runmode") == "dev" {
		ip = "http://192.168.0.206"
	} else {
		// label + namespace
		ip = "http://" + name + "." + user.Name
	}

	// connect
	var svc JobService
	err = svc.ConnectToPod(ip, containerPort)
	if err != nil {
		beego.Debug("connect to pod failed")
	}
	// send data
	beego.Debug("== send data to pod", ip)
	err = svc.SendDataToPod(ip, containerPort, &data)
	if err != nil {
		beego.Debug("send data failed")
	}

	ch <- "finish pod :" + newPod.Name + "-" + module.Name + "-" + job.Name
}

func (this *JobService) ConnectToPod(ip string, containerPort int) error {
	beego.Debug("connect to pod")
	beego.Debug("address :", ip, containerPort)
	return nil
}

func (this *JobService) SendDataToPod(ip string, containerPort int, data *model.Module) error {
	beego.Debug("send data to pod :", data.Pods[0].Name)

	var err error
	// send data to processor
	var msgReq model.MessageRequest
	// // header
	msgReq.Token = ""
	msgReq.SessionId = ""

	// // from
	var msgReqFrom model.MessageRole
	msgReqFrom.Type = model.MSG_ROLE_TYPE_SERVICE
	msgReqFrom.Id = model.BUSSINESS_SERVICE
	msgReq.From = &msgReqFrom

	// // to
	var msgReqTo model.MessageRole
	msgReqTo.Type = model.MSG_ROLE_TYPE_SYS_JOB_PROCESSOR
	msgReqTo.Id = data.Pods[0].Id
	msgReq.To = &msgReqTo

	// // parameter
	var msgParameter model.MessageParameter
	// parameter ---- action
	msgParameter.Action = model.MSG_PARAM_ACTION_CREATE
	// parameter ---- target
	msgParameter.Target = model.MSG_PARAM_TYPE_POD
	// parameter ---- data
	var msgParamData1 model.MessageParameterData
	msgParamData1.Type = model.MSG_PARAM_TYPE_POD
	var moduleProcessorJsonBytes []byte
	moduleProcessorJsonBytes, err = json.Marshal(data)
	if err != nil {
		// todo handle error
		beego.Debug("erro: ", err)
	}
	msgParamData1.Content = string(moduleProcessorJsonBytes)
	msgParameter.Data = append(msgParameter.Data, &msgParamData1)

	msgReq.Parameter = &msgParameter

	// convert data to string
	var msgReqJsonBytes []byte
	msgReqJsonBytes, _ = json.Marshal(&msgReq)
	//beego.Debug(string(msgReqJsonBytes))

	// send data to processor
	port := strconv.Itoa(containerPort)
	url := ip + ":" + port
	var resp []byte
	for i := 0; i < 20; i++ {
		resp, err = httpclient.Post(url, msgReqJsonBytes)
		if err != nil {
			beego.Debug(err)
			beego.Debug("wating for 5s")
			time.Sleep(5 * time.Second)
			beego.Debug("retry send data")
		} else {
			break
		}
	}

	if err != nil {
		return err
	}

	beego.Debug("recieve message from pod ", string(resp))

	return nil
}
