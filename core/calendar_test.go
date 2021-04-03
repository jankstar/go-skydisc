package core

import (
	"reflect"
	"testing"
	"time"
)

func TestTime2Date(t *testing.T) {
	lTime, _ := time.Parse(time.RFC3339, "2020-03-23T08:10:12Z")
	type args struct {
		iTime time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Time2Date",
			args: args{
				iTime: lTime,
			},
			want: "2020-03-23",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Time2Date(tt.args.iTime); got != tt.want {
				t.Errorf("Time2Date() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTime2Time(t *testing.T) {
	lTime, _ := time.Parse(time.RFC3339, "2020-03-23T08:10:12Z")
	type args struct {
		iTime time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Time2Time",
			args: args{
				iTime: lTime,
			},
			want: "08:10:12",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Time2Time(tt.args.iTime); got != tt.want {
				t.Errorf("Time2Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDuration(t *testing.T) {
	type args struct {
		iDay   time.Time
		iStart string
		iEnd   string
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "Test CalculateDurationByStr()",
			args: args{
				iDay:   time.Now(),
				iStart: "08:10:00",
				iEnd:   "09:10:00",
			},
			want: time.Hour,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateDurationByStr(tt.args.iDay, tt.args.iStart, tt.args.iEnd); got != tt.want {
				t.Errorf("CalculateDurationByStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDayTime(t *testing.T) {
	lTime, _ := time.Parse(time.RFC3339, "2020-03-23T08:10:12Z")
	lTimeWant, _ := time.Parse(time.RFC3339, "2020-03-23T10:10:05Z")

	type args struct {
		iDay     time.Time
		iTimeStr string
	}
	tests := []struct {
		name      string
		args      args
		wantETime time.Time
		wantErr   bool
	}{
		{
			name: "Test GetDayTime()",
			args: args{
				iDay:     lTime,
				iTimeStr: "10:10:05",
			},
			wantETime: lTimeWant,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotETime, err := GetDayTime(tt.args.iDay, tt.args.iTimeStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDayTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotETime, tt.wantETime) {
				t.Errorf("GetDayTime() = %v, want %v", gotETime, tt.wantETime)
			}
		})
	}
}
