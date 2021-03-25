package catalog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/jankstar/go-skydisc/lib"
	"gorm.io/gorm"
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

func InitCatalogDB(iDB *gorm.DB, iMode int) error {

	//Catalogs
	iDB.AutoMigrate(&CatOrderClass{})
	iDB.AutoMigrate(&CatTrade{})
	iDB.AutoMigrate(&CatQualification{})

	if iMode == 1 {
		//Test-Modus - Daten initialisieren

		iDB.Where("class <> ''").Delete(&CatOrderClass{})
		iDB.Where("trade <> ''").Delete(&CatTrade{})
		iDB.Where("qualification <> ''").Delete(&CatQualification{})

		loadTestData(iDB)

	}

	return iDB.Error
}

func loadTestData(iDB *gorm.DB) {
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
	iDB.Save(&test.OrderClassList)
	iDB.Save(&test.QualificationList)
	iDB.Save(&test.TradeList)
	fmt.Println("loadTestData Catalog: in DB verbucht")
}
