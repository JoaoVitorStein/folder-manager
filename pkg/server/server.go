package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JoaoVitorStein/folder-manager/pkg/manager"
	"github.com/gorilla/mux"
)

type server struct {
	folderLister folderLister
	router       *mux.Router
}

func New(folderLister folderLister) *server {
	s := &server{
		folderLister: folderLister,
	}
	s.setupRouter()
	return s
}

func (s *server) setupRouter() {
	router := mux.NewRouter()
	router.Path("/folders").HandlerFunc(handleListFolders(s.folderLister)).Methods(http.MethodGet)
	s.router = router
}

func (s *server) Start(ctx context.Context) {
	http.ListenAndServe(":8080", s.router)
}

type folderLister interface {
	List(ctx context.Context, nameFilter string, sortInfo manager.SortInfo, page, size int) ([]manager.Folder, error)
}

func handleListFolders(folderLister folderLister) http.HandlerFunc {
	type responseItemAttributes struct {
		ParentID int    `json:"parentId"`
		Name     string `json:"name"`
		Priority int    `json:"priority"`
		PathName string `json:"pathName"`
	}
	type responseItem struct {
		Type       string                 `json:"type"`
		ID         int                    `json:"id"`
		Attributes responseItemAttributes `json:"attributes"`
	}

	type response struct {
		Data []responseItem `json:"data"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		sort := r.FormValue("sort")
		sortInfo := manager.SortInfo{}
		if string(sort[0]) == "-" {
			sortInfo.Order = "DESC"
			sortInfo.Field = sort[1:]
		} else {
			sortInfo.Order = "ASC"
			sortInfo.Field = sort
		}
		nameFilter := r.FormValue("filter")
		page, err := strconv.Atoi(r.FormValue("page"))
		if err != nil {
			writeResponse(rw, http.StatusInternalServerError, nil)
			return
		}
		size, err := strconv.Atoi(r.FormValue("size"))
		if err != nil {
			writeResponse(rw, http.StatusInternalServerError, nil)
			return
		}

		result, err := folderLister.List(r.Context(), nameFilter, sortInfo, page, size)
		if err != nil {
			writeResponse(rw, http.StatusInternalServerError, nil)
			return
		}
		response := response{
			Data: make([]responseItem, 0),
		}
		for _, v := range result {
			response.Data = append(response.Data, responseItem{
				Type: "folder",
				ID:   v.ID,
				Attributes: responseItemAttributes{
					ParentID: int(v.ParentID.Int32),
					Name:     v.Name,
					Priority: v.Priority,
					PathName: v.FullPath,
				},
			})
		}
		writeResponse(rw, http.StatusOK, response)
	}
}

func writeResponse(rw http.ResponseWriter, code int, payload interface{}) {
	res, _ := json.Marshal(payload)
	rw.Header().Set("Content-Type", "application/vnd.api+json")
	rw.WriteHeader(code)
	rw.Write(res)
}
