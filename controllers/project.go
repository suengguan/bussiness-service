package controllers

import (
	"app-service/bussiness-service/models"
	"app-service/bussiness-service/service"
	"encoding/json"
	"fmt"
	"model"

	"github.com/astaxie/beego"
)

// Operations about Project
type ProjectController struct {
	beego.Controller
}

// @Title Create
// @Description create project
// @Param	body		body 	models.Project	true		"body for project content"
// @Success 200 {object} models.Response
// @Failure 403 body is empty
// @router / [post]
func (this *ProjectController) Create() {
	var err error
	var project model.Project
	var response models.Response

	// unmarshal
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &project)
	if err == nil {
		var svc service.ProjectService
		var result []byte
		var newProject *model.Project
		newProject, err = svc.Create(&project)
		if err == nil {
			result, err = json.Marshal(newProject)
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

// @Title DeleteById
// @Description delete the project by id
// @Param	id		path 	int64	true		"The int you want to delete"
// @Success 200 {object} models.Response
// @Failure 403 :id is invalid
// @router /id/:id [delete]
func (this *ProjectController) DeleteById() {
	var err error
	var response models.Response

	var id int64
	id, err = this.GetInt64(":id")
	beego.Debug("DeleteById", id)
	if id > 0 && err == nil {
		var svc service.ProjectService
		err = svc.DeleteById(id)
		if err == nil {
			response.Status = model.MSG_RESULTCODE_SUCCESS
			response.Reason = "success"
			response.Result = ""
		}
	} else {
		beego.Debug(err)
		err = fmt.Errorf("%s", "project id is invalid")
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
// @Description update the project
// @Param	body		body 	models.Project	true		"body for user content"
// @Success 200 {object} models.Response
// @Failure 403 :id is invalid
// @router / [put]
func (this *ProjectController) Update() {
	var err error
	var project model.Project
	var response models.Response

	// unmarshal
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &project)
	if err == nil {
		var svc service.ProjectService
		var result []byte
		var newProject *model.Project
		newProject, err = svc.Update(&project)
		if err == nil {
			result, err = json.Marshal(newProject)
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

// @Title GetAll
// @Description get all user's projects
// @Param	userid		path 	int64	true		"The key for staticblock"
// @Success 200 {object} models.Response
// @router /:userId [get]
func (this *ProjectController) GetAll() {
	var err error
	var response models.Response

	var userId int64
	userId, err = this.GetInt64(":userId")
	beego.Debug("GetAll", userId)
	if userId > 0 && err == nil {
		var svc service.ProjectService
		var projects []*model.Project
		var result []byte
		projects, err = svc.GetAll(userId)
		if err == nil {
			result, err = json.Marshal(projects)
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

// @Title GetById
// @Description get project by id
// @Param	id		path 	int64	true		"The key for staticblock"
// @Success 200 {object} models.Response
// @Failure 403 :id is empty
// @router /id/:id [get]
func (this *ProjectController) GetById() {
	var err error
	var response models.Response

	var id int64
	id, err = this.GetInt64(":id")
	beego.Debug("GetById", id)
	if id > 0 && err == nil {
		var svc service.ProjectService
		var project *model.Project
		var result []byte
		project, err = svc.GetById(id)
		if err == nil {
			result, err = json.Marshal(project)
			if err == nil {
				response.Status = model.MSG_RESULTCODE_SUCCESS
				response.Reason = "success"
				response.Result = string(result)
			}
		}
	} else {
		beego.Debug(err)
		err = fmt.Errorf("%s", "project id is invalid")
	}

	if err != nil {
		response.Status = model.MSG_RESULTCODE_FAILED
		response.Reason = err.Error()
		response.RetryCount = 3
	}
	this.Data["json"] = &response

	this.ServeJSON()
}
