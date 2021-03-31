package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

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

//DatLocationBuffer is the DB Buffer table for location and geo-location
// so we not need to recall bing REST
type DataLocationBuffer struct {
	ID           uint `json:"id" gorm:"primaryKey; autoIncrement"`
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

//InitLocationDB(iDB *gorm.DB) error
// initiates the DB tables for the job and all the
// dependent tables
func InitLocationDB(iMode int) error {

	//Data
	Server.DB.AutoMigrate(&DataLocationBuffer{})

	if iMode == 1 {
		//Test-Modus - Daten initialisieren
		Server.DB.Where("id <> ''").Delete(&DataLocationBuffer{})

		loadTestDataLocation()
	}

	return Server.DB.Error
}

//*********************************

func loadTestDataLocation() {

}

//GetGeoLocationFromBing call Bing REST for GeoLocation
func (me *TLocation) GetGeoLocationFromBing() {
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

	if me.CountryCode == "" ||
		me.PostCode == "" ||
		me.Town == "" ||
		me.Street == "" ||
		Server.BingApiKey == "" ||
		me.GeoLatitude != 0 ||
		me.GeoLongitude != 0 {
		return //we need location data and geo-location is initial
	}

	//first we search in buffer
	var lLocSearch DataLocationBuffer
	Server.DB.Where(&DataLocationBuffer{
		CountryCode:  me.CountryCode,
		PostCode:     me.PostCode,
		RegionCode:   me.RegionCode,
		Town:         me.Town,
		Street:       me.Street,
		StreetNumber: me.StreetNumber,
		BuildingName: me.BuildingName,
	}).First(&lLocSearch)
	if Server.DB.Error == nil &&
		lLocSearch.GeoLatitude != 0 &&
		lLocSearch.GeoLongitude != 0 {
		me.GeoAltitude = lLocSearch.GeoAltitude
		me.GeoLatitude = lLocSearch.GeoLatitude
		me.GeoLongitude = lLocSearch.GeoLongitude
		me.TimeZone = lLocSearch.TimeZone
		return
	}

	//** call BING Maps for ge-location
	{
		lURL := fmt.Sprintf(
			Server.BingURLLocation,
			me.CountryCode,
			me.PostCode,
			me.Town,
			me.Street+"%20"+me.StreetNumber,
			1, //we seache for only one result
			Server.BingApiKey)
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
					me.GeoLatitude = dstLocation.ResourceSets[0].Resources[0].Point.Coordinates[0]
					me.GeoLongitude = dstLocation.ResourceSets[0].Resources[0].Point.Coordinates[1]
					me.GeoTimestamp = time.Now()
					me.GeoServices = "bing"
				}

			}
		}
	}

	//call BING maps for timezone
	{
		lURL := fmt.Sprintf(
			Server.BingURLTimezone,
			me.PostCode+"%20"+me.Town+"%20"+me.CountryCode,
			Server.BingApiKey)
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
			me.TimeZone = dstTimezone.ResourceSets[0].Resources[0].TimeZoneAtLocation[0].TimeZone[0].Abbreviation
			fmt.Println(me.TimeZone)
		}
	}

	//save location to buffer
	var lLoc = DataLocationBuffer{
		CountryCode:  me.CountryCode,
		PostCode:     me.PostCode,
		RegionCode:   me.RegionCode,
		Town:         me.Town,
		Street:       me.Street,
		StreetNumber: me.StreetNumber,
		BuildingName: me.BuildingName,
		TimeZone:     me.TimeZone,
		GeoAltitude:  me.GeoAltitude,
		GeoLatitude:  me.GeoLatitude,
		GeoLongitude: me.GeoLongitude,
		GeoTimestamp: me.GeoTimestamp,
		GeoServices:  me.GeoServices,
	}
	Server.DB.Create(&lLoc)

}
