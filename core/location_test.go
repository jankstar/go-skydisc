package core

import (
	"testing"
	"time"
)

func TestDatOrder_GetGeoLocationFromBing(t *testing.T) {

	ServerInit(1, "../")

	type fields struct {
		Location TLocation
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Test bing ge0-location",
			fields: fields{

				Location: TLocation{
					CountryCode:   "DE",
					PostCode:      "10178",
					RegionCode:    "",
					Town:          "Berlin",
					Street:        "Alexanderplatz",
					StreetNumber:  "1",
					BuildingName:  "",
					BuildingFloor: "",
					BuildingRoom:  "",
					GeoAltitude:   0,
					GeoLatitude:   0,
					GeoLongitude:  0,
					GeoTimestamp:  time.Time{},
					GeoServices:   "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			me := tt.fields.Location
			me.GetGeoLocationFromBing()
			if me.GeoLatitude != 52.521915 ||
				me.GeoLongitude != 13.415063 {
				t.Errorf("GetGeoLocationFromBing() Latitude 52.521915 != %v, Longitude 13.415063 != %v",
					me.GeoLatitude, me.GeoLongitude)
			}
		})
	}
}
