package healthcheck

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elliotxx/healthcheck/checks"
	"github.com/gin-gonic/gin"
)

func Example_register() {
	// Example code
	r := gin.New()
	_ = Register(&r.RouterGroup)

	// Simulate request
	dummyRequest(r, "/healthz")

	// Output:
	// 200
	// OK
	//
}

func Example_registerFor() {
	// Example code
	r := gin.New()

	config := NewDefaultConfig()
	config.Verbose = true
	_ = RegisterFor(&r.RouterGroup, config)

	// Simulate request
	dummyRequest(r, "/healthz")

	// Output:
	// 200
	// [+] Ping ok
	// health check passed
	//
}

func Example_customHandler() {
	// Example code
	r := gin.New()

	r.GET("livez", NewHandler(NewDefaultHandlerConfig()))
	readyzChecks := []checks.Check{checks.NewPingCheck(), checks.NewEnvCheck("DB_HOST")}
	r.GET("readyz", NewHandler(NewDefaultHandlerConfigFor(readyzChecks...)))

	// Simulate request
	dummyRequest(r, "/livez")
	dummyRequest(r, "/readyz?verbose")
	dummyRequest(r, "/readyz?verbose&excludes=Env-DB_HOST")

	// Output:
	// 200
	// OK
	//
	// 503
	// [+] Ping ok
	// [-] Env-DB_HOST fail
	// health check failed
	//
	// 200
	// [+] Ping ok
	// health check passed
	//
}

func BenchmarkDefaultHealthCheck(b *testing.B) {
	// Example code
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	r := gin.New()
	_ = Register(&r.RouterGroup)

	// Simulate request for benchmark
	req, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func BenchmarkDefaultHealthCheckWithVerbose(b *testing.B) {
	// Example code
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	r := gin.New()
	_ = Register(&r.RouterGroup)

	// Simulate request for benchmark
	req, _ := http.NewRequest(http.MethodGet, "/healthz?verbose", nil)
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func BenchmarkDefaultHealthCheckWithExcludes(b *testing.B) {
	// Example code
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	r := gin.New()
	_ = Register(&r.RouterGroup)

	// Simulate request for benchmark
	req, _ := http.NewRequest(http.MethodGet, "/healthz?excludes=Ping", nil)
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func BenchmarkDefaultHealthCheckWithVerboseExcludes(b *testing.B) {
	// Example code
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	r := gin.New()
	_ = Register(&r.RouterGroup)

	// Simulate request for benchmark
	req, _ := http.NewRequest(http.MethodGet, "/healthz?verbose&excludes=Ping", nil)
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func dummyRequest(r *gin.Engine, endpoint string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	r.ServeHTTP(w, req)

	fmt.Println(w.Code)
	fmt.Println(w.Body.String())
	fmt.Println()
}
