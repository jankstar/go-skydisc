package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

//DataServiceArea Service area
//describes e.g. a branch office, temporal allocation of certain resources
//* Name, Designation
//* Address, geo-coordinates for starting point
//* Capacities per qualification
type DataServiceArea struct {
	ID          string `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string            `json:"name"`
	Location    TLocation         `json:"location" gorm:"embedded"`
	Requirement []DataRequirement `json:"requirement" gorm:"foreignKey:ServiceAreaRef"`
}

func InitOrgaDB(iMode int) error {

	//Data
	Server.DB.AutoMigrate(&DataServiceArea{})

	if iMode == 1 {
		//Test-Modus - Daten initialisieren
		Server.DB.Where("id <> ''").Delete(&DataServiceArea{})
		//lib.Server.DB.Where("id <> '' AND service_area_refer <> ''").Delete(&order.DataRequirement{})

		loadTestDataOrga()
	}
	return Server.DB.Error
}

//*********************************

func loadTestDataOrga() {
	var test struct {
		ServiceArea []DataServiceArea `json:"service_area_list"`
	}
	fmt.Println("loadTestData: orga.json")
	data, err := ioutil.ReadFile(Server.Path + Server.TestfileOrga)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	err = json.Unmarshal(data, &test)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	Server.DB.Save(&test.ServiceArea)

	fmt.Println("loadTestData: in DB verbucht")
}
