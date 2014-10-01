package mock

import (
	"net/http"
	"net/http/httptest"
)

type Route struct {
	Pattern string
	Handler func(w http.ResponseWriter, r *http.Request)
}

func RunHydraServerMock(routes []Route) *httptest.Server {
	hydraMux := http.NewServeMux()
	for _, route := range routes {
		hydraMux.Handle(route.Pattern, http.HandlerFunc(route.Handler))
	}
	return httptest.NewServer(hydraMux)
}
