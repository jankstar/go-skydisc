package core

import (
	"time"
)

//DatOrderRequirement are the order requirements as table
type DataRequirement struct {
	ID               uint `json:"id" gorm:"primaryKey; autoIncrement"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	TradeRef         string           `json:"trade_ref"`
	Trade            CatTrade         `json:"trade" gorm:"foreignKey:TradeRef"`
	QualificationRef string           `json:"qualification_ref"`
	Qualification    CatQualification `json:"qualification" gorm:"foreignKey:QualificationRef"`
	NumOfResources   int              `json:"num_of_resources"`
	ServiceAreaRef   string
	ResourceRef      string
}

//InitRequirementDB(iDB *gorm.DB) error
// initiates the DB tables for the job and all the
// dependent tables
func InitRequirementDB(iMode int) error {

	//Data
	Server.DB.AutoMigrate(&DataRequirement{})

	if iMode == 1 {
		//Test-Modus - Daten initialisieren
		//lib.Server.DB.Where("id <> '' AND order_ref <> ''").Delete(&DataRequirement{})

		loadTestDataRequirement()
	}

	return Server.DB.Error
}

//*********************************

func loadTestDataRequirement() {

}
