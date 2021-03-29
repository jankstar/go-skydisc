package main

import (
	"fmt"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jankstar/go-skydisc/catalog"
	"github.com/jankstar/go-skydisc/lib"
	"github.com/jankstar/go-skydisc/order"
)

//returns context index.html
func indexFunc(iCon *gin.Context) {
}

func main() {

	//init Server and DB
	err := lib.ServerInit(1, "")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Init DB in Mode %v\n", lib.Server.Mode)

	catalog.InitCatalogDB(lib.Server.Mode)
	order.InitOrderDB(lib.Server.Mode)

	gin.SetMode(gin.DebugMode) //gin.ReleaseMode)
	oRouter := gin.New()
	oRouter.Use(gin.Recovery())

	//Change template delimiter, because {{}} is used by vue
	oRouter.Delims("<(", ")>")
	oRouter.StaticFile("favicon.ico", "favicon.ico")
	oRouter.StaticFile("lookinlogo.png", "lookinlogo.png")
	oRouter.Use(static.Serve("/vendor", static.LocalFile("./client/vendor", false)))
	oRouter.Use(static.Serve("/icon", static.LocalFile("./client", false)))
	oRouter.LoadHTMLGlob("client/*.html")

	//routerGroup(oRouter, "/user")

	oRouter.GET("/", indexFunc)
	oRouter.RunTLS(lib.Server.Port, "./key/server.pem", "./key/server.key")
	//oRouter.Run(port)
}
