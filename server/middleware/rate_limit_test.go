package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func limiterRouter(t *testing.T, limiter *ipLimiter) *gin.Engine {
	t.Helper()
	previous := global.I18
	global.I18 = &global.I18n{Bundle: i18n.NewBundle(language.English)}
	t.Cleanup(func() { global.I18 = previous })
	router := gin.New()
	router.Use(limiter.Handle)
	router.GET("/", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	return router
}

func limitedRequest(router http.Handler, ip string) int {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.RemoteAddr = ip + ":12345"
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response.Code
}

func TestIPLimiterContinuousTrafficReplenishesBudget(t *testing.T) {
	router := limiterRouter(t, newIPLimiter(200*time.Millisecond, 3))
	for i := 0; i < 4; i++ {
		if i > 0 {
			time.Sleep(80 * time.Millisecond)
		}
		if got := limitedRequest(router, "192.0.2.1"); got != http.StatusNoContent {
			t.Fatalf("request %d blocked despite crossing into the next window: %d", i+1, got)
		}
	}
}

func TestIPLimiterRejectsBurstAndPreservesNeighbor(t *testing.T) {
	router := limiterRouter(t, newIPLimiter(time.Hour, 3))
	for i := 0; i < 3; i++ {
		if got := limitedRequest(router, "192.0.2.1"); got != http.StatusNoContent {
			t.Fatalf("request inside budget rejected: %d", got)
		}
	}
	if got := limitedRequest(router, "192.0.2.1"); got != http.StatusTooManyRequests {
		t.Fatalf("burst outside budget accepted: %d", got)
	}
	if got := limitedRequest(router, "192.0.2.2"); got != http.StatusNoContent {
		t.Fatalf("neighbor IP blocked: %d", got)
	}
}
