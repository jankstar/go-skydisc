package order

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jankstar/go-skydisc/catalog"
	"github.com/jankstar/go-skydisc/lib"
)

//DatLocationBuffer is the DB Buffer table for location and geo-location
// so we not need to recall bing REST
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
	TimeZone     string    `json:"timeZone"`
	GeoAltitude  float64   `json:"geo_altitude"`
	GeoLatitude  float64   `json:"geo_latitude"`
	GeoLongitude float64   `json:"geo_longitude"`
	GeoTimestamp time.Time `json:"geo_timestamp"`
	GeoServices  string    `json:"GeoServices"`
}

//DatOrderRequirement are the order requirements as table
type DatOrderRequirement struct {
	ID             uint `json:"id" gorm:"primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Trade          catalog.CatTrade         `json:"trade" gorm:"embedded"`
	Qualification  catalog.CatQualification `json:"qualification" gorm:"embedded;embeddedPrefix:cat_"`
	NumOfResources int                      `json:"num_of_resources"`
	OrderRefer     string
}

//TLocation is the sub-struct for order location
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
	TimeZone      string    `json:"time_zone"`
	GeoAltitude   float64   `json:"geo_altitude"`
	GeoLatitude   float64   `json:"geo_latitude"`
	GeoLongitude  float64   `json:"geo_longitude"`
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
func InitOrderDB(iMode int) error {

	//Data
	lib.Server.DB.AutoMigrate(&DatOrder{})
	lib.Server.DB.AutoMigrate(&DatLocationBuffer{})
	lib.Server.DB.AutoMigrate(&DatOrderRequirement{})

	if iMode == 1 {
		//Test-Modus - Daten initialisieren
		lib.Server.DB.Where("id <> ''").Delete(&DatOrder{})
		lib.Server.DB.Where("id <> ''").Delete(&DatLocationBuffer{})
		lib.Server.DB.Where("id <> '").Delete(&DatOrderRequirement{})

		loadTestData()
	}

	return lib.Server.DB.Error
}

//*********************************

func loadTestData() {
	var test struct {
		OrderList []DatOrder `json:"order_list"`
	}
	fmt.Println("loadTestData: order.json")
	data, err := ioutil.ReadFile(lib.Server.TestfileOrder)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	err = json.Unmarshal(data, &test)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	lib.Server.DB.Save(&test.OrderList)

	fmt.Println("loadTestData: in DB verbucht")
}

//************************************
//Save(iDB *gorm.DB) (err error)
// save der Order itself
func (me *DatOrder) Save() (err error) {
	lib.Server.DB.Save(&me)
	return lib.Server.DB.Error
}

//Check() (err error)
func (me *DatOrder) Check() (err error) {
	return nil
}

//GetGeoLocationFromBing call Bing REST for GeoLocation
func (me *DatOrder) GetGeoLocationFromBing() {
	var dstLocation struct {
		ResourceSets []struct {
			EstimatedTotal int `json:"estimatedTotal"`
			Resources      []struct {
				Point struct {
					Coordinates []float64 `json:"coordinates"`
				} `json:"point"`
			} `json:"resources"`
		} `json:"resourceSets"`
	}
	var dstTimezone struct {
		ResourceSets []struct {
			EstimatedTotal int `json:"estimatedTotal"`
			Resources      []struct {
				TimeZoneAtLocation []struct {
					TimeZone []struct {
						Abbreviation string `json:"abbreviation"`
					} `json:"timeZone"`
				} `json:"timeZoneAtLocation"`
			} `json:"resources"`
		} `json:"resourceSets"`
	}

	if me.Location.CountryCode == "" ||
		me.Location.PostCode == "" ||
		me.Location.Town == "" ||
		me.Location.Street == "" ||
		lib.Server.BingApiKey == "" ||
		me.Location.GeoLatitude != 0 ||
		me.Location.GeoLongitude != 0 {
		return //we need location data and geo-location is initial
	}

	//first we search in buffer
	var lLocSearch DatLocationBuffer
	lib.Server.DB.Where(&DatLocationBuffer{
		CountryCode:  me.Location.CountryCode,
		PostCode:     me.Location.PostCode,
		RegionCode:   me.Location.RegionCode,
		Town:         me.Location.Town,
		Street:       me.Location.Street,
		StreetNumber: me.Location.StreetNumber,
		BuildingName: me.Location.BuildingName,
	}).First(&lLocSearch)
	if lib.Server.DB.Error == nil &&
		lLocSearch.GeoLatitude != 0 &&
		lLocSearch.GeoLongitude != 0 {
		me.Location.GeoAltitude = lLocSearch.GeoAltitude
		me.Location.GeoLatitude = lLocSearch.GeoLatitude
		me.Location.GeoLongitude = lLocSearch.GeoLongitude
		me.Location.TimeZone = lLocSearch.TimeZone
		return
	}

	//** call BING Maps for ge-location
	{
		lURL := fmt.Sprintf(
			lib.Server.BingURLLocation,
			me.Location.CountryCode,
			me.Location.PostCode,
			me.Location.Town,
			me.Location.Street+"%20"+me.Location.StreetNumber,
			1, //we seache for only one result
			lib.Server.BingApiKey)
		fmt.Println(lURL)

		res, err1 := http.Get(lURL)
		if err1 != nil || res.StatusCode != 200 {
			return
		}
		defer res.Body.Close()

		dec := json.NewDecoder(res.Body)
		//dec.DisallowUnknownFields()

		err2 := dec.Decode(&dstLocation)
		if err2 != nil {
			return
		}
		if len(dstLocation.ResourceSets) > 0 {
			if len(dstLocation.ResourceSets[0].Resources) > 0 {
				if len(dstLocation.ResourceSets[0].Resources[0].Point.Coordinates) == 2 {
					fmt.Printf("x:%f y:%f \n",
						dstLocation.ResourceSets[0].Resources[0].Point.Coordinates[0],
						dstLocation.ResourceSets[0].Resources[0].Point.Coordinates[1])
					me.Location.GeoLatitude = dstLocation.ResourceSets[0].Resources[0].Point.Coordinates[0]
					me.Location.GeoLongitude = dstLocation.ResourceSets[0].Resources[0].Point.Coordinates[1]
					me.Location.GeoTimestamp = time.Now()
					me.Location.GeoServices = "bing"
				}

			}
		}
	}

	//call BING maps for timezone
	{
		lURL := fmt.Sprintf(
			lib.Server.BingURLTimezone,
			me.Location.PostCode+"%20"+me.Location.Town+"%20"+me.Location.CountryCode,
			lib.Server.BingApiKey)
		fmt.Println(lURL)

		res, err1 := http.Get(lURL)
		if err1 != nil || res.StatusCode != 200 {
			return
		}
		defer res.Body.Close()

		dec := json.NewDecoder(res.Body)
		//dec.DisallowUnknownFields()

		err2 := dec.Decode(&dstTimezone)
		if err2 != nil {
			return
		}
		if len(dstTimezone.ResourceSets) > 0 &&
			len(dstTimezone.ResourceSets[0].Resources) > 0 &&
			len(dstTimezone.ResourceSets[0].Resources[0].TimeZoneAtLocation) > 0 &&
			len(dstTimezone.ResourceSets[0].Resources[0].TimeZoneAtLocation[0].TimeZone) > 0 {
			me.Location.TimeZone = dstTimezone.ResourceSets[0].Resources[0].TimeZoneAtLocation[0].TimeZone[0].Abbreviation
			fmt.Println(me.Location.TimeZone)
		}
	}

	//save location to buffer
	var lLoc = DatLocationBuffer{
		CountryCode:  me.Location.CountryCode,
		PostCode:     me.Location.PostCode,
		RegionCode:   me.Location.RegionCode,
		Town:         me.Location.Town,
		Street:       me.Location.Street,
		StreetNumber: me.Location.StreetNumber,
		BuildingName: me.Location.BuildingName,
		TimeZone:     me.Location.TimeZone,
		GeoAltitude:  me.Location.GeoAltitude,
		GeoLatitude:  me.Location.GeoLatitude,
		GeoLongitude: me.Location.GeoLongitude,
		GeoTimestamp: me.Location.GeoTimestamp,
		GeoServices:  me.Location.GeoServices,
	}
	lib.Server.DB.Create(&lLoc)

}
