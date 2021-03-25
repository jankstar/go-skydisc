package order

import (
	"time"

	"github.com/jankstar/go-skydisc/catalog"
	"gorm.io/gorm"
)

type DatLocationBuffer struct {
	ID           uint `json:"id" gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
	ID             uint `json:"id" gorm:"primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Trade          catalog.CatTrade         `json:"trade" gorm:"embedded"`
	Qualification  catalog.CatQualification `json:"qualification" gorm:"embedded;embeddedPrefix:cat_"`
	NumOfResources int                      `json:"num_of_resources"`
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
	ID            uint `json:"id" gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Description   string                `json:"description"`
	OrderType     catalog.CatOrderClass `json:"order_type" gorm:"embedded;embeddedPrefix:cat_"`
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
func InitOrderDB(iDB *gorm.DB, iMode int) error {

	//Data
	iDB.AutoMigrate(&DatOrder{})
	iDB.AutoMigrate(&DatLocationBuffer{})
	iDB.AutoMigrate(&DatOrderRequirement{})

	if iMode == 1 {
		//Test-Modus - Daten initialisieren
		iDB.Where("id <> ''").Delete(&DatOrder{})
		iDB.Where("id <> ''").Delete(&DatLocationBuffer{})
		iDB.Where("id <> '").Delete(&DatOrderRequirement{})

		loadTestData(iDB)
	}

	return iDB.Error
}

func loadTestData(iDB *gorm.DB) {
	// var test struct {
	// 	TradeList         []catalog.CatTrade         `json:"trade_list"`
	// 	QualificationList []catalog.CatQualification `json:"qualification_list"`
	// 	OrderClassList    []catalog.CatOrderClass    `json:"class_list"`
	// }
	// fmt.Println("loadTestData: test.json")
	// data, err := ioutil.ReadFile(lib.Server.Testfile)
	// if err != nil {
	// 	fmt.Println("error:", err)
	// 	return
	// }
	// err = json.Unmarshal(data, &test)
	// if err != nil {
	// 	fmt.Println("error:", err)
	// 	return
	// }
	// iDB.Save(&test.OrderClassList)
	// iDB.Save(&test.QualificationList)
	// iDB.Save(&test.TradeList)
	// fmt.Println("loadTestData: in DB verbucht")
}

//Save(iDB *gorm.DB) (err error)
// save der Order itself
func (me *DatOrder) Save(iDB *gorm.DB) (err error) {
	iDB.Save(&me).Commit()
	return iDB.Error
}

//Check() (err error)
func (me *DatOrder) Check() (err error) {
	return nil
}
