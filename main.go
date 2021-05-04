package main

import (
	"fmt"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/jankstar/go-skydisc/core"
)

var (
	goRouter *gin.Engine
	goServer *core.TServer
)

func main() {
	var err error
	//init Server and DB
	goServer, err = core.ServerInit(0, "")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Init DB in Mode %v\n", core.Server.Mode)

	gin.SetMode(gin.DebugMode) //gin.ReleaseMode)
	goRouter = gin.New()
	goRouter.Use(gin.Recovery())

	//Change template delimiter, because {{}} is used by vue
	goRouter.Delims("<(", ")>")
	// goRouter.StaticFile("favicon.ico", "favicon.ico")
	// goRouter.StaticFile("skydisc.png", "skydisc.png")
	goRouter.Use(static.Serve("/vendor", static.LocalFile("./client/vendor", false)))
	goRouter.Use(static.Serve("/icon", static.LocalFile("./client", false)))
	goRouter.LoadHTMLGlob("client/*.html")

	//routerGroup(goRouter, "/user")

	core.HTTPRouter(goRouter)

	goRouter.RunTLS(core.Server.Port, "./key/server.pem", "./key/server.key")
	//goRouter.Run(port)
}
