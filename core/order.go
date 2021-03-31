package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type DataProjekt struct {
	Number        string `json:"project_number" gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string    `json:"project_name"`
	EarliestStart time.Time `json:"earliest_start"`
	LatestEnd     time.Time `json:"latest_end"`
}

//DatOrder - define data Order entity
type DataOrder struct {
	ID               uint `json:"id" gorm:"primaryKey; autoIncrement"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Description      string            `json:"description"`
	OrderTypeRef     string            `json:"order_type_ref"`
	OrderType        CatOrderClass     `json:"order_type" gorm:"foreignKey:OrderTypeRef"`
	EarliestStart    time.Time         `json:"earliest_start"`
	LatestEnd        time.Time         `json:"latest_end"`
	Distress         bool              `json:"distress"`
	Priority         int               `json:"priority"`
	ProjectRef       string            `json:"project_ref"`
	Project          DataProjekt       `json:"project" gorm:"foreignKey:ProjectRef"`
	Duration         time.Duration     `json:"duration"`
	Location         TLocation         `json:"location" gorm:"embedded"`
	ContactPerson    string            `json:"contact_person"`
	Client           string            `json:"client"`
	Requirement      []DataRequirement `json:"requirement" gorm:"foreignKey:OrderRefer"`
	PredecessorRef   uint              `json:"predecessor_ref"`
	ServiceAreaRefer string            `json:"service_area_refer"`
	ServiceArea      DataServiceArea   `json:"service_area" gorm:"foreignKey:ServiceAreaRefer"`
}

//InitOrderDB(iDB *gorm.DB) error
// initiates the DB tables for the job and all the
// dependent tables
func InitOrderDB(iMode int) error {

	//Data
	Server.DB.AutoMigrate(&DataOrder{})
	Server.DB.AutoMigrate(&DataProjekt{})

	if iMode == 1 {
		//Test-Modus - Daten initialisieren
		Server.DB.Where("id <> ''").Delete(&DataOrder{})
		Server.DB.Where("number <> ''").Delete(&DataProjekt{})

		loadTestDataOrder()
	}

	return Server.DB.Error
}

//*********************************

func loadTestDataOrder() {
	var test struct {
		OrderList   []DataOrder   `json:"order_list"`
		ProjectList []DataProjekt `json:"project_list"`
	}
	fmt.Println("loadTestData: order.json")
	data, err := ioutil.ReadFile(Server.Path + Server.TestfileOrder)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	err = json.Unmarshal(data, &test)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	Server.DB.Save(&test.OrderList)
	Server.DB.Save(&test.ProjectList)

	fmt.Println("loadTestData: in DB verbucht")
}

//************************************
//Save(iDB *gorm.DB) (err error)
// save der Order itself
func (me *DataOrder) Save() (err error) {
	Server.DB.Save(&me)
	return Server.DB.Error
}

//Check() (err error)
func (me *DataOrder) Check() (err error) {
	return nil
}
