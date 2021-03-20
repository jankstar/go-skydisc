package main

import (
	"fmt"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jankstar/go-skydisc/lib"
	"github.com/jankstar/go-skydisc/order"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//returns context index.html
func indexFunc(iCon *gin.Context) {
}

func main() {

	err := lib.ServerInit(1, "")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Init DB in Mode %v\n", lib.Server.Mode)
	loDB, _ := gorm.Open(sqlite.Open(lib.Server.DBName), lib.Server.DBConfig)
	order.InitOrderDB(loDB, lib.Server.Mode)

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
