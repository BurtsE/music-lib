package api

import (
	"context"
	"net/http"
	"strconv"
)

const PageIDKey = "page"

// Хэндлер для обеспечения пагинации
type pagination struct {
	next http.Handler
}

func paginate(handler http.Handler) http.Handler {
	return &pagination{next: handler}
}

func (p *pagination) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pageID := req.URL.Query().Get(PageIDKey)
	intPageID := 0
	var err error
	if pageID != "" {
		intPageID, err = strconv.Atoi(pageID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	ctx := context.WithValue(req.Context(), PageIDKey, intPageID)
	p.next.ServeHTTP(w, req.WithContext(ctx))
}
