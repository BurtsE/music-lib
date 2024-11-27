package api

import (
	"bufio"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPagination_ServeHTTP(t *testing.T) {
	tests := []struct {
		name   string
		URL    string
		pageId int
	}{
		{
			name:   "page not provided",
			URL:    "/url/",
			pageId: 0,
		},
		{
			name:   "page zero",
			URL:    "/url/?page=0",
			pageId: 0,
		},
		{
			name:   "page three",
			URL:    "/url/?page=3",
			pageId: 3,
		},
	}
	p := pagination{}
	w := new(httptest.ResponseRecorder)
	ctx := context.Background()
	for _, tt := range tests {
		p.next = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if pageID, ok := req.Context().Value(PageIDKey).(int); !ok || pageID != tt.pageId {
				t.Errorf("paginator did not added valid page id to request. expected: %d, got: %d",
					tt.pageId, pageID)
			}
		})
		r, err := http.NewRequestWithContext(ctx, "GET", tt.URL, &bufio.Reader{})
		if err != nil {
			t.Fatal(err.Error())
		}
		p.ServeHTTP(w, r)
	}
}
