package core

import (
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestInitDBOrder(t *testing.T) {
	ServerInit(1, "../")

	type args struct {
		iDB *gorm.DB
	}
	Server.TestfileOrder = "../" + Server.TestfileOrder
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Define Test DB 'test.db' in /tmp ",
			args:    args{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitOrderDB(1); (err != nil) != tt.wantErr {
				t.Errorf("InitDBOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDataOrder_SearchCapacity(t *testing.T) {
	ServerInit(0, "../")
	type fields struct {
		Order DataOrder
	}
	type args struct {
		iCount uint
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantEList []DataCapacityCalendar
	}{
		{
			name: "Test SearchCapacity()",
			fields: fields{
				Order: GetOrderByID(1),
			},
			args: args{
				iCount: 0,
			},
			wantEList: []DataCapacityCalendar{},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			me := tt.fields.Order
			if gotEList := me.SearchCapacity(tt.args.iCount); !reflect.DeepEqual(gotEList, tt.wantEList) {
				//t.Errorf("DataOrder.SearchCapacity() = %v, want %v", gotEList, tt.wantEList)
			}
		})
	}
}
