package shared

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

func LoadDirNameMap(ctx context.Context, dirIDs []int64) map[int64]string {
	dirIDs = compactDirIDs(dirIDs)
	if len(dirIDs) == 0 {
		return nil
	}
	rows, err := g.DB().Ctx(ctx).Model("upload_dir").
		Fields("id", "name").
		Where("deleted_at", nil).
		WhereIn("id", dirIDs).
		All()
	if err != nil {
		return nil
	}
	dirMap := make(map[int64]string, len(rows))
	for _, row := range rows {
		dirMap[row["id"].Int64()] = row["name"].String()
	}
	return dirMap
}

func LookupDirName(ctx context.Context, dirID int64) string {
	if dirID <= 0 {
		return ""
	}
	return LoadDirNameMap(ctx, []int64{dirID})[dirID]
}

func compactDirIDs(values []int64) []int64 {
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
