package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/yourusername/vigilant-web/internal/alerts"
	"github.com/yourusername/vigilant-web/internal/crawler"
	"github.com/yourusername/vigilant-web/internal/db"
)

var logger = logrus.New()

type Server struct {
	DB     *sql.DB
	Config struct {
		HTTPPort  int
		TorProxy  string
		JWTSecret string
	}
	AlertEngine *alerts.AlertEngine
}

func NewServer(database *sql.DB, httpPort int, torProxy, jwtSecret string) *Server {
	InitAuth(jwtSecret) // from auth_middleware
	return &Server{
		DB: database,
		Config: struct {
			HTTPPort  int
			TorProxy  string
			JWTSecret string
		}{
			HTTPPort:  httpPort,
			TorProxy:  torProxy,
			JWTSecret: jwtSecret,
		},
		AlertEngine: &alerts.AlertEngine{DB: database},
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Production Hardening: add security headers
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
		c.Next()
	})

	// Example of rate limiting? There's a third-party library e.g. "github.com/aviddiviner/gin-limit"
	// or you can roll your own. We'll skip that for brevity.

	// Serve static index at root
	r.StaticFile("/", "./web/index.html")

	// Public route for user login or signup
	r.POST("/login", s.handleLogin)

	// Auth group
	authGroup := r.Group("/api")
	authGroup.Use(AuthRequired()) // all routes in here require JWT
	{
		authGroup.GET("/pages", s.handleListPages)
		authGroup.GET("/alerts", s.handleListAlerts)

		authGroup.POST("/watchlist", RequireRole("admin"), s.handleAddWatchlist) // only admin can add patterns

		authGroup.POST("/crawl/onion", s.handleCrawlOnion)
		authGroup.POST("/crawl/rss", s.handleCrawlRSS)
	}

	return r
}

func (s *Server) Run() error {
	r := s.SetupRoutes()
	addr := fmt.Sprintf(":%d", s.Config.HTTPPort)
	logger.Infof("Starting server on %s", addr)
	return r.Run(addr)
}

// Example: handleLogin - we skip real DB user check for brevity
func (s *Server) handleLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}
	// TODO: check user in DB, compare hashed password, etc.
	if req.Username == "admin" && req.Password == "admin123" {
		token, err := GenerateToken("admin", "admin")
		if err != nil {
			logger.Error("Failed to generate token: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
		return
	} else if req.Username == "user" && req.Password == "user123" {
		token, err := GenerateToken("user", "user")
		if err != nil {
			logger.Error("Failed to generate token: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
}

func (s *Server) handleListPages(c *gin.Context) {
	pages, err := db.ListPages(s.DB, 50)
	if err != nil {
		logger.Error("ListPages error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pages)
}

func (s *Server) handleListAlerts(c *gin.Context) {
	alertsList, err := db.ListAlerts(s.DB, 50)
	if err != nil {
		logger.Error("ListAlerts error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, alertsList)
}

func (s *Server) handleAddWatchlist(c *gin.Context) {
	var req struct {
		Pattern     string `json:"pattern"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	if err := db.AddWatchlistPattern(s.DB, req.Pattern, req.Description); err != nil {
		logger.Error("AddWatchlistPattern error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "pattern added"})
}

func (s *Server) handleCrawlOnion(c *gin.Context) {
	var req struct {
		URL string `json:"url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	if req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing onion url"})
		return
	}
	results, err := crawler.CrawlOnionSite(req.URL, s.Config.TorProxy)
	if err != nil {
		logger.Error("CrawlOnionSite error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "crawl error"})
		return
	}

	var pagesCount, alertsCount int
	for url, html := range results {
		if err := db.InsertPage(s.DB, url, "", html); err != nil {
			logger.Error("InsertPage error: ", err)
			continue
		}
		pagesCount++
		alertFired := s.AlertEngine.EvaluateContent(url, html)
		alertsCount += alertFired
	}
	c.JSON(http.StatusOK, gin.H{
		"status":         "onion crawl complete",
		"pages_inserted": pagesCount,
		"alerts_fired":   alertsCount,
	})
}

func (s *Server) handleCrawlRSS(c *gin.Context) {
	var req struct {
		FeedURL string `json:"feed_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	items, err := crawler.CrawlRSSFeed(req.FeedURL)
	if err != nil {
		logger.Error("CrawlRSSFeed error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "crawl error"})
		return
	}

	var pagesCount, alertsCount int
	for _, item := range items {
		link := item.Link
		title := item.Title
		content := item.Content
		if content == "" {
			content = item.Description
		}

		if err := db.InsertPage(s.DB, link, title, content); err != nil {
			logger.Error("InsertPage error: ", err)
			continue
		}
		pagesCount++
		alertFired := s.AlertEngine.EvaluateContent(link, title+" "+content)
		alertsCount += alertFired
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         "rss crawl complete",
		"pages_inserted": pagesCount,
		"alerts_fired":   alertsCount,
	})
}
