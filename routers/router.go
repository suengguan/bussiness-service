// @APIVersion 1.0.0
// @Title bussiness-service API
// @Description bussiness-service only serve account register/delete/update/get
// @Contact qsg@corex-tek.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"app-service/bussiness-service/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/project",
			beego.NSInclude(
				&controllers.ProjectController{},
			),
		),
		beego.NSNamespace("/job",
			beego.NSInclude(
				&controllers.JobController{},
			),
		),
		beego.NSNamespace("/pod",
			beego.NSInclude(
				&controllers.PodController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
