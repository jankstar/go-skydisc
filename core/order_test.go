package core

import (
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
