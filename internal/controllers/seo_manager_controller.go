package controllers

import (
	"fmt"
	"net/http"
	"time"

	"yoyo-server/internal/config"
	"yoyo-server/internal/services"

	"github.com/gin-gonic/gin"
)

type SEOManagerController struct {
	cfg      *config.Config
	services *services.Services
}

func NewSEOManagerController(cfg *config.Config, s *services.Services) *SEOManagerController {
	return &SEOManagerController{cfg: cfg, services: s}
}

func (ctl *SEOManagerController) Sitemap(c *gin.Context) {
	baseUrl := "https://yoyofunnfood.com" // Should come from config in production

	xml := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`

	// 1. Static Pages
	statics := []string{"", "/tickets", "/contact", "/gallery", "/faq"}
	for _, p := range statics {
		xml += fmt.Sprintf(`
	<url>
		<loc>%s%s</loc>
		<lastmod>%s</lastmod>
		<priority>0.8</priority>
	</url>`, baseUrl, p, time.Now().Format("2006-01-02"))
	}

	// 2. Dynamic Content Pages
	contentPages, _ := ctl.services.Content.List(c.Request.Context())
	for _, p := range contentPages {
		if p.IsPublished {
			xml += fmt.Sprintf(`
	<url>
		<loc>%s/%s</loc>
		<lastmod>%s</lastmod>
		<priority>0.6</priority>
	</url>`, baseUrl, p.Slug, p.UpdatedAt.Format("2006-01-02"))
		}
	}

	// 3. Suites
	suites, _ := ctl.services.Suites.ListPublic(c.Request.Context())
	for _, s := range suites {
		xml += fmt.Sprintf(`
	<url>
		<loc>%s/suites/%s</loc>
		<lastmod>%s</lastmod>
		<priority>0.7</priority>
	</url>`, baseUrl, s.Slug, s.UpdatedAt.Format("2006-01-02"))
	}

	xml += "\n</urlset>"

	c.Header("Content-Type", "application/xml")
	c.String(http.StatusOK, xml)
}

func (ctl *SEOManagerController) Robots(c *gin.Context) {
	baseUrl := "https://yoyofunnfood.com"
	robots := fmt.Sprintf(`User-agent: *
Allow: /
Disallow: /admin/
Disallow: /api/

Sitemap: %s/sitemap.xml`, baseUrl)

	c.Header("Content-Type", "text/plain")
	c.String(http.StatusOK, robots)
}
