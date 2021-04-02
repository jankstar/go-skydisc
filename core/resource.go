package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type DataResource struct {
	ID               string    `json:"id" gorm:"primaryKey"`
	ValidFrom        time.Time `json:"valid_from" gorm:"primaryKey"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ValidTo          time.Time               `json:"valid_to"`
	Name             string                  `json:"name"`
	HomeLocation     TLocation               `json:"home_location" gorm:"embedded"`
	Capacity         []DataRequirement       `json:"capacity" gorm:"foreignKey:ResourceRef"`
	CalendarRef      string                  `json:"calendar_refer"`
	Calendar         DataWorkingTimeCalendar `json:"calendar" gorm:"foreignKey:CalendarRef"`
	CapacityCalendar []DataCapacityCalendar  `json:"capacity_calendar" gorm:"foreignKey:CalendarRef"`
}

type DataCapacityCalendar struct {
	RecourceRef    string     `json:"id" gorm:"primaryKey"`
	Date           string     `json:"date" gorm:"primaryKey"`
	SectionRef     uint       `json:"section_ref" gorm:"primaryKey"`
	Section        CatSection `json:"section" gorm:"foreignKey:SectionRef"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DurationTarget time.Duration `json:"DurationTarget"`
	DurationRest   time.Duration `json:"DurationRest"`
}

func InitResourceDB(iMode int) error {

	//Data
	Server.DB.AutoMigrate(&DataResource{})

	if iMode == 1 {
		//Test-Modus - Daten initialisieren
		Server.DB.Where("id <> ''").Delete(&DataResource{})

		loadTestDataResource()
	}
	return Server.DB.Error
}

func loadTestDataResource() {
	var test struct {
		Resources []DataResource `json:"resource_list"`
	}
	fmt.Println("loadTestData: resource.json")
	data, err := ioutil.ReadFile(Server.Path + Server.TestfileResource)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	err = json.Unmarshal(data, &test)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	Server.DB.Save(&test.Resources)

	fmt.Println("loadTestData: in DB verbucht")
	return
}
