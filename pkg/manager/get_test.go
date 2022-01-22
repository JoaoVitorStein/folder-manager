package manager

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
)

func Test_manager_List(t *testing.T) {
	pg := dockerPg{}
	db := pg.Start(t)
	defer pg.Stop(t)
	manager := manager{
		db: db,
	}
	folderOne, _ := manager.Save(context.Background(), Folder{
		ID:       1,
		Name:     "Test",
		Priority: 1,
	})
	folderTwo, _ := manager.Save(context.Background(), Folder{
		ID:       2,
		Name:     "Another folder",
		Priority: 1,
		ParentID: sql.NullInt32{Int32: int32(folderOne.ID), Valid: true},
		FullPath: "Test/Another folder",
	})
	tests := []struct {
		name            string
		givenNameFilter string
		givenSortInfo   SortInfo
		givenPage       int
		givenSize       int
		want            []Folder
		wantErr         bool
	}{
		{
			name: "when given an invalid sort order should return an error",
			givenSortInfo: SortInfo{
				Field: "aa",
				Order: "asv",
			},
			wantErr: true,
		},
		{
			name: "when given a name that doesn't exist should return an empty array",
			givenSortInfo: SortInfo{
				Field: "priority",
				Order: "ASC",
			},
			givenPage:       1,
			givenSize:       10,
			givenNameFilter: "not test",
			wantErr:         false,
		},
		{
			name: "when no filter should return all data with pagination",
			givenSortInfo: SortInfo{
				Field: "priority",
				Order: "ASC",
			},
			givenPage: 1,
			givenSize: 10,
			wantErr:   false,
			want:      []Folder{folderOne, folderTwo},
		},
		{
			name: "when size 1 pages should have one item",
			givenSortInfo: SortInfo{
				Field: "priority",
				Order: "ASC",
			},
			givenPage: 2,
			givenSize: 1,
			wantErr:   false,
			want:      []Folder{folderTwo},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := manager.List(context.Background(), tt.givenNameFilter, tt.givenSortInfo, tt.givenPage, tt.givenSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("manager.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("manager.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
