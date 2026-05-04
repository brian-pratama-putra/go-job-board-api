package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var v_scanner_patterns = []string{"censys", "zgrab", "modtscanner", "nmap", "masscan"}
var v_invalid_paths    = []*regexp.Regexp{
	regexp.MustCompile(`(?i)/\.env`),
	regexp.MustCompile(`(?i)/security\.txt`),
	regexp.MustCompile(`(?i)/\.git`),
	regexp.MustCompile(`(?i)/phpmyadmin`),
	regexp.MustCompile(`(?i)/wp-.*`),
}

func SecurityMiddleware() gin.HandlerFunc {
	return func(p_c *gin.Context) {
		v_user_agent := strings.ToLower(p_c.GetHeader("User-Agent"))
		for _, v_pattern := range v_scanner_patterns {
			if strings.Contains(v_user_agent, v_pattern) {
				p_c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"detail": "Blocked"})
				return
			}
		}
		v_path := p_c.Request.URL.Path
		for _, v_re := range v_invalid_paths {
			if v_re.MatchString(v_path) {
				p_c.AbortWithStatus(http.StatusNotFound)
				return
			}
		}
		p_c.Next()
	}
}

func SecurityHeaders() gin.HandlerFunc {
	return func(p_c *gin.Context) {
		p_c.Next()
		p_c.Header("X-XSS-Protection", "1; mode=block")
		p_c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		p_c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		p_c.Header("X-Frame-Options", "DENY")
		p_c.Header("X-Content-Type-Options", "nosniff")
		p_c.Header("Permissions-Policy", "accelerometer=(), autoplay=(), camera=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), payment=(), usb=(), browsing-topics=()")
		p_c.Header("Cross-Origin-Embedder-Policy", "require-corp")
		p_c.Header("Cross-Origin-Opener-Policy", "same-origin")
		p_c.Header("Cross-Origin-Resource-Policy", "same-origin")
	}
}

func RateLimiter() gin.HandlerFunc {
	return func(p_c *gin.Context) {
		p_c.Next()
	}
}

func ParseBody() gin.HandlerFunc {
	return func(p_c *gin.Context) {
		var v_body map[string]interface{}
		if v_err := p_c.ShouldBindJSON(&v_body); v_err == nil {
			p_c.Set("parsed_body", v_body)
		} else {
			p_c.Set("parsed_body", map[string]interface{}{})
		}
		p_c.Next()
	}
}
