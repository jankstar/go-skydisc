package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	Server.DB.AutoMigrate(&CatOrderClass{})
	Server.DB.AutoMigrate(&CatTrade{})
	Server.DB.AutoMigrate(&CatQualification{})

	if iMode == 1 {
		//Test-Modus - Daten initialisieren

		Server.DB.Where("class <> ''").Delete(&CatOrderClass{})
		Server.DB.Where("trade <> ''").Delete(&CatTrade{})
		Server.DB.Where("qualification <> ''").Delete(&CatQualification{})

		loadTestDataCatalog()

	}

	return Server.DB.Error
}

func loadTestDataCatalog() {
	var test struct {
		TradeList         []CatTrade         `json:"trade_list"`
		QualificationList []CatQualification `json:"qualification_list"`
		OrderClassList    []CatOrderClass    `json:"class_list"`
	}
	fmt.Println("loadTestData catalog: ", Server.TestfileCatalog)
	data, err := ioutil.ReadFile(Server.Path + Server.TestfileCatalog)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	err = json.Unmarshal(data, &test)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	Server.DB.Save(&test.OrderClassList)
	Server.DB.Save(&test.QualificationList)
	Server.DB.Save(&test.TradeList)
	fmt.Println("loadTestData Catalog: in DB verbucht")
}
