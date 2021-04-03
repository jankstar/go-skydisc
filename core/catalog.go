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

//CatOrderClass - define Order class as customization
type CatOrderStatus struct {
	Status  string `json:"status" gorm:"primaryKey"`
	Name    string `json:"name"`
	Default bool   `json:"default"`
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

//CatSection - define Appointment Section as customization
type CatSection struct {
	Section uint   `json:"section" gorm:"primaryKey"`
	Name    string `json:"name"`
}

func InitCatalogDB(iMode int) error {

	//Catalogs
	Server.DB.AutoMigrate(&CatOrderClass{})
	Server.DB.AutoMigrate(&CatOrderStatus{})
	Server.DB.AutoMigrate(&CatTrade{})
	Server.DB.AutoMigrate(&CatQualification{})
	Server.DB.AutoMigrate(&CatSection{})

	if iMode == 1 {
		//Test-Modus - Daten initialisieren

		Server.DB.Where("class <> ''").Delete(&CatOrderClass{})
		Server.DB.Where("status <> ''").Delete(&CatOrderStatus{})
		Server.DB.Where("trade <> ''").Delete(&CatTrade{})
		Server.DB.Where("qualification <> ''").Delete(&CatQualification{})
		Server.DB.Where("section <> ''").Delete(&CatSection{})

		loadTestDataCatalog()

	}

	return Server.DB.Error
}

func loadTestDataCatalog() {
	var test struct {
		TradeList         []CatTrade         `json:"trade_list"`
		QualificationList []CatQualification `json:"qualification_list"`
		OrderClassList    []CatOrderClass    `json:"class_list"`
		OrderStatusList   []CatOrderStatus   `json:"status_list"`
		SectionList       []CatSection       `json:"section_list"`
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
	Server.DB.Save(&test.OrderStatusList)
	Server.DB.Save(&test.QualificationList)
	Server.DB.Save(&test.TradeList)
	Server.DB.Save(&test.SectionList)
	fmt.Println("loadTestData Catalog: in DB verbucht")
}
