package shared

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

func LoadTitleMap(ctx context.Context, table string, ids []int64) map[int64]string {
	ids = compactPositiveIDs(ids)
	if table == "" || len(ids) == 0 {
		return nil
	}
	rows, err := g.DB().Ctx(ctx).Model(table).
		Fields("id", "title").
		Where("deleted_at", nil).
		WhereIn("id", ids).
		All()
	if err != nil {
		return nil
	}
	titleMap := make(map[int64]string, len(rows))
	for _, row := range rows {
		titleMap[row["id"].Int64()] = row["title"].String()
	}
	return titleMap
}

func LookupTitle(ctx context.Context, table string, id int64) string {
	if id <= 0 {
		return ""
	}
	return LoadTitleMap(ctx, table, []int64{id})[id]
}

func LoadExistingIDSet(ctx context.Context, table string, ids []int64) map[int64]struct{} {
	ids = compactPositiveIDs(ids)
	if table == "" || len(ids) == 0 {
		return nil
	}
	rows, err := g.DB().Ctx(ctx).Model(table).
		Fields("id").
		Where("deleted_at", nil).
		WhereIn("id", ids).
		All()
	if err != nil {
		return nil
	}
	idSet := make(map[int64]struct{}, len(rows))
	for _, row := range rows {
		idSet[row["id"].Int64()] = struct{}{}
	}
	return idSet
}

func ContainsAllIDs(ctx context.Context, table string, ids []int64) bool {
	ids = compactPositiveIDs(ids)
	if len(ids) == 0 {
		return true
	}
	idSet := LoadExistingIDSet(ctx, table, ids)
	if len(idSet) != len(ids) {
		return false
	}
	for _, id := range ids {
		if _, ok := idSet[id]; !ok {
			return false
		}
	}
	return true
}

func compactPositiveIDs(values []int64) []int64 {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(values))
	normalized := make([]int64, 0, len(values))
	for _, value := range values {
		if value <= 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		normalized = append(normalized, value)
	}
	if len(normalized) == 0 {
		return nil
	}
	return normalized
}
