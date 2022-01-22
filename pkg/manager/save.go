package manager

import (
	"context"
	"fmt"
)

func (m *manager) Save(ctx context.Context, data Folder) (Folder, error) {
	if data.ParentID.Int32 != 0 {
		parent, err := m.GetById(ctx, int(data.ParentID.Int32))
		if err != nil {
			return Folder{}, fmt.Errorf("error getting folder parent %v", err)
		}
		data.FullPath += parent.FullPath + "/" + data.Name
	} else {
		data.FullPath = data.Name
	}
	query := `INSERT INTO folder(id, parent, name, priority, full_path) 
			VALUES ($1, $2, $3, $4, $5) 
			RETURNING id, parent, name, priority, full_path`

	var folder Folder
	if err := m.db.GetContext(ctx, &folder, query, data.ID, data.ParentID, data.Name, data.Priority, data.FullPath); err != nil {
		return Folder{}, fmt.Errorf("error saving folder: %v", err)
	}
	return folder, nil
}
