package catalog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/jankstar/go-skydisc/lib"
)

//CatOrderClass - define Order class as customization
type CatOrderClass struct {
	Class string `json:"class" gorm:"primaryKey"`
	Name  string `json:"name"`
}

//CatTrade - define Order trade as customization
type CatTrade struct {
	Trade string `json:"trade" gorm:"primaryKey"`
	Name  string `json:"name"`
}

//CatQualification - define Qualification as customization
type CatQualification struct {
	Qualification string `json:"qualification" gorm:"primaryKey"`
	Name          string `json:"name"`
}

func InitCatalogDB(iMode int) error {

	//Catalogs
	lib.Server.DB.AutoMigrate(&CatOrderClass{})
	lib.Server.DB.AutoMigrate(&CatTrade{})
	lib.Server.DB.AutoMigrate(&CatQualification{})

	if iMode == 1 {
		//Test-Modus - Daten initialisieren

		lib.Server.DB.Where("class <> ''").Delete(&CatOrderClass{})
		lib.Server.DB.Where("trade <> ''").Delete(&CatTrade{})
		lib.Server.DB.Where("qualification <> ''").Delete(&CatQualification{})

		loadTestData()

	}

	return lib.Server.DB.Error
}

func loadTestData() {
	var test struct {
		TradeList         []CatTrade         `json:"trade_list"`
		QualificationList []CatQualification `json:"qualification_list"`
		OrderClassList    []CatOrderClass    `json:"class_list"`
	}
	fmt.Println("loadTestData catalog: ", lib.Server.TestfileCatalog)
	data, err := ioutil.ReadFile(lib.Server.TestfileCatalog)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	err = json.Unmarshal(data, &test)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	lib.Server.DB.Save(&test.OrderClassList)
	lib.Server.DB.Save(&test.QualificationList)
	lib.Server.DB.Save(&test.TradeList)
	fmt.Println("loadTestData Catalog: in DB verbucht")
}
