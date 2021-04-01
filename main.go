package main

import (
	"fmt"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/jankstar/go-skydisc/core"
)

//returns context index.html
func indexFunc(iCon *gin.Context) {
}

func main() {

	//init Server and DB
	err := core.ServerInit(1, "")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Init DB in Mode %v\n", core.Server.Mode)

	core.InitCatalogDB(core.Server.Mode)
	core.InitLocationDB(core.Server.Mode)
	core.InitRequirementDB(core.Server.Mode)
	core.InitOrgaDB(core.Server.Mode)
	core.InitCalendarDB(core.Server.Mode)
	core.InitResourceDB(core.Server.Mode)
	core.InitOrderDB(core.Server.Mode)
	core.InitAppointmentDB(core.Server.Mode)

	gin.SetMode(gin.DebugMode) //gin.ReleaseMode)
	oRouter := gin.New()
	oRouter.Use(gin.Recovery())

	//Change template delimiter, because {{}} is used by vue
	oRouter.Delims("<(", ")>")
	// oRouter.StaticFile("favicon.ico", "favicon.ico")
	// oRouter.StaticFile("lookinlogo.png", "lookinlogo.png")
	oRouter.Use(static.Serve("/vendor", static.LocalFile("./client/vendor", false)))
	oRouter.Use(static.Serve("/icon", static.LocalFile("./client", false)))
	oRouter.LoadHTMLGlob("client/*.html")

	//routerGroup(oRouter, "/user")

	oRouter.GET("/", indexFunc)
	oRouter.RunTLS(core.Server.Port, "./key/server.pem", "./key/server.key")
	//oRouter.Run(port)
}
