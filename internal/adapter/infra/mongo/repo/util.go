package repo

func calculateSkip(page, size int64) int64 {
	if page < 0 {
		page = 0
	}
	if size <= 0 {
		size = 10
	}
	return page * size
}
