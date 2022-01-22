package importer

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/JoaoVitorStein/folder-manager/pkg/manager"
)

type mockFolderSaver struct {
	count int
}

func (m *mockFolderSaver) Save(ctx context.Context, data manager.Folder) (manager.Folder, error) {
	m.count += 1
	return manager.Folder{}, nil
}
func Test_importer_ImportFile(t *testing.T) {

	tests := []struct {
		name         string
		givenReader  io.Reader
		wantImported int
		wantErr      bool
	}{
		{
			name: "when given a reader should import all lines",
			givenReader: strings.NewReader(`1,0,heading 1,3
			1,0,heading 1,3
			1,0,heading 1,3
			1,0,heading 1,3`),
			wantImported: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockFolderSaver{}
			i := &importer{
				folderSaver: mock,
			}
			if err := i.ImportFile(tt.givenReader); (err != nil) != tt.wantErr {
				t.Errorf("importer.ImportFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantImported != mock.count {
				t.Errorf("importer.ImportFile() imported = %v, want %v", mock.count, tt.wantImported)
			}
		})
	}
}
