package core

import "testing"

func TestInitCatalogDB(t *testing.T) {
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
			name: "Test InitCatalogDB()",
			args: args{
				iMode: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitCatalogDB(tt.args.iMode); (err != nil) != tt.wantErr {
				t.Errorf("InitCatalogDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
