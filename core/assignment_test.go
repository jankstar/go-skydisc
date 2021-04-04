package core

import (
	"testing"
	"time"
)

func TestProcessAssignment(t *testing.T) {
	ServerInit(0, "../")
	lStart, _ := time.Parse(time.RFC3339, "2020-04-05T08:00:00Z")
	lEnd, _ := time.Parse(time.RFC3339, "2020-04-05T09:00:00Z")
	type args struct {
		iAppoint DataAssignment
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test Assignment",
			args: args{
				iAppoint: DataAssignment{
					Start:       lStart,
					End:         lEnd,
					TimeFix:     true,
					ResourceFix: false,
					OrderRef:    1,
					Order:       DataOrder{},
					ResourceRef: "MUELLERH",
					SectionRef:  1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ProcessAssignment(tt.args.iAppoint)
		})
	}
}
