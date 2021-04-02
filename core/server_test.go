package core

import (
	"testing"
)

func TestGetExternalIP(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "Test",
			want:    "192.168.0.119",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetExternalIP()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetExternalIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetExternalIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerInit(t *testing.T) {

	type args struct {
		iMode int
		iPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test ServerInit()",
			args: args{
				iMode: 1,
				iPath: "../",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ServerInit(tt.args.iMode, tt.args.iPath); (err != nil) != tt.wantErr {
				t.Errorf("ServerInit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
