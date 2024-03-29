package web

import (
	"github.com/gin-gonic/gin"
	innerControllers "go-swe/src/web/controllers"
	"gopkg.in/go-mixed/go-common.v1/web.v1/controllers"
)

func RegisterRouter(r *gin.Engine) {

	RegisterControllers()

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "home/index.html", gin.H{})
	})

	r.GET("/jd", controllers.ControllerHandler("JDController", "Convert"))

	r.GET("/solar/terms/:year", controllers.ControllerHandler("SolarController", "TermsByYear"))

	r.GET("/solar/terms", controllers.ControllerHandler("SolarController", "TermsByRange"))

	r.GET("/lunar/phases/", controllers.ControllerHandler("LunarController", "PhasesByRange"))
	r.GET("/lunar/phases/:year", controllers.ControllerHandler("LunarController", "PhasesByYear"))
	r.GET("/lunar/months/:year", controllers.ControllerHandler("LunarController", "MonthsByYear"))
}

func RegisterControllers() {
	controllers.RegisterController("SolarController", func(ctx *gin.Context) controllers.IController {
		return &innerControllers.SolarController{Controller: controllers.Controller{Context: ctx}}
	})

	controllers.RegisterController("JDController", func(ctx *gin.Context) controllers.IController {
		return &innerControllers.JDController{Controller: controllers.Controller{Context: ctx}}
	})

	controllers.RegisterController("LunarController", func(ctx *gin.Context) controllers.IController {
		return &innerControllers.LunarController{Controller: controllers.Controller{Context: ctx}}
	})
}
