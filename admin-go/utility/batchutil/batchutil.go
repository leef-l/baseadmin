package batchutil

import "gbaseadmin/utility/snowflake"

type TreeRow struct {
	ID       int64
	ParentID int64
}

func CompactIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
	if len(ids) == 0 {
		return nil
	}
	seen := make(map[snowflake.JsonInt64]struct{}, len(ids))
	normalized := make([]snowflake.JsonInt64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		normalized = append(normalized, id)
	}
	if len(normalized) == 0 {
		return nil
	}
	return normalized
}

func ToInt64s(ids []snowflake.JsonInt64) []int64 {
	if len(ids) == 0 {
		return nil
	}
	values := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		values = append(values, int64(id))
	}
	if len(values) == 0 {
		return nil
	}
	return values
}

func IDSet(ids []snowflake.JsonInt64) map[int64]struct{} {
	if len(ids) == 0 {
		return nil
	}
	values := make(map[int64]struct{}, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		values[int64(id)] = struct{}{}
	}
	if len(values) == 0 {
		return nil
	}
	return values
}

func ExpandTreeDeleteOrder(selected []snowflake.JsonInt64, rows []TreeRow) []snowflake.JsonInt64 {
	selected = CompactIDs(selected)
	if len(selected) == 0 || len(rows) == 0 {
		return nil
	}
	rowMap := make(map[int64]struct{}, len(rows))
	childrenMap := make(map[int64][]int64, len(rows))
	for _, row := range rows {
		if row.ID <= 0 {
			continue
		}
		rowMap[row.ID] = struct{}{}
		childrenMap[row.ParentID] = append(childrenMap[row.ParentID], row.ID)
	}
	visited := make(map[int64]struct{}, len(rows))
	order := make([]snowflake.JsonInt64, 0, len(rows))
	var walk func(id int64)
	walk = func(id int64) {
		if id <= 0 {
			return
		}
		if _, ok := rowMap[id]; !ok {
			return
		}
		if _, ok := visited[id]; ok {
			return
		}
		visited[id] = struct{}{}
		for _, childID := range childrenMap[id] {
			walk(childID)
		}
		order = append(order, snowflake.JsonInt64(id))
	}
	for _, id := range selected {
		walk(int64(id))
	}
	if len(order) == 0 {
		return nil
	}
	return order
}
