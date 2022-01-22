package manager

import (
	"context"
	"reflect"
	"testing"
)

func Test_manager_Save(t *testing.T) {
	pg := dockerPg{}
	db := pg.Start(t)
	defer pg.Stop(t)
	tests := []struct {
		name        string
		givenFolder Folder
		wantFolder  Folder
		wantErr     bool
	}{
		{
			name: "when given a valid folder should save in database",
			givenFolder: Folder{
				ID:       1,
				Name:     "Test",
				Priority: 1,
			},
			wantFolder: Folder{
				ID:       1,
				Name:     "Test",
				Priority: 1,
				FullPath: "Test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := manager{
				db: db,
			}
			got, err := manager.Save(context.Background(), tt.givenFolder)
			if (err != nil) != tt.wantErr {
				t.Fatalf("manager.Save() err = %v, wantErr = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.wantFolder) {
				t.Fatalf("manager.Save() = %v, want = %v", got, tt.wantFolder)
			}
		})
	}
}
