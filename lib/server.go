package lib

//Library for central processing and server functions

import (
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	GfHost = "127.0.0.1"
	GfPort = ":8000"

	GfDBName   = "tmp/test.db"
	GsDBConfig = &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		Logger: logger.Default.LogMode(logger.Silent),
	}
)

//externalIP determines the external IP address
func getExternalIP() (string, error) {
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

func ServerInit(iMode int, iPath string) (err error) {

	godotenv.Load()

	if os.Getenv("port") != "" {
		GfPort = ":" + os.Getenv("port")
	}

	//IP Determine address
	if ip, err := getExternalIP(); err == nil {
		GfHost = ip
		fmt.Println("IP address set:", ip)
	}

	return
}
