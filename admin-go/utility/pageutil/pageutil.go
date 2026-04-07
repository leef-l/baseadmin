package pageutil

const (
	DefaultPageNum  = 1
	DefaultPageSize = 10
	MaxPageSize     = 200
)

func Normalize(pageNum, pageSize int) (int, int) {
	if pageNum <= 0 {
		pageNum = DefaultPageNum
	}
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return pageNum, pageSize
}
