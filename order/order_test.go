package order

import (
	"testing"
	"time"

	"github.com/jankstar/go-skydisc/catalog"
	"github.com/jankstar/go-skydisc/lib"
	"gorm.io/gorm"
)

func TestInitDBOrder(t *testing.T) {
	lib.ServerInit(1, "../")
	type args struct {
		iDB *gorm.DB
	}
	lib.Server.TestfileOrder = "../" + lib.Server.TestfileOrder
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Define Test DB 'test.db' in /tmp ",
			args:    args{lib.Server.DB},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitOrderDB(tt.args.iDB, 1); (err != nil) != tt.wantErr {
				t.Errorf("InitDBOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatOrder_GetGeoLocationFromBing(t *testing.T) {
	type fields struct {
		ID            uint
		CreatedAt     time.Time
		UpdatedAt     time.Time
		Description   string
		OrderType     catalog.CatOrderClass
		EarliestStart time.Time
		EatestEnd     time.Time
		Duration      time.Duration
		Location      TLocation
		ContactPerson string
		Client        string
		Requirement   []DatOrderRequirement
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Test bing ge0-location",
			fields: fields{
				ID:          1,
				CreatedAt:   time.Time{},
				UpdatedAt:   time.Time{},
				Description: "Testorder",
				OrderType: catalog.CatOrderClass{
					Class: "KLR",
					Name:  "Kleinreparatur",
				},
				EarliestStart: time.Time{},
				EatestEnd:     time.Time{},
				Duration:      6000,
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
				ContactPerson: "",
				Client:        "",
				Requirement:   []DatOrderRequirement{},
			},
		},
	}
	lib.ServerInit(1, "../")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			me := &DatOrder{
				ID:            tt.fields.ID,
				CreatedAt:     tt.fields.CreatedAt,
				UpdatedAt:     tt.fields.UpdatedAt,
				Description:   tt.fields.Description,
				OrderType:     tt.fields.OrderType,
				EarliestStart: tt.fields.EarliestStart,
				EatestEnd:     tt.fields.EatestEnd,
				Duration:      tt.fields.Duration,
				Location:      tt.fields.Location,
				ContactPerson: tt.fields.ContactPerson,
				Client:        tt.fields.Client,
				Requirement:   tt.fields.Requirement,
			}
			me.GetGeoLocationFromBing(lib.Server.DB)
			if me.Location.GeoLatitude != 52.521915 ||
				me.Location.GeoLongitude != 13.415063 {
				t.Errorf("GetGeoLocationFromBing() Latitude 52.521915 != %v, Longitude 13.415063 != %v",
					me.Location.GeoLatitude, me.Location.GeoLongitude)
			}
		})
	}
}
