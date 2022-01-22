package server

import (
	"context"
	"database/sql"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JoaoVitorStein/folder-manager/pkg/manager"
)

type mockFolderLister struct {
	respose []manager.Folder
	err     error
}

func (m *mockFolderLister) List(ctx context.Context, nameFilter string, sortInfo manager.SortInfo, page, size int) ([]manager.Folder, error) {
	return m.respose, m.err
}

func Test_handleListFolders(t *testing.T) {

	tests := []struct {
		name             string
		givenMock        mockFolderLister
		wantStatusCode   int
		wantResponseBody string
	}{
		{
			name: "test",
			givenMock: mockFolderLister{
				respose: []manager.Folder{
					{
						ID:       1,
						ParentID: sql.NullInt32{Int32: int32(2), Valid: true},
						Name:     "teste",
						Priority: 1,
						FullPath: "test",
					},
				},
			},
			wantStatusCode: 200,
			wantResponseBody: `{"data":[{"type":"folder","id":1,"attributes":{"parentId":2,"name":"teste","priority":1,"pathName":"test"}}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &server{
				folderLister: &tt.givenMock,
			}
			server.setupRouter()
			url := "/folders?filter=11&sort=priority&page=1&size=1"
			req, _ := http.NewRequest("GET", url, strings.NewReader(``))
			rr := httptest.NewRecorder()
			server.router.ServeHTTP(rr, req)
			if rr.Code != tt.wantStatusCode {
				t.Errorf("handleGetAccount() status code = %v, want %v", rr.Code, tt.wantStatusCode)
			}
			if got, _ := ioutil.ReadAll(rr.Result().Body); string(got) != tt.wantResponseBody {
				t.Errorf("handleGetAccount() body=%v, want=%v", string(got), tt.wantResponseBody)

			}
		})
	}
}
