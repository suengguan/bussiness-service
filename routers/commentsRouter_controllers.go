package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["app-service/bussiness-service/controllers:JobController"] = append(beego.GlobalControllerRouter["app-service/bussiness-service/controllers:JobController"],
		beego.ControllerComments{
			Method: "Create",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["app-service/bussiness-service/controllers:JobController"] = append(beego.GlobalControllerRouter["app-service/bussiness-service/controllers:JobController"],
		beego.ControllerComments{
			Method: "Update",
			Router: `/`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["app-service/bussiness-service/controllers:PodController"] = append(beego.GlobalControllerRouter["app-service/bussiness-service/controllers:PodController"],
		beego.ControllerComments{
			Method: "GetCurrent",
			Router: `/current/:userid`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["app-service/bussiness-service/controllers:PodController"] = append(beego.GlobalControllerRouter["app-service/bussiness-service/controllers:PodController"],
		beego.ControllerComments{
			Method: "UpdateStatus",
			Router: `/status/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["app-service/bussiness-service/controllers:ProjectController"] = append(beego.GlobalControllerRouter["app-service/bussiness-service/controllers:ProjectController"],
		beego.ControllerComments{
			Method: "Create",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["app-service/bussiness-service/controllers:ProjectController"] = append(beego.GlobalControllerRouter["app-service/bussiness-service/controllers:ProjectController"],
		beego.ControllerComments{
			Method: "DeleteById",
			Router: `/id/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["app-service/bussiness-service/controllers:ProjectController"] = append(beego.GlobalControllerRouter["app-service/bussiness-service/controllers:ProjectController"],
		beego.ControllerComments{
			Method: "Update",
			Router: `/`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["app-service/bussiness-service/controllers:ProjectController"] = append(beego.GlobalControllerRouter["app-service/bussiness-service/controllers:ProjectController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/:userId`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["app-service/bussiness-service/controllers:ProjectController"] = append(beego.GlobalControllerRouter["app-service/bussiness-service/controllers:ProjectController"],
		beego.ControllerComments{
			Method: "GetById",
			Router: `/id/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

}
