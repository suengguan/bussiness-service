package controllers

import (
	"app-service/bussiness-service/models"
	"app-service/bussiness-service/service"
	"encoding/json"
	"fmt"
	"model"

	"github.com/astaxie/beego"
)

// Operations about Job
type JobController struct {
	beego.Controller
}

// @Title Create
// @Description create job
// @Param	body		body 	models.Job	true		"body for job content"
// @Success 200 {object} models.Response
// @Failure 403 body is empty
// @router / [post]
func (this *JobController) Create() {
	var err error
	var job model.Job
	var response models.Response

	// unmarshal
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &job)
	if err == nil {
		var svc service.JobService
		var result []byte
		var newJob *model.Job
		newJob, err = svc.Create(&job)
		if err == nil {
			result, err = json.Marshal(newJob)
			if err == nil {
				response.Status = model.MSG_RESULTCODE_SUCCESS
				response.Reason = "success"
				response.Result = string(result)
			}
		}
	} else {
		beego.Debug("Unmarshal data failed")
		err = fmt.Errorf("%s", "Unmarshal data failed")
	}

	if err != nil {
		response.Status = model.MSG_RESULTCODE_FAILED
		response.Reason = err.Error()
		response.RetryCount = 3
	}

	this.Data["json"] = &response

	this.ServeJSON()
}

// @Title Update
// @Description update the job
// @Param	body		body 	models.Job	true		"body for user content"
// @Success 200 {object} models.Response
// @Failure 403 :id is invalid
// @router / [put]
func (this *JobController) Update() {
	var err error
	var job model.Job
	var response models.Response

	// unmarshal
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &job)
	if err == nil {
		var svc service.JobService
		var result []byte
		var newJob *model.Job
		newJob, err = svc.Update(&job)
		if err == nil {
			result, err = json.Marshal(newJob)
			if err == nil {
				response.Status = model.MSG_RESULTCODE_SUCCESS
				response.Reason = "success"
				response.Result = string(result)
			}
		}
	} else {
		beego.Debug("Unmarshal data failed")
	}

	if err != nil {
		response.Status = model.MSG_RESULTCODE_FAILED
		response.Reason = err.Error()
		response.RetryCount = 3
	}

	this.Data["json"] = &response

	this.ServeJSON()
}
