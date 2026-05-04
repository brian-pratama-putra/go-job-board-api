package main

import (
	"net/http"
	"os"

	"go-job-board-api/config"
	controllers "go-job-board-api/controllers/api"
	"go-job-board-api/middleware"
	"go-job-board-api/others"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config.GetPgPool()
	config.GetRedisClient()

	if os.Getenv("STATUS_APP") != "DEV" {
		gin.SetMode(gin.ReleaseMode)
	}

	v_router := gin.New()
	v_router.Use(gin.Recovery())
	v_router.Use(middleware.SecurityMiddleware())
	v_router.Use(middleware.SecurityHeaders())

	v_router.NoRoute(func(p_c *gin.Context) {
		p_c.JSON(http.StatusOK, gin.H{
			"err_code": 404,
			"err_msg":  "Route not found.",
			"datetime": others.NowJakartaStr(),
		})
	})

	v_router.GET("/", func(p_c *gin.Context) {
		p_c.JSON(http.StatusOK, gin.H{"status": "running"})
	})

	v_router.GET("/health", func(p_c *gin.Context) {
		p_c.JSON(http.StatusOK, gin.H{"status": "running"})
	})

	v_inquiry := v_router.Group("/")
	v_inquiry.Use(middleware.ParseBody())
	v_inquiry.Use(middleware.RateLimiter())
	{
		v_inquiry.POST("/inquiry", InquiryHandler)
	}

	v_router.Run(":" + config.GetPort())
}

func InquiryHandler(p_c *gin.Context) {
	v_body, _ := p_c.Get("parsed_body")
	v_data     := v_body.(map[string]interface{})

	v_method    := ""
	v_datetime  := ""
	if v_m, v_ok := v_data["method"].(string); v_ok {
		v_method = v_m
	}
	if v_dt, v_ok := v_data["datetime"].(string); v_ok {
		v_datetime = v_dt
	}

	if v_datetime != "" {
		v_is_valid, v_msg := others.ValidateRequestDatetime(v_datetime)
		if !v_is_valid {
			others.ResponseService(p_c, v_method, 405, v_msg, nil)
			return
		}
	}

	switch v_method {
	case "post_job":
		controllers.PostJob(p_c)
	case "get_job_list":
		controllers.GetJobList(p_c)
	case "get_job_detail":
		controllers.GetJobDetail(p_c)
	case "update_job":
		controllers.UpdateJob(p_c)
	case "close_job":
		controllers.CloseJob(p_c)
	case "apply_job":
		controllers.ApplyJob(p_c)
	case "update_application_status":
		controllers.UpdateApplicationStatus(p_c)
	case "get_my_application":
		controllers.GetMyApplication(p_c)
	case "get_job_applicants":
		controllers.GetJobApplicants(p_c)
	case "get_matched_jobs":
		controllers.GetMatchedJobs(p_c)
	case "get_job_summary":
		controllers.GetJobSummary(p_c)
	default:
		others.ResponseService(p_c, v_method, 405, "Invalid Method", nil)
	}
}
