package core

import "time"

type DataAppointment struct {
	ID          uint `json:"id" gorm:"primaryKey; autoIncrement"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	OrderRef    uint         `json:"order_ref"`
	Order       DataOrder    `json:"order" gorm:"foreignKey:OrderRef"`
	ResourceRef string       `json:"resource_ref"`
	Resource    DataResource `json:"resource" gorm:"foreignKey:ResourceRef"`
	SectionRef  uint         `json:"section_ref"`
	Section     CatSection   `json:"section" gorm:"foreignKey:SectionRef"`
}

func InitAppointmentDB(iMode int) error {

	//Data
	Server.DB.AutoMigrate(&DataAppointment{})

	if iMode == 1 {
		//Test-Modus - Daten initialisieren
		Server.DB.Where("id <> ''").Delete(&DataAppointment{})

		loadTestDataAppointment()
	}
	return Server.DB.Error
}

func loadTestDataAppointment() {

}
