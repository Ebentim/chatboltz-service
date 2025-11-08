package handler

import (
	"net/http"

	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/alpinesboltltd/boltz-ai/internal/scraper"
	"github.com/gin-gonic/gin"
)

type ScraperHandler struct {
	svc *scraper.Service
}

func NewScraperHandler(svc *scraper.Service) *ScraperHandler {
	return &ScraperHandler{svc: svc}
}

// POST /api/v1/scrape
func (h *ScraperHandler) Scrape(c *gin.Context) {
	var req struct {
		URL      string   `json:"url" binding:"required"`
		Trace    bool     `json:"trace"`
		Exclude  []string `json:"exclude"`
		MaxPages int      `json:"max_pages"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "Scrape")
		return
	}

	opts := scraper.ScrapeOptions{Trace: req.Trace, Exclude: req.Exclude, MaxPages: req.MaxPages}
	res, err := h.svc.Scrape(c.Request.Context(), req.URL, opts)
	if err != nil {
		appErrors.HandleError(c, appErrors.NewInternalError("Scrape failed", err.Error()), "Scrape")
		return
	}

	c.JSON(http.StatusOK, res)
}
