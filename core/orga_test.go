package core

import (
	"testing"
)

func TestInitOrgaDB(t *testing.T) {
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
			name: "Test InitOrgaDB()",
			args: args{
				iMode: 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitOrgaDB(tt.args.iMode); (err != nil) != tt.wantErr {
				t.Errorf("InitOrgaDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
