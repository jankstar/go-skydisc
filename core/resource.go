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
	Qualification    []DataRequirement       `json:"qualification" gorm:"foreignKey:ResourceRef"`
	CalendarRef      string                  `json:"calendar_ref"`
	Calendar         DataWorkingTimeCalendar `json:"calendar" gorm:"foreignKey:CalendarRef"`
	CapacityCalendar []DataCapacityCalendar  `json:"capacity_calendar" gorm:"foreignKey:RecourceRef"`
	AbsencePeriods   []DataAbsencePeriods    `json:"absence_periods" gorm:"foreignKey:RecourceRef"`
	Appointments     []DataAppointment       `json:"appointments" gorm:"foreignKey:RecourceRef"`
}

type DataCapacityCalendar struct {
	ID               uint `json:"id" gorm:"primaryKey; autoIncrement"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	RecourceRef      string           `json:"recource_ref" gorm:"index"`
	Date             string           `json:"date" gorm:"index"`
	QualificationRef string           `json:"qualification_ref"`
	Qualification    CatQualification `json:"qualification" gorm:"foreignKey:QualificationRef"`
	SectionRef       uint             `json:"section_ref" gorm:"index"`
	Section          CatSection       `json:"section" gorm:"foreignKey:SectionRef"`
	StartTime        time.Time        `json:"start_time"`
	EndTime          time.Time        `json:"end_time"`
	DurationTarget   time.Duration    `json:"DurationTarget"`
	DurationRest     time.Duration    `json:"DurationRest"`
}

func InitResourceDB(iMode int) error {

	//Data
	Server.DB.AutoMigrate(&DataResource{})
	Server.DB.AutoMigrate(&DataCapacityCalendar{})

	if iMode == 1 {
		//Test-Modus - Daten initialisieren
		Server.DB.Delete(&DataResource{})
		Server.DB.Delete(&DataCapacityCalendar{})
		Server.DB.Where("resource_ref <> ''").Delete(&DataAbsencePeriods{})
		Server.DB.Where("resource_ref <> ''").Delete(&DataRequirement{})

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

//GetResourceByID returns a resource for an ID
func GetResourceByID(iID string) (eResource DataResource) {

	err := Server.DB.Preload("Qualification").First(&eResource,
		" id = ? AND valid_from <= ? AND ( valid_to > ? OR valid_to = '0001-01-01 00:00:00+00:00' ) ",
		iID, time.Now(), time.Now()).Error
	if err != nil {
		return DataResource{}
	}
	//read CapazityCalendar and AbsencePeriods
	eResource.LoadValidCalender()
	Server.DB.Where(" recource_ref = ? ", eResource.ID).Find(&eResource.CapacityCalendar)
	Server.DB.Where(" recource_ref = ? ", eResource.ID).Find(&eResource.AbsencePeriods)
	return
}

//InitAllResourcesFromNow calculate alle resources from now for @iCount days and delete all older one
func InitAllResourcesFromNow(iCount uint, iDelete bool) {
	var ltResourceID []string
	Server.DB.Model(&DataResource{}).Select("id").Where(
		" valid_from <= ? AND ( valid_to > ? OR valid_to = '0001-01-01 00:00:00+00:00' ) ",
		time.Now(), time.Now()).Find(&ltResourceID)

	for _, lResourceID := range ltResourceID {

		lReource := GetResourceByID(lResourceID)
		if lReource.ID != "" {
			lReource.BuildCapacityCalendarRange(time.Now(), iCount, iDelete)
		}

	}
}

//isInitial() check Resource
func (me *DataResource) isInitial() bool {
	return (me.ID == "")
}

//LoadValidCalender load the valid calender
func (me *DataResource) LoadValidCalender() {
	if me.CalendarRef != "" {
		Server.DB.Where("id = ? AND valid_from < ?", me.CalendarRef, time.Now()).Order("valid_from desc").First(me.Calendar)
	}
}

//BuildCapacityCalendarRange calculates the capacity for the resource from
//a @iStart start date for @iCount days and writes it to the calendar.
func (me *DataResource) BuildCapacityCalendarRange(iStart time.Time, iCount uint, iDelete bool) {

	if iDelete == true {
		//delete old entries
		var lDelCapCal []DataCapacityCalendar
		var lRemainderCapCal []DataCapacityCalendar

		for _, item := range me.CapacityCalendar {
			if Date2Time(item.Date).Before(iStart) == true {
				//elemente of CapacityCalendar is before
				lDelCapCal = append(lDelCapCal, item)
			} else {
				lRemainderCapCal = append(lRemainderCapCal, item)
			}
		}
		Server.DB.Delete(&lDelCapCal)
		me.CapacityCalendar = append([]DataCapacityCalendar{}, lRemainderCapCal...)
	}

	lDate := iStart
	var i uint
	for i = 1; i <= iCount; i++ {
		me.BuildCapacityCalendarByDay(lDate, true)
		lDate = lDate.Add(time.Hour * 24)
	}
}

//BuildCapacityCalendar
func (me *DataResource) BuildCapacityCalendarByDay(iDay time.Time, iReBuild bool) {
	var lDelCapCal []DataCapacityCalendar
	var lRemainderCapCal []DataCapacityCalendar
	justDefine := false
	for _, item := range me.CapacityCalendar {
		lStr := Time2Date(iDay)
		if item.Date == lStr {
			justDefine = true
			if iReBuild == true {
				lDelCapCal = append(lDelCapCal, item)
			} else {
				break
			}
		} else {
			lRemainderCapCal = append(lRemainderCapCal, item)
		}
	}

	//for rebuild == true delete all entries for the date
	if iReBuild == true && len(lDelCapCal) > 0 {
		Server.DB.Delete(&lDelCapCal)
		me.CapacityCalendar = append([]DataCapacityCalendar{}, lRemainderCapCal...)
	}

	if justDefine == false || iReBuild == true {
		//define new capacity
		lMorningDuration, lAfternoonDuration,
			lStartMorning, lEndMorning, lStartAfternoon, lEndAfternoon := me.Calendar.GetDurationByDay(iDay)
		for _, itemQuali := range me.Qualification {
			newCapacityMorning := &DataCapacityCalendar{
				RecourceRef:      me.ID,
				Date:             Time2Date(iDay),
				QualificationRef: itemQuali.QualificationRef,
				SectionRef:       1,
				StartTime:        lStartMorning,
				EndTime:          lEndMorning,
				DurationTarget:   lMorningDuration,
				DurationRest:     lMorningDuration,
			}
			newCapacityAfternoon := &DataCapacityCalendar{
				RecourceRef:      me.ID,
				Date:             Time2Date(iDay),
				QualificationRef: itemQuali.QualificationRef,
				SectionRef:       2,
				StartTime:        lStartAfternoon,
				EndTime:          lEndAfternoon,
				DurationTarget:   lAfternoonDuration,
				DurationRest:     lAfternoonDuration,
			}
			newCapacityDay := &DataCapacityCalendar{
				RecourceRef:      me.ID,
				Date:             Time2Date(iDay),
				QualificationRef: itemQuali.QualificationRef,
				SectionRef:       3,
				StartTime:        lStartMorning,
				EndTime:          lEndAfternoon,
				DurationTarget:   lMorningDuration + lAfternoonDuration,
				DurationRest:     lMorningDuration + lAfternoonDuration,
			}

			//create Date if duration not null
			if newCapacityDay.DurationTarget != 0 {
				Server.DB.Create(&newCapacityDay)
				if newCapacityMorning.DurationTarget != 0 {
					Server.DB.Create(&newCapacityMorning)
				}
				if newCapacityAfternoon.DurationTarget != 0 {
					Server.DB.Create(&newCapacityAfternoon)
				}
				//Consider absences
				me.MergeAbsencePeriodsByDay(iDay, newCapacityDay.StartTime, newCapacityDay.EndTime)
			}
		}
	}
}

//MergeAbsencePeriodsByDay consider absences by day
func (me *DataResource) MergeAbsencePeriodsByDay(iDay time.Time, iStart time.Time, iEnd time.Time) {
	for _, element := range me.AbsencePeriods {
		//
		absenceStart_LE_Start := element.Start.Before(iStart) == true || element.Start.Equal(iStart) == true
		absenceEnd_GE_End := iEnd.Before(element.End) == true || element.End.Equal(iEnd) == true

		if absenceStart_LE_Start && absenceEnd_GE_End {
			// kompletter Bereich
			newAppointment := DataAppointment{
				Start:       iStart,
				End:         iEnd,
				TimeFix:     true,
				ResourceFix: true,
				OrderRef:    0,
				ResourceRef: me.ID,
				SectionRef:  3,
			}
			ProcessAppointment(newAppointment)
			me.Appointments = append(me.Appointments, newAppointment)

		}

	}

}
