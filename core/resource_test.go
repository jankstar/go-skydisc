package core

import (
	"testing"
	"time"
)

func TestInitResourceDB(t *testing.T) {
	ServerInit(1, "../")

	type args struct {
		iMode int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test Init Resource DB",
			args: args{
				iMode: 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitResourceDB(tt.args.iMode); (err != nil) != tt.wantErr {
				t.Errorf("InitResourceDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDataResource_BuildCapacityCalendarAtDay(t *testing.T) {
	ServerInit(0, "../")
	lTestResource := GetResourceByID("SCHMIDTK")
	type fields struct {
		ME DataResource
	}
	type args struct {
		iDate    time.Time
		iReBuild bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test BuildCapacityCalendar for Object resource",
			fields: fields{
				ME: lTestResource,
			},
			args: args{
				iDate:    time.Now(),
				iReBuild: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			me := tt.fields.ME
			me.BuildCapacityCalendarByDay(tt.args.iDate, tt.args.iReBuild)
		})
	}
}

func TestDataResource_BuildCapacityCalendarRange(t *testing.T) {
	ServerInit(0, "../")
	lTestResource := GetResourceByID("MUELLERH")
	type fields struct {
		ME DataResource
	}
	type args struct {
		iStart  time.Time
		iCount  uint
		iDelete bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test BuildCapacityCalendarRange",
			fields: fields{
				ME: lTestResource,
			},
			args: args{
				iStart:  time.Now(),
				iCount:  14,
				iDelete: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			me := tt.fields.ME
			me.BuildCapacityCalendarRange(tt.args.iStart, tt.args.iCount, tt.args.iDelete)
		})
	}
}

func TestInitAllResourcesFromNow(t *testing.T) {
	ServerInit(0, "../")
	type args struct {
		iCount  uint
		iDelete bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test InitAllResourcesFromNow()",
			args: args{
				iCount:  14,
				iDelete: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitAllResourcesFromNow(tt.args.iCount, tt.args.iDelete)
		})
	}
}
