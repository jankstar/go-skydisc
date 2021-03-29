package lib

import "testing"

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
