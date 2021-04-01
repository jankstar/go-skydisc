package core

//Library for central processing and server functions

import (
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	Server = TServer{
		Mode: 0,
		Host: "127.0.0.1",
		Port: ":8000",
		//
		DBName: "tmp/test.db",
		DBConfig: &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
			},
			Logger: logger.Default.LogMode(logger.Silent)},
		//
		Path:             "",
		TestfileCatalog:  "tmp/catalog.json",
		TestfileOrder:    "tmp/order.json",
		TestfileOrga:     "tmp/orga.json",
		TestfileResource: "tmp/resource.json",
		TestfileCalendar: "tmp/calendar.json",
		//
		BingURLLocation: "https://dev.virtualearth.net/REST/v1/Locations/%s/%s/%s/%s?" +
			"includeNeighborhood=1&include=ciso2&maxResults=%d&key=%s",
		BingURLTimezone: "https://dev.virtualearth.net/REST/v1/TimeZone/?query=%s&key=%s",
		BingApiKey:      "", //put api key in .env file
	}
)

type TServer struct {
	Mode             int
	Host             string
	Port             string
	DB               *gorm.DB
	DBName           string
	DBConfig         *gorm.Config
	Path             string
	TestfileCatalog  string
	TestfileOrder    string
	TestfileOrga     string
	TestfileResource string
	TestfileCalendar string
	BingURLLocation  string
	BingURLTimezone  string
	BingApiKey       string
}

//externalIP determines the external IP address
func GetExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("Are you connected to the network?")
}

//ServerInit Initialize server, read .env variables
// Initialize DB
// iMode int -> 0:run or 1:Testdaten
// iPath string -> path for data files
func ServerInit(iMode int, iPath string) (err error) {

	Server.Mode = iMode
	Server.Path = iPath

	//load and handle .env variables
	err = godotenv.Load(iPath + ".env")

	if os.Getenv("port") != "" {
		Server.Port = ":" + os.Getenv("port")
	}

	//IP Determine address
	if ip, err := GetExternalIP(); err == nil {
		Server.Host = ip
		fmt.Println("IP address set:", ip)
	}

	if os.Getenv("bing_api_key") != "" {
		Server.BingApiKey = os.Getenv("bing_api_key")
	}

	//init DB
	Server.DB, err = gorm.Open(sqlite.Open(iPath+Server.DBName), Server.DBConfig)
	return
}
