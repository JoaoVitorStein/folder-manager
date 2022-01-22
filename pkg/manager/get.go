package manager

import (
	"context"
	"fmt"
)

type SortInfo struct {
	Field string
	Order string
}

func (m *manager) List(ctx context.Context, nameFilter string, sortInfo SortInfo, page, size int) ([]Folder, error) {
	// For simplicity i'm not validating sort info data
	// this can be a security breach
	query := `SELECT *  FROM folder
		WHERE 1=1 %v 
		ORDER BY %v %v 
		LIMIT %v OFFSET %v`
	params := make([]interface{}, 0)
	filter := ""
	if nameFilter != "" {
		filter += " AND name ILIKE $1"
		params = append(params, nameFilter+"%")
	}

	var folders []Folder
	err := m.db.SelectContext(ctx, &folders, fmt.Sprintf(query, filter, sortInfo.Field, sortInfo.Order, size, size*(page-1)), params...)
	if err != nil {
		return nil, fmt.Errorf("error getting folders from database: %v", err)
	}
	return folders, nil
}

func (m *manager) GetById(ctx context.Context, id int) (Folder, error) {
	query := "SELECT * FROM folder WHERE id = $1"
	var result Folder
	err := m.db.GetContext(ctx, &result, query, id)
	return result, err
}
