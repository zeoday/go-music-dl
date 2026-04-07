package web

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAppJSMediaSessionArtworkUsesCoverProxy(t *testing.T) {
	content, err := templateFS.ReadFile("templates/static/js/app.js")
	if err != nil {
		t.Fatalf("ReadFile(app.js): %v", err)
	}

	js := string(content)
	if !strings.Contains(js, "function buildMediaSessionCoverURL(audio = getCurrentAPlayerAudio())") {
		t.Fatal("app.js missing buildMediaSessionCoverURL helper")
	}
	if !strings.Contains(js, "cover_proxy") {
		t.Fatal("app.js missing cover_proxy media session artwork path")
	}
	if !strings.Contains(js, "function scheduleMediaSessionSync(audio = getCurrentAPlayerAudio(), delayMs = 160)") {
		t.Fatal("app.js missing delayed media session resync helper")
	}
}

func TestCoverProxyReturnsInlineImage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	imageBytes := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x01, 0x02, 0x03, 0x04}
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(imageBytes)
	}))
	defer upstream.Close()

	router := gin.New()
	RegisterMusicRoutes(router.Group(RoutePrefix))

	req := httptest.NewRequest(http.MethodGet, RoutePrefix+"/cover_proxy?url="+url.QueryEscape(upstream.URL), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("Content-Type"); !strings.Contains(got, "image/png") {
		t.Fatalf("Content-Type = %q, want image/png", got)
	}
	if got := rec.Header().Get("Cache-Control"); got != "public, max-age=21600" {
		t.Fatalf("Cache-Control = %q, want public, max-age=21600", got)
	}
	if !bytes.Equal(rec.Body.Bytes(), imageBytes) {
		t.Fatal("cover_proxy returned unexpected body")
	}
}
