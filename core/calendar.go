package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type TWorkingTimeDay struct {
	StartMorning            string          `json:"start_morning"`
	EndMoning               string          `json:"end_moning"`
	ServiceAreaMorningRef   string          `json:"service_area_morning_ref"`
	ServiceAreaMorning      DataServiceArea `json:"service_area_morning" gorm:"foreignKey:ServiceAreaMorningRef"`
	StartAfternoon          string          `json:"start_afternoon"`
	EndAfternoon            string          `json:"end_afternoon"`
	ServiceAreaAfternoonRef string          `json:"service_area_afternoon_ref"`
	ServiceAreaAfternoon    DataServiceArea `json:"service_area_afternoon" gorm:"foreignKey:ServiceAreaAfternoonRef"`
}

type TWorkingTimeWeek struct {
	Sunday    TWorkingTimeDay `json:"sunday" gorm:"embedded;embeddedPrefix:sun_"`
	Monday    TWorkingTimeDay `json:"monday" gorm:"embedded;embeddedPrefix:mon_"`
	Tuesday   TWorkingTimeDay `json:"tuesday" gorm:"embedded;embeddedPrefix:tue_"`
	Wednesday TWorkingTimeDay `json:"wednesday" gorm:"embedded;embeddedPrefix:wen_"`
	Thursday  TWorkingTimeDay `json:"thursday" gorm:"embedded;embeddedPrefix:thu_"`
	Friday    TWorkingTimeDay `json:"friday" gorm:"embedded;embeddedPrefix:fri_"`
	Saturday  TWorkingTimeDay `json:"saturday" gorm:"embedded;embeddedPrefix:sat_"`
}

type DataWorkingTimeCalendar struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	ValidFrom time.Time `json:"valid_from" gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string           `json:"name"`
	Week      TWorkingTimeWeek `json:"week" gorm:"embedded"`
	OddWeek   TWorkingTimeWeek `json:"odd_week" gorm:"embedded;embeddedPrefix:odd_"`
}

func InitCalendarDB(iMode int) error {

	//Data
	Server.DB.AutoMigrate(&DataWorkingTimeCalendar{})

	if iMode == 1 {
		//Test-Modus - Daten initialisieren
		Server.DB.Where("id <> ''").Delete(&DataWorkingTimeCalendar{})

		loadTestDataCalendar()
	}
	return Server.DB.Error
}

func loadTestDataCalendar() {
	var test struct {
		Calendar []DataWorkingTimeCalendar `json:"calendar_list"`
	}
	fmt.Println("loadTestData: calendar.json")
	data, err := ioutil.ReadFile(Server.Path + Server.TestfileCalendar)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	err = json.Unmarshal(data, &test)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	Server.DB.Save(&test.Calendar)

	fmt.Println("loadTestData: in DB verbucht")
	return
}
