package controllers

import (
	"app-service/bussiness-service/models"
	"app-service/bussiness-service/service"
	"encoding/json"
	"fmt"
	"model"

	"github.com/astaxie/beego"
)

// Operations about Pod
type PodController struct {
	beego.Controller
}

// @Title GetCurrent
// @Description get all current pods
// @Param	userid		path 	int64	true		"The key for staticblock"
// @Success 200 {object} models.Response
// @Failure 403 :userid is invalid
// @router /current/:userid [get]
func (this *PodController) GetCurrent() {
	var err error
	var response models.Response

	var userid int64
	userid, err = this.GetInt64(":userid")
	beego.Debug("GetCurrent", userid)
	if userid > 0 && err == nil {
		var svc service.PodService
		var pods []*model.Pod
		var result []byte
		pods, err = svc.GetCurrent(userid)
		if err == nil {
			result, err = json.Marshal(pods)
			if err == nil {
				response.Status = model.MSG_RESULTCODE_SUCCESS
				response.Reason = "success"
				response.Result = string(result)
			}
		}
	} else {
		beego.Debug(err)
		err = fmt.Errorf("%s", "user id is invalid")
	}

	if err != nil {
		response.Status = model.MSG_RESULTCODE_FAILED
		response.Reason = err.Error()
		response.RetryCount = 3
	}
	this.Data["json"] = &response

	this.ServeJSON()
}

// @Title UpdateStatus
// @Description update pod status
// @Param	body		body 	models.PodStatus	true		"body for pod status content"
// @Success 200 {object} models.Response
// @Failure 403 :id is invalid
// @router /status/ [post]
func (this *PodController) UpdateStatus() {
	var err error
	var podStatus model.PodStatus
	var response models.Response

	// unmarshal
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &podStatus)
	if err == nil {
		var svc service.PodService
		err = svc.UpdateStatus(&podStatus)
		if err == nil {
			response.Status = model.MSG_RESULTCODE_SUCCESS
			response.Reason = "success"
			response.Result = ""
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
