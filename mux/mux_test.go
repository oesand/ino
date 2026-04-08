package mux_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/oesand/ino/mux"
)

func TestNew(t *testing.T) {
	m := mux.New(
		mux.Get("/hello", simpleHandler("hello")),
		mux.Get("/world", simpleHandler("world")),
	)
	if m == nil {
		t.Fatal("New() with routes returned nil")
	}

	routeCount := 0
	for range m.Routes() {
		routeCount++
	}
	if routeCount != 2 {
		t.Errorf("expected 2 routes, got %d", routeCount)
	}
}

func TestInclude(t *testing.T) {
	m := mux.New()
	m.Include(
		mux.Get("/test", simpleHandler("test")),
		mux.Post("/submit", simpleHandler("submit")),
	)

	routeCount := 0
	for range m.Routes() {
		routeCount++
	}
	if routeCount != 2 {
		t.Errorf("expected 2 routes, got %d", routeCount)
	}
}

func TestRouteCreation(t *testing.T) {
	tests := []struct {
		name   string
		route  mux.Route
		method string
	}{
		{"Get", mux.Get("/test", simpleHandler("test")), http.MethodGet},
		{"Post", mux.Post("/test", simpleHandler("test")), http.MethodPost},
		{"Put", mux.Put("/test", simpleHandler("test")), http.MethodPut},
		{"Delete", mux.Delete("/test", simpleHandler("test")), http.MethodDelete},
		{"Options", mux.Options("/test", simpleHandler("test")), http.MethodOptions},
		{"Head", mux.Head("/test", simpleHandler("test")), http.MethodHead},
		{"Connect", mux.Connect("/test", simpleHandler("test")), http.MethodConnect},
		{"Patch", mux.Patch("/test", simpleHandler("test")), http.MethodPatch},
		{"Trace", mux.Trace("/test", simpleHandler("test")), http.MethodTrace},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.route.Method() != tt.method {
				t.Errorf("expected method %s, got %s", tt.method, tt.route.Method())
			}
			if tt.route.Pattern() != "/test" {
				t.Errorf("expected pattern /test, got %s", tt.route.Pattern())
			}
		})
	}
}

func TestRouteAttributes(t *testing.T) {
	attr1 := "auth_required"
	attr2 := 42
	route := mux.Get("/protected", simpleHandler("test"), attr1, attr2)
	attrs := route.Attrs()
	if len(attrs) != 2 {
		t.Errorf("expected 2 attributes, got %d", len(attrs))
	}
	if attrs[0] != attr1 {
		t.Errorf("expected %v, got %v", attr1, attrs[0])
	}
	if attrs[1] != attr2 {
		t.Errorf("expected %v, got %v", attr2, attrs[1])
	}
}

func TestHandleInvalidPattern(t *testing.T) {
	tests := []struct {
		name        string
		pattern     string
		shouldPanic bool
	}{
		{"empty pattern", "", true},
		{"no leading slash", "hello", true},
		{"valid pattern", "/hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.shouldPanic {
						t.Errorf("unexpected panic: %v", r)
					}
				} else if tt.shouldPanic {
					t.Error("expected panic but got none")
				}
			}()
			mux.Get(tt.pattern, simpleHandler("test"))
		})
	}
}

func TestHandleNilHandler(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil handler")
		}
	}()
	mux.Handle(http.MethodGet, "/test", nil)
}

func TestServeHTTP(t *testing.T) {
	m := mux.New(
		mux.Get("/", simpleHandler("root")),
		mux.Get("/hello", simpleHandler("hello")),
		mux.Post("/submit", simpleHandler("submit")),
	)

	tests := []struct {
		method   string
		path     string
		expected string
	}{
		{http.MethodGet, "/", "root"},
		{http.MethodGet, "/hello", "hello"},
		{http.MethodPost, "/submit", "submit"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s %s", tt.method, tt.path), func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			m.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
			}
			body, _ := io.ReadAll(w.Body)
			if string(body) != tt.expected {
				t.Errorf("expected body %s, got %s", tt.expected, string(body))
			}
		})
	}
}

func TestNotFound(t *testing.T) {
	m := mux.New(
		mux.Get("/hello", simpleHandler("hello")),
	)

	m.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "custom not found")
	}))

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
	body, _ := io.ReadAll(w.Body)
	if string(body) != "custom not found" {
		t.Errorf("expected body 'custom not found', got %s", string(body))
	}
}

func TestDefaultNotFound(t *testing.T) {
	m := mux.New(
		mux.Get("/hello", simpleHandler("hello")),
	)

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestURLParameters(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		path     string
		expected map[string]string
	}{
		{
			"simple param",
			"/users/{id}",
			"/users/123",
			map[string]string{"id": "123"},
		},
		{
			"multiple params",
			"/posts/{year}/{slug}",
			"/posts/2024/hello-world",
			map[string]string{"year": "2024", "slug": "hello-world"},
		},
		{
			"param with regex",
			"/users/{id:\\d+}",
			"/users/456",
			map[string]string{"id": "456"},
		},
		{
			"multiple params with regex",
			"/posts/{year:\\d{4}}/{slug:[^/]+}",
			"/posts/2024/my-slug",
			map[string]string{"year": "2024", "slug": "my-slug"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mux.New(
				mux.Get(tt.pattern, func(w http.ResponseWriter, r *http.Request) {
					params := mux.UrlParams(r.Context())
					for k, v := range tt.expected {
						if val, ok := params[k]; !ok || val != v {
							w.WriteHeader(http.StatusBadRequest)
							return
						}
					}
					w.WriteHeader(http.StatusOK)
				}),
			)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			m.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
			}
		})
	}
}

func TestMatchedRoute(t *testing.T) {
	m := mux.New(
		mux.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			route := mux.MatchedRoute(r.Context())
			if route == nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if route.Method() != http.MethodGet {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if route.Pattern() != "/test" {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestPrefixRoutes(t *testing.T) {
	routes := mux.PrefixRoutes("/api",
		mux.Get("/users", simpleHandler("users")),
		mux.Get("/posts", simpleHandler("posts")),
	)
	m := mux.New(routes...)

	tests := []struct {
		path     string
		expected string
	}{
		{"/api/users", "users"},
		{"/api/posts", "posts"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			m.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
			}
			body, _ := io.ReadAll(w.Body)
			if string(body) != tt.expected {
				t.Errorf("expected body %s, got %s", tt.expected, string(body))
			}
		})
	}
}

func TestNestedPrefixRoutes(t *testing.T) {
	routes := mux.PrefixRoutes("/api",
		mux.PrefixRoutes("/v1",
			mux.Get("/users", simpleHandler("v1_users")),
		),
		mux.PrefixRoutes("/v2",
			mux.Get("/users", simpleHandler("v2_users")),
		),
	)
	m := mux.New(routes...)

	tests := []struct {
		path     string
		expected string
	}{
		{"/api/v1/users", "v1_users"},
		{"/api/v2/users", "v2_users"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			m.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
			}
			body, _ := io.ReadAll(w.Body)
			if string(body) != tt.expected {
				t.Errorf("expected body %s, got %s", tt.expected, string(body))
			}
		})
	}
}

func TestMiddlewareOrder(t *testing.T) {
	callOrder := []string{}
	m := mux.New(
		mux.Get("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "handler")
			w.WriteHeader(http.StatusOK)
		})),
	)

	mw1 := mux.Middleware(func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		callOrder = append(callOrder, "middleware1")
		next.ServeHTTP(w, r)
	})
	m.Middleware(mw1)
	mw2 := mux.Middleware(func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		callOrder = append(callOrder, "middleware2")
		next.ServeHTTP(w, r)
	})
	m.Middleware(mw2)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	expected := "middleware1,middleware2,handler"
	if got := strings.Join(callOrder, ","); got != expected {
		t.Fatalf("expected call order %s, got %s", expected, got)
	}
}

func TestMiddlewareShortCircuit(t *testing.T) {
	called := false
	m := mux.New(
		mux.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		}),
	)

	blocker := mux.Middleware(func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "blocked")
	})
	m.Middleware(blocker)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, req)

	if called {
		t.Fatal("expected route handler not to be called when middleware short-circuits")
	}
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, w.Code)
	}
	if body, _ := io.ReadAll(w.Body); string(body) != "blocked" {
		t.Fatalf("expected blocked response body, got %s", string(body))
	}

}

func TestRoutesIterator(t *testing.T) {
	m := mux.New(
		mux.Get("/a", simpleHandler("a")),
		mux.Post("/b", simpleHandler("b")),
		mux.Put("/c", simpleHandler("c")),
	)

	routeCount := 0
	methods := make(map[string]bool)
	for route := range m.Routes() {
		routeCount++
		methods[route.Method()] = true
	}

	if routeCount != 3 {
		t.Errorf("expected 3 routes, got %d", routeCount)
	}
	if !methods[http.MethodGet] || !methods[http.MethodPost] || !methods[http.MethodPut] {
		t.Errorf("expected all methods in routes, got %v", methods)
	}
}

func TestTrailingSlash(t *testing.T) {
	m := mux.New(
		mux.Get("/hello", simpleHandler("hello")),
	)

	tests := []string{"/hello", "/hello/"}

	for _, path := range tests {
		t.Run(fmt.Sprintf("path %s", path), func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			w := httptest.NewRecorder()
			m.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
			}
			body, _ := io.ReadAll(w.Body)
			if string(body) != "hello" {
				t.Errorf("expected body 'hello', got %s", string(body))
			}
		})
	}
}

func TestRouteSorting(t *testing.T) {
	// Routes are sorted by depth and number of parameters
	// This test ensures more specific routes (deeper or with fewer params) take precedence
	m := mux.New(
		mux.Get("/posts", simpleHandler("list")),
		mux.Get("/posts/{id}", simpleHandler("detail")),
		mux.Get("/posts/{id}/comments", simpleHandler("comments")),
	)

	tests := []struct {
		path     string
		expected string
	}{
		{"/posts", "list"},
		{"/posts/123", "detail"},
		{"/posts/123/comments", "comments"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			m.ServeHTTP(w, req)

			body, _ := io.ReadAll(w.Body)
			if string(body) != tt.expected {
				t.Errorf("path %s: expected %s, got %s", tt.path, tt.expected, string(body))
			}
		})
	}
}

func TestMultipleRouteRegistration(t *testing.T) {
	m := mux.New(
		mux.Get("/hello", simpleHandler("hello")),
	)

	// Register more routes
	m.Include(
		mux.Post("/hello", simpleHandler("posted")),
		mux.Put("/hello", simpleHandler("updated")),
	)

	tests := []struct {
		method   string
		expected string
	}{
		{http.MethodGet, "hello"},
		{http.MethodPost, "posted"},
		{http.MethodPut, "updated"},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/hello", nil)
			w := httptest.NewRecorder()
			m.ServeHTTP(w, req)

			body, _ := io.ReadAll(w.Body)
			if string(body) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(body))
			}
		})
	}
}

// simpleHandler returns a handler that writes the given text to the response
func simpleHandler(text string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, text)
	}
}
