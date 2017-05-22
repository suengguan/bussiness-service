package service

import (
	"model"

	daoApi "api/dao_service"

	"github.com/astaxie/beego"
)

type ProjectService struct {
}

func (this *ProjectService) Create(project *model.Project) (*model.Project, error) {
	var err error
	var newProject *model.Project

	beego.Debug("->create project")
	newProject, err = daoApi.BussinessDaoApi.CreateProject(project)
	if err != nil {
		return nil, err
	}

	// create project dir

	beego.Debug("result:", *newProject)

	return newProject, err
}

func (this *ProjectService) GetAll(userId int64) ([]*model.Project, error) {
	var err error
	var projects []*model.Project

	projects, err = daoApi.BussinessDaoApi.GetAllProjects(userId)
	if err != nil {
		return nil, err
	}

	return projects, err
}

func (this *ProjectService) GetById(id int64) (*model.Project, error) {
	var err error
	var project *model.Project

	project, err = daoApi.BussinessDaoApi.GetProjectById(id)
	if err != nil {
		return nil, err
	}

	return project, err
}

func (this *ProjectService) Update(project *model.Project) (*model.Project, error) {
	var err error
	var newProject *model.Project

	newProject, err = daoApi.BussinessDaoApi.UpdateProject(project)
	if err != nil {
		return nil, err
	}

	return newProject, err
}

func (this *ProjectService) DeleteById(id int64) error {
	var err error

	err = daoApi.BussinessDaoApi.DeleteProjectById(id)
	if err != nil {
		return err
	}

	return err
}

func (this *ProjectService) DeleteAll(userId int64) error {
	var err error

	err = daoApi.BussinessDaoApi.DeleteAllProjects(userId)
	if err != nil {
		return err
	}

	return err
}
