package importer

import (
	"context"
	"database/sql"
	"encoding/csv"
	"io"
	"strconv"

	"github.com/JoaoVitorStein/folder-manager/pkg/manager"
)

type FolderSaver interface {
	Save(ctx context.Context, data manager.Folder) (manager.Folder, error)
}

type importer struct {
	folderSaver FolderSaver
}

func New(folderSaver FolderSaver) *importer {
	return &importer{
		folderSaver: folderSaver,
	}
}

func (i *importer) ImportFile(fileReader io.Reader) error {
	reader := csv.NewReader(fileReader)
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		id, _ := strconv.Atoi(row[0])
		parentID, _ := strconv.Atoi(row[1])
		priority, _ := strconv.Atoi(row[3])
		folder := manager.Folder{
			ID:       id,
			ParentID: sql.NullInt32{Int32: int32(parentID), Valid: parentID != 0},
			Name:     row[2],
			Priority: priority,
		}
		i.folderSaver.Save(context.Background(), folder)
	}
	return nil
}
