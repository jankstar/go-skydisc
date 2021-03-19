package main

import (
	"fmt"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jankstar/go-skydisc/lib"
	"github.com/jankstar/go-skydisc/order"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	host = "192.168.0.119"
	port = ":8000"
)

//returns context index.html
func indexFunc(iCon *gin.Context) {
}

func main() {
	godotenv.Load()

	if os.Getenv("port") != "" {
		port = ":" + os.Getenv("port")
	}

	//IP Determine address
	if ip, err := lib.GetExternalIP(); err == nil {
		host = ip
		fmt.Println("IP address set:", ip)
	}

	loDB, err := gorm.Open(sqlite.Open("test.db"), lib.GsDBConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	order.InitOrderDB(loDB)

	gin.SetMode(gin.DebugMode) //gin.ReleaseMode)
	oRouter := gin.New()
	oRouter.Use(gin.Recovery())

	//Change template delimiter, because {{}} is used by vue
	oRouter.Delims("<(", ")>")
	oRouter.StaticFile("favicon.ico", "favicon.ico")
	oRouter.StaticFile("lookinlogo.png", "lookinlogo.png")
	//oRouter.Use(static.Serve("/tmp", static.LocalFile(tmpDir, false)))
	oRouter.Use(static.Serve("/vendor", static.LocalFile("./client/vendor", false)))
	oRouter.Use(static.Serve("/icon", static.LocalFile("./client", false)))
	oRouter.LoadHTMLGlob("client/*.html")

	//routerGroup(oRouter, "/user")

	oRouter.GET("/", indexFunc)
	oRouter.RunTLS(port, "./key/server.pem", "./key/server.key")
	//oRouter.Run(port)
}
