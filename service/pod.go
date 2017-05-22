package service

import (
	"model"

	daoApi "api/dao_service"
	systemApi "api/system_service"

	"github.com/astaxie/beego"
)

type PodService struct {
}

func (this *PodService) GetCurrent(userId int64) ([]*model.Pod, error) {
	var err error
	var pods []*model.Pod

	// get user
	var user *model.User
	user, err = daoApi.UserDaoApi.GetById(userId)
	if err != nil {
		beego.Debug("get user failed")
		return nil, err
	}

	// get all projects
	var projects []*model.Project
	projects, err = daoApi.BussinessDaoApi.GetAllProjects(userId)
	if err != nil {
		beego.Debug("get project by user failed")
		return nil, err
	}

	// get current pods name
	var currentPodsName []string
	currentPodsName, err = systemApi.ApiGetPodsByNamespace(user.Name)
	if err != nil {
		beego.Debug("get current pods name failed")
		return nil, err
	}

	// check current names
	if len(currentPodsName) <= 0 {
		return nil, err
	}

	for _, project := range projects {
		for _, j := range project.Jobs {
			for _, m := range j.Modules {
				for _, p := range m.Pods {
					//p.PodName
					for _, pn := range currentPodsName {
						if p.PodName == pn {
							pods = append(pods, p)
							break
						}
					}
				}
			}
		}
	}

	return pods, err
}

func (this *PodService) UpdateStatus(status *model.PodStatus) error {
	var err error

	// find pod by id
	beego.Debug("->get pod")
	var pod *model.Pod
	pod, err = daoApi.BussinessDaoApi.GetPodById(status.Id)
	if err != nil {
		beego.Debug(err)
		return err
	}

	// find module by id
	beego.Debug("->get module")
	var module *model.Module
	module, err = daoApi.BussinessDaoApi.GetModuleById(pod.Module.Id)
	if err != nil {
		beego.Debug(err)
		return err
	}

	// get job
	beego.Debug("->get job")
	var job *model.Job
	job, err = daoApi.BussinessDaoApi.GetJobById(module.Job.Id)
	if err != nil {
		beego.Debug(err)
		return err
	}

	// get project
	beego.Debug("->get project")
	var project *model.Project
	project, err = daoApi.BussinessDaoApi.GetProjectById(job.Project.Id)
	if err != nil {
		beego.Debug(err)
		return err
	}

	// get user
	beego.Debug("->get user")
	var user *model.User
	user, err = daoApi.UserDaoApi.GetById(project.User.Id)
	if err != nil {
		beego.Debug(err)
		return err
	}

	if status.Status == model.JOB_STATUS_TYPE_FINISH ||
		status.Status == model.JOB_STATUS_TYPE_ERROR {
		// 通过status and reason => 删除pod/重启pod 等操作
		if status.Status == model.JOB_STATUS_TYPE_FINISH {
			beego.Debug("pod run finish")
		} else if status.Status == model.JOB_STATUS_TYPE_ERROR {
			beego.Debug("pod run error")
		}

		beego.Debug("->delete rc")
		err = systemApi.ApiDeleteRc(user.Name, pod.RcName)
		if err == nil {
			beego.Debug("delete rc success", pod.RcName)
		} else {
			beego.Debug("delete rc failed", pod.RcName)
		}

		beego.Debug("->delete svc")
		err = systemApi.ApiDeleteSvc(user.Name, pod.SvcName)
		if err == nil {
			beego.Debug("delete service success", pod.SvcName)
		} else {
			beego.Debug("delete service failed", pod.SvcName)
		}

		// update resource used
		beego.Debug("->get resource")
		var resource *model.Resource
		resource, err = daoApi.ResourceDaoApi.GetByUserId(user.Id)
		if err != nil {
			beego.Debug("get resource failed")
			return err
		}

		// update resource
		beego.Debug("->update resource")
		resource.CpuUsageResource -= pod.Cpu
		resource.MemoryUsageResource -= pod.Memory
		_, err = daoApi.ResourceDaoApi.Update(resource)
		if err != nil {
			beego.Debug("update resource failed")
			return err
		}

	}

	return err
}
