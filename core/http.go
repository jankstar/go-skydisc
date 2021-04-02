package core

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HTTPRouter(iRouter *gin.Engine) {
	iRouter.GET("/", HTTPindexFunc)
	iRouter.GET("/api/:table", HTTPapiFunc)
}

//returns context index.html
func HTTPindexFunc(iCon *gin.Context) {
	iCon.HTML(http.StatusOK, "index.html", gin.H{
		"title": "go-skaydisc",
	})
}

func ParseQuery(iTable string, iSearch string) func(db *gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {

		myDB := db
		m, _ := url.ParseQuery(iSearch)

		if m["q"] != nil && m["q"][0] != "" {
			//query: where Condition
			if mq := strings.Split(m["q"][0], ":"); len(mq) >= 2 {

				mq[1] = strings.ReplaceAll(mq[1], "*", "%")
				//myDB = myDB.Where(" ? like ? ", mq[0], mq[1])
				//switch mq[0] {
				//case "body", "subject", "status", "amount", "sender_name", "recipient_name", "category":
				myDB = myDB.Where(" "+mq[0]+" like ? ", mq[1])
				//}
			}
		}

		//No deleted entries
		if iTable == "record" {
			myDB = myDB.Where("deleted = ? ", "0001-01-01 00:00:00+00:00")
		}

		if m["sort"] != nil && m["sort"][0] != "" {
			//order by
			myDB = myDB.Order(m["sort"][0])
		}

		if m["rows"] != nil && m["rows"][0] != "" {
			//Limit
			lLimit, err := strconv.Atoi(m["rows"][0])
			if err == nil {
				myDB = myDB.Limit(lLimit)
			}
		}

		return myDB
	}
}

//returns api
func HTTPapiFunc(iCon *gin.Context) {
	var lTable string = iCon.Param("table")
	var lSearch string = iCon.Request.URL.RawQuery

	if lTable == "" || lSearch == "" {
		iCon.JSON(http.StatusBadRequest, gin.H{
			"data": "",
		})
		return
	}

	if lTable == "DataOrder" {
		var ltDataOrder []DataOrder
		Server.DB.Scopes(ParseQuery(lTable, lSearch)).Preload("OrderType").Preload("Project").Preload("ServiceArea").Find(&ltDataOrder)
		if ltDataOrder != nil {
			iCon.JSON(http.StatusOK, gin.H{
				"data": ltDataOrder,
				"visible_columns": []string{"id", "description",
					"order_type_ref", "earliest_start", "latest_end", "distress", "priority",
					"project.project_number", "project.project_name",
					"location.country_code", "location.post_code", "location.town", "location.street", "location.street_number"},
			})
		} else {
			iCon.JSON(http.StatusBadRequest, gin.H{
				"data":            "",
				"visible_columns": "",
			})
		}

		return
	}
}
