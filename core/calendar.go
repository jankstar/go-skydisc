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

type DataAbsencePeriods struct {
	ID             uint `json:"id" gorm:"primaryKey; autoIncrement"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Start          time.Time `json:"start"`
	End            time.Time `json:"end"`
	Description    string    `json:"description"`
	ServiceAreaRef string
	RecourceRef    string
}

func InitCalendarDB(iMode int) error {

	//Data
	Server.DB.AutoMigrate(&DataWorkingTimeCalendar{})
	Server.DB.AutoMigrate(&DataAbsencePeriods{})

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

func GetDayTime(iDay time.Time, iTimeStr string) (eTime time.Time, err error) {
	return time.Parse(time.RFC3339, iDay.Format("2006-01-02")+"T"+iTimeStr+"Z")
}

func CalculateDurationByStr(iDay time.Time, iStart string, iEnd string) time.Duration {
	lStart, err := GetDayTime(iDay, iStart)
	lEnd, err2 := GetDayTime(iDay, iEnd)
	if err == nil && err2 == nil && lEnd.After(lStart) == true {
		return lEnd.Sub(lStart)
	}
	return 0
}

func CalculateDurationByTime(iDay time.Time, iStart time.Time, iEnd time.Time) time.Duration {
	if iEnd.After(iStart) == true {
		return iEnd.Sub(iStart)
	}
	return 0
}

//Time2Date convert Time to date as string, e.g. "2000-01-01"
func Time2Date(iTime time.Time) string {
	return iTime.Format("2006-01-02")
}

func Date2Time(iDate string) (eTime time.Time) {
	eTime, _ = time.Parse("2006-01-02", iDate)
	return
}

//Time2Time convert Time to date as string e.g. "08:00:00"
func Time2Time(iTime time.Time) string {
	return iTime.Format("15:04:05")
}

//isInitial checks if the calendar is still initial
func (me *TWorkingTimeWeek) isInitial() bool {

	return me.Sunday.StartMorning == "" && me.Sunday.StartAfternoon == "" &&
		me.Monday.StartMorning == "" && me.Monday.StartAfternoon == "" &&
		me.Tuesday.StartMorning == "" && me.Tuesday.StartAfternoon == "" &&
		me.Wednesday.StartMorning == "" && me.Wednesday.StartAfternoon == "" &&
		me.Thursday.StartMorning == "" && me.Thursday.StartAfternoon == "" &&
		me.Friday.StartMorning == "" && me.Friday.StartAfternoon == "" &&
		me.Saturday.StartMorning == "" && me.Saturday.StartAfternoon == ""
}

//GetDrationByDay provides the duration for one day
func (me *DataWorkingTimeCalendar) GetDurationByDay(iDay time.Time) (eMorningDuration time.Duration, eAfternoonDuration time.Duration,
	eStartMorning time.Time, eEndMoning time.Time, eStartAfternoon time.Time, eEndAfternoon time.Time) {
	if me.ValidFrom.After(time.Now()) != true {
		var myWeek *TWorkingTimeWeek
		_, lWeek := time.Now().ISOWeek()
		lWeekRemainder := lWeek % 2
		if lWeekRemainder == 1 &&
			me.OddWeek.isInitial() == false {
			//odd week
			myWeek = &me.OddWeek
		} else {
			myWeek = &me.Week
		}
		switch iDay.Weekday() {
		case time.Sunday:
			eStartMorning, _ = GetDayTime(iDay, myWeek.Sunday.StartMorning)
			eEndMoning, _ = GetDayTime(iDay, myWeek.Sunday.EndMoning)
			eStartAfternoon, _ = GetDayTime(iDay, myWeek.Sunday.StartAfternoon)
			eEndAfternoon, _ = GetDayTime(iDay, myWeek.Sunday.EndAfternoon)

		case time.Monday:
			eStartMorning, _ = GetDayTime(iDay, myWeek.Monday.StartMorning)
			eEndMoning, _ = GetDayTime(iDay, myWeek.Monday.EndMoning)
			eStartAfternoon, _ = GetDayTime(iDay, myWeek.Monday.StartAfternoon)
			eEndAfternoon, _ = GetDayTime(iDay, myWeek.Monday.EndAfternoon)

		case time.Tuesday:
			eStartMorning, _ = GetDayTime(iDay, myWeek.Tuesday.StartMorning)
			eEndMoning, _ = GetDayTime(iDay, myWeek.Tuesday.EndMoning)
			eStartAfternoon, _ = GetDayTime(iDay, myWeek.Tuesday.StartAfternoon)
			eEndAfternoon, _ = GetDayTime(iDay, myWeek.Tuesday.EndAfternoon)

		case time.Wednesday:
			eStartMorning, _ = GetDayTime(iDay, myWeek.Wednesday.StartMorning)
			eEndMoning, _ = GetDayTime(iDay, myWeek.Wednesday.EndMoning)
			eStartAfternoon, _ = GetDayTime(iDay, myWeek.Wednesday.StartAfternoon)
			eEndAfternoon, _ = GetDayTime(iDay, myWeek.Wednesday.EndAfternoon)

		case time.Thursday:
			eStartMorning, _ = GetDayTime(iDay, myWeek.Thursday.StartMorning)
			eEndMoning, _ = GetDayTime(iDay, myWeek.Thursday.EndMoning)
			eStartAfternoon, _ = GetDayTime(iDay, myWeek.Thursday.StartAfternoon)
			eEndAfternoon, _ = GetDayTime(iDay, myWeek.Thursday.EndAfternoon)

		case time.Friday:
			eStartMorning, _ = GetDayTime(iDay, myWeek.Friday.StartMorning)
			eEndMoning, _ = GetDayTime(iDay, myWeek.Friday.EndMoning)
			eStartAfternoon, _ = GetDayTime(iDay, myWeek.Friday.StartAfternoon)
			eEndAfternoon, _ = GetDayTime(iDay, myWeek.Friday.EndAfternoon)

		case time.Saturday:
			eStartMorning, _ = GetDayTime(iDay, myWeek.Saturday.StartMorning)
			eEndMoning, _ = GetDayTime(iDay, myWeek.Saturday.EndMoning)
			eStartAfternoon, _ = GetDayTime(iDay, myWeek.Saturday.StartAfternoon)
			eEndAfternoon, _ = GetDayTime(iDay, myWeek.Saturday.EndAfternoon)

		}

		eMorningDuration = CalculateDurationByTime(iDay, eStartMorning, eEndMoning)
		eAfternoonDuration = CalculateDurationByTime(iDay, eStartAfternoon, eEndAfternoon)
	}
	return
}
