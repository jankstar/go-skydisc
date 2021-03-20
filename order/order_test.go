package order

import (
	"testing"

	"github.com/jankstar/go-skydisc/lib"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestInitDBOrder(t *testing.T) {
	loDB, _ := gorm.Open(sqlite.Open("../"+lib.Server.DBName), lib.Server.DBConfig)
	type args struct {
		iDB *gorm.DB
	}
	lib.Server.Testfile = "../" + lib.Server.Testfile
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Define Test DB 'test.db' in /tmp ",
			args:    args{loDB},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitOrderDB(tt.args.iDB, 1); (err != nil) != tt.wantErr {
				t.Errorf("InitDBOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
