package web

import (
	"github.com/gin-gonic/gin"
	"go-common-web/controllers"
	controllers2 "go-swe/src/web/controllers"
)

func RegisterRouter(r *gin.Engine) {

	RegisterControllers()

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "home/index.html", gin.H{})
	})

	r.GET("/jd", controllers.ControllerHandle("JDController", "Convert"))

	r.GET("/solar/terms/:year", controllers.ControllerHandle("SolarController", "TermsByYear"))

	r.GET("/solar/terms", controllers.ControllerHandle("SolarController", "TermsByRange"))

	r.GET("/lunar/phases/", controllers.ControllerHandle("LunarController", "PhasesByRange"))
	r.GET("/lunar/phases/:year", controllers.ControllerHandle("LunarController", "PhasesByYear"))
	r.GET("/lunar/months/:year", controllers.ControllerHandle("LunarController", "MonthsByYear"))
}

func RegisterControllers() {
	controllers.RegisterController("SolarController", func(ctx *gin.Context) controllers.ControllerInterface {
		return &controllers2.SolarController{Controller: controllers.Controller{Context: ctx}}
	})

	controllers.RegisterController("JDController", func(ctx *gin.Context) controllers.ControllerInterface {
		return &controllers2.JDController{Controller: controllers.Controller{Context: ctx}}
	})

	controllers.RegisterController("LunarController", func(ctx *gin.Context) controllers.ControllerInterface {
		return &controllers2.LunarController{Controller: controllers.Controller{Context: ctx}}
	})
}
