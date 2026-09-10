package handler

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/fahrigunadi/ytdl/internal/ytdlp"
)

type InfoHandler struct {
	svc ytdlp.Downloader
}

func NewInfoHandler(svc ytdlp.Downloader) *InfoHandler {
	return &InfoHandler{svc: svc}
}

func (h *InfoHandler) Handle(c *gin.Context) {
	rawURL := strings.TrimSpace(c.PostForm("url"))
	if !IsSupportedURL(rawURL) {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Error": "Invalid URL. Only YouTube and TikTok are supported.",
		})
		return
	}

	info, err := h.svc.GetFormats(c.Request.Context(), rawURL)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Error": "Failed to fetch video info: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "formats.html", gin.H{
		"Info": info,
		"URL":  rawURL,
	})
}

// IsYouTubeURL exported for testing.
func IsYouTubeURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme == "" {
		return false
	}
	host := strings.ToLower(u.Hostname())
	return host == "youtube.com" || host == "www.youtube.com" || host == "youtu.be"
}

// IsTikTokURL exported for testing. Matches tiktok.com and any subdomain
// (www, m, vt, vm, and regional hosts like us.tiktok.com).
func IsTikTokURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme == "" {
		return false
	}
	host := strings.ToLower(u.Hostname())
	return host == "tiktok.com" || strings.HasSuffix(host, ".tiktok.com")
}

// IsSupportedURL reports whether rawURL points at a source we can download from.
func IsSupportedURL(rawURL string) bool {
	return IsYouTubeURL(rawURL) || IsTikTokURL(rawURL)
}
