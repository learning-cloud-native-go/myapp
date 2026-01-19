package paramsutil

import (
	"net/http"
	"strconv"
)

const (
	defaultPageSize int64 = 10
	maxPageSize     int64 = 100
)

func LimitOffset(r *http.Request) (limit int64, offset int64) {
	q := r.URL.Query()

	l, err := strconv.ParseInt(q.Get("pageSize"), 10, 64)
	if err != nil || l <= 0 {
		limit = defaultPageSize
	} else {
		limit = l
	}
	if limit > maxPageSize {
		limit = maxPageSize
	}

	page, err := strconv.ParseInt(q.Get("page"), 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	offset = (page - 1) * limit

	return limit, offset
}
