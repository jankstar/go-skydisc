package order

import (
	"time"

	"gorm.io/gorm"
)

//CatOrderClass - define Order class as customization
type CatOrderClass struct {
	gorm.Model
	Class string `json:"class"`
}

//CatTrade - define Order trade as customization
type CatTrade struct {
	gorm.Model
	Trade string `json:"trade"`
}

//CatQualification - define Qualification as customization
type CatQualification struct {
	gorm.Model
	Qualification string `json:"qualification"`
}

type DatLocationBuffer struct {
	gorm.Model
	CountryCode  string    `json:"country_code"`
	PostCode     string    `json:"post_code"`
	RegionCode   string    `json:"region_code"`
	Town         string    `json:"town"`
	Street       string    `json:"street"`
	StreetNumber string    `json:"street_number"`
	BuildingName string    `json:"building_name"`
	GeoAltitude  int64     `json:"geo_altitude"`
	GeoLatitude  int64     `json:"geo_latitude"`
	GeoLongitude int64     `json:"geo_longitude"`
	GeoTimestamp time.Time `json:"geo_timestamp"`
	GeoServices  string    `json:"GeoServices"`
}

type DatOrderRequirement struct {
	gorm.Model
	Trade          CatTrade         `json:"trade" gorm:"embedded"`
	Qualification  CatQualification `json:"qualification" gorm:"embedded;embeddedPrefix:cat_"`
	NumOfResources int              `json:"num_of_resources"`
	OrderRefer     string
}

type TLocation struct {
	CountryCode   string    `json:"country_code"`
	PostCode      string    `json:"post_code"`
	RegionCode    string    `json:"region_code"`
	Town          string    `json:"town"`
	Street        string    `json:"street"`
	StreetNumber  string    `json:"street_number"`
	BuildingName  string    `json:"building_name"`
	BuildingFloor string    `json:"building_floor"`
	BuildingRoom  string    `json:"building_room"`
	GeoAltitude   int64     `json:"geo_altitude"`
	GeoLatitude   int64     `json:"geo_latitude"`
	GeoLongitude  int64     `json:"geo_longitude"`
	GeoTimestamp  time.Time `json:"geo_timestamp"`
	GeoServices   string    `json:"GeoServices"`
}

//DatOrder - define data Order entity
type DatOrder struct {
	gorm.Model
	Description   string                `json:"description"`
	OrderType     CatOrderClass         `json:"order_type" gorm:"embedded;embeddedPrefix:cat_"`
	EarliestStart time.Time             `json:"earliest_start"`
	EatestEnd     time.Time             `json:"latest_end"`
	Duration      time.Duration         `json:"duration"`
	Location      TLocation             `json:"location" gorm:"embedded"`
	ContactPerson string                `json:"contact person"`
	Client        string                `json:"client"`
	Requirement   []DatOrderRequirement `json:"requirement" gorm:"foreignKey:OrderRefer"`
}

//InitOrderDB(iDB *gorm.DB) error
// initiates the DB tables for the job and all the
// dependent tables
func InitOrderDB(iDB *gorm.DB) error {

	//Data
	iDB.AutoMigrate(&DatOrder{})
	iDB.AutoMigrate(&DatLocationBuffer{})
	iDB.AutoMigrate(&DatOrderRequirement{})

	//Catalogs
	iDB.AutoMigrate(&CatOrderClass{})
	iDB.AutoMigrate(&CatTrade{})
	iDB.AutoMigrate(&CatQualification{})
	return iDB.Error
}

//Save(iDB *gorm.DB) (err error)
// save der Order itself
func (me *DatOrder) Save(iDB *gorm.DB) (err error) {
	iDB.Save(&me).Commit()
	return iDB.Error
}
