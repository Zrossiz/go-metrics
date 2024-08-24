package get

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Zrossiz/go-metrics/internal/server/storage"
	memstorage "github.com/Zrossiz/go-metrics/internal/server/storage/memStorage"
	"github.com/go-chi/chi/v5"
)

func TestMetric(t *testing.T) {

	store := memstorage.NewMemStorage()

	store.Metrics = []storage.Metric{
		{Name: "testMetric", Type: storage.GaugeType, Value: 42},
	}

	req := httptest.NewRequest(http.MethodGet, "/value/testMetric", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "testMetric")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	Metric(rr, req, *store)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	expected := "42"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("expected body %s, got %s", expected, rr.Body.String())
	}
}
