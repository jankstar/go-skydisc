package core

import "time"

type DataAssignment struct {
	ID          uint `json:"id" gorm:"primaryKey; autoIncrement"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Start       time.Time    `json:"start"`
	End         time.Time    `json:"end"`
	TimeFix     bool         `json:"time_fix"`
	ResourceFix bool         `json:"resource_fix"`
	OrderRef    uint         `json:"order_ref"`
	Order       DataOrder    `json:"order" gorm:"foreignKey:OrderRef"`
	ResourceRef string       `json:"resource_ref"`
	Resource    DataResource `json:"resource" gorm:"foreignKey:ResourceRef"`
	SectionRef  uint         `json:"section_ref"`
	Section     CatSection   `json:"section" gorm:"foreignKey:SectionRef"`
}

func InitAssignmentDB(iMode int) error {

	//Data
	Server.DB.AutoMigrate(&DataAssignment{})

	if iMode == 1 {
		//Test-Modus - Daten initialisieren
		Server.DB.Where("id <> 0").Delete(&DataAssignment{})

		loadTestDataAssignment()
	}
	return Server.DB.Error
}

func loadTestDataAssignment() {

}

func ProcessAssignment(iAppoint DataAssignment) {
	//save appointment
	var ltDataCapacityCalendar []DataCapacityCalendar
	Server.DB.Create(&iAppoint)

	lDuration := iAppoint.End.Sub(iAppoint.Start)

	//Reduce capacity of the resource
	Server.DB.Where(
		"recource_ref = ? AND date = ?",
		iAppoint.OrderRef, Time2Date(iAppoint.Start)).Find(&ltDataCapacityCalendar)
	var Duration3 struct {
		value time.Duration
		valid bool
	}
	var Duration2 struct {
		value time.Duration
		valid bool
		rest  time.Duration
	}
	var Duration1 struct {
		value time.Duration
		valid bool
		rest  time.Duration
	}

	//first day and morning
	for _, element := range ltDataCapacityCalendar {
		//day
		if element.SectionRef == 3 {
			if Duration3.valid != true {
				if element.DurationRest > lDuration {
					element.DurationRest = element.DurationRest - lDuration
					Duration3.value = element.DurationRest
					Duration3.valid = true
				} else {
					element.DurationRest = 0
					Duration3.value = element.DurationRest
					Duration3.valid = true
				}
			} else {
				element.DurationRest = Duration3.value
			}
		}

		//morning
		if element.SectionRef == 1 {
			if Duration1.valid != true {
				if element.DurationRest > lDuration {
					element.DurationRest = element.DurationRest - lDuration
					Duration1.value = element.DurationRest
					Duration1.valid = true
					Duration1.rest = 0
				} else {
					Duration1.rest = lDuration - element.DurationRest
					element.DurationRest = 0
					Duration1.value = element.DurationRest
					Duration1.valid = true

				}
			} else {
				element.DurationRest = Duration1.value
			}
		}
	}

	//afternoon
	for _, element := range ltDataCapacityCalendar {
		if element.SectionRef == 3 {
			if Duration3.valid != true {
				if Duration1.valid == false {
					//no morning found
					if element.DurationRest > lDuration {
						element.DurationRest = element.DurationRest - lDuration
						Duration3.value = element.DurationRest
						Duration3.valid = true

					} else {

						element.DurationRest = 0
						Duration3.value = element.DurationRest
						Duration3.valid = true

					}
				} else {
					if element.DurationRest > Duration2.rest {
						element.DurationRest = element.DurationRest - Duration2.rest
						Duration3.value = element.DurationRest
						Duration3.valid = true

					} else {

						element.DurationRest = 0
						Duration3.value = element.DurationRest
						Duration3.valid = true

					}
				}
			} else {
				element.DurationRest = Duration3.value
			}
		}
	}

	//save all changes
	Server.DB.Save(&ltDataCapacityCalendar)
}
