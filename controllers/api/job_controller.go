package controllers

import (
	"os"
	"strconv"

	dao "go-job-board-api/dao/api"
	"go-job-board-api/others"

	"github.com/gin-gonic/gin"
)

func PostJob(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method        := others.GetStr(v_data, "method")
	v_title         := others.GetStr(v_data, "title")
	v_description   := others.GetStr(v_data, "description")
	v_company       := others.GetStr(v_data, "company")
	v_location      := others.GetStr(v_data, "location")
	v_job_type      := others.GetStr(v_data, "job_type")
	v_salary_min    := others.GetStr(v_data, "salary_min")
	v_salary_max    := others.GetStr(v_data, "salary_max")
	v_expired_date  := others.GetStr(v_data, "expired_date")
	v_session_key   := others.GetStr(v_data, "session_key")
	v_datetime      := others.GetStr(v_data, "datetime")
	v_checksum      := others.GetStr(v_data, "checksum")
	v_skills        := others.GetSliceStr(v_data, "skills")

	if v_title == "" || v_company == "" || v_location == "" || v_job_type == "" || v_expired_date == "" || v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	if v_job_type != "full_time" && v_job_type != "part_time" && v_job_type != "contract" && v_job_type != "internship" && v_job_type != "remote" {
		others.ResponseService(p_c, v_method, 422, "Job type harus full_time, part_time, contract, internship, atau remote", nil)
		return
	}

	v_salary_min_float, _  := strconv.ParseFloat(v_salary_min, 64)
	v_salary_max_float, _  := strconv.ParseFloat(v_salary_max, 64)

	v_app_payload   := v_method + "#" + v_title + "#" + v_company + "#" + v_job_type + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	v_user_id   := v_session_data["user_id"].(string)
	v_hasil, _  := dao.InsertJob(v_user_id, v_title, v_description, v_company, v_location, v_job_type, v_expired_date, v_salary_min_float, v_salary_max_float, v_skills)

	if v_hasil["status"] == "T" {
		dao.CacheDeletePattern("job:list:*")
		v_result := v_hasil["result"].(map[string]interface{})
		others.ResponseService(p_c, v_method, 200, "Success", map[string]interface{}{
			"job_id": v_result["job_id"],
		})
	} else {
		others.ResponseService(p_c, v_method, 400, v_hasil["message"].(string), nil)
	}
}

func GetJobList(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method        := others.GetStr(v_data, "method")
	v_keyword       := others.GetStr(v_data, "keyword")
	v_location      := others.GetStr(v_data, "location")
	v_job_type      := others.GetStr(v_data, "job_type")
	v_page          := others.GetStr(v_data, "page")
	v_limit         := others.GetStr(v_data, "limit")
	v_session_key   := others.GetStr(v_data, "session_key")
	v_datetime      := others.GetStr(v_data, "datetime")
	v_checksum      := others.GetStr(v_data, "checksum")

	if v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	v_app_payload   := v_method + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	if v_page == "" {
		v_page = "1"
	}
	if v_limit == "" {
		v_limit = "20"
	}

	_, v_limit_int, v_offset_int := others.ParsePagination(v_page, v_limit)
	v_hasil, _ := dao.GetJobList(v_keyword, v_location, v_job_type, v_limit_int, v_offset_int)

	if v_hasil["status"] == "T" {
		others.ResponseService(p_c, v_method, 200, "Success", map[string]interface{}{
			"data": v_hasil["result"],
		})
	} else {
		others.ResponseService(p_c, v_method, 400, v_hasil["message"].(string), nil)
	}
}

func GetJobDetail(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method        := others.GetStr(v_data, "method")
	v_job_id        := others.GetStr(v_data, "job_id")
	v_session_key   := others.GetStr(v_data, "session_key")
	v_datetime      := others.GetStr(v_data, "datetime")
	v_checksum      := others.GetStr(v_data, "checksum")

	if v_job_id == "" || v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	v_app_payload   := v_method + "#" + v_job_id + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	v_hasil, _ := dao.GetJobDetail(v_job_id)

	if v_hasil["status_code"] == 204 {
		others.ResponseService(p_c, v_method, 204, "Data tidak ditemukan", nil)
	} else if v_hasil["status"] == "T" {
		others.ResponseService(p_c, v_method, 200, "Success", map[string]interface{}{
			"data": v_hasil["result"],
		})
	} else {
		others.ResponseService(p_c, v_method, 400, v_hasil["message"].(string), nil)
	}
}

func UpdateJob(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method        := others.GetStr(v_data, "method")
	v_job_id        := others.GetStr(v_data, "job_id")
	v_session_key   := others.GetStr(v_data, "session_key")
	v_datetime      := others.GetStr(v_data, "datetime")
	v_checksum      := others.GetStr(v_data, "checksum")

	if v_job_id == "" || v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	v_app_payload   := v_method + "#" + v_job_id + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	v_update := map[string]interface{}{}
	if v_title := others.GetStr(v_data, "title"); v_title != "" {
		v_update["title"] = v_title
	}
	if v_desc := others.GetStr(v_data, "description"); v_desc != "" {
		v_update["description"] = v_desc
	}
	if v_loc := others.GetStr(v_data, "location"); v_loc != "" {
		v_update["location"] = v_loc
	}
	if v_sm := others.GetStr(v_data, "salary_min"); v_sm != "" {
		if v_sm_float, v_err := strconv.ParseFloat(v_sm, 64); v_err == nil {
			v_update["salary_min"] = v_sm_float
		}
	}
	if v_sx := others.GetStr(v_data, "salary_max"); v_sx != "" {
		if v_sx_float, v_err := strconv.ParseFloat(v_sx, 64); v_err == nil {
			v_update["salary_max"] = v_sx_float
		}
	}
	if v_ed := others.GetStr(v_data, "expired_date"); v_ed != "" {
		v_update["expired_date"] = v_ed
	}

	if len(v_update) == 0 {
		others.ResponseService(p_c, v_method, 422, "Tidak ada data yang diupdate", nil)
		return
	}

	v_user_id   := v_session_data["user_id"].(string)
	v_hasil, _  := dao.UpdateJob(v_job_id, v_user_id, v_update)

	if v_hasil["status"] == "T" {
		dao.CacheDelete("job:detail:" + v_job_id)
		dao.CacheDelete("job:summary:" + v_job_id)
		dao.CacheDeletePattern("job:list:*")
		others.ResponseService(p_c, v_method, 200, "Success", nil)
	} else {
		others.ResponseService(p_c, v_method, 400, v_hasil["message"].(string), nil)
	}
}

func CloseJob(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method        := others.GetStr(v_data, "method")
	v_job_id        := others.GetStr(v_data, "job_id")
	v_session_key   := others.GetStr(v_data, "session_key")
	v_datetime      := others.GetStr(v_data, "datetime")
	v_checksum      := others.GetStr(v_data, "checksum")

	if v_job_id == "" || v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	v_app_payload   := v_method + "#" + v_job_id + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	v_user_id   := v_session_data["user_id"].(string)
	v_hasil, _  := dao.CloseJob(v_job_id, v_user_id)

	if v_hasil["status"] == "T" {
		dao.CacheDelete("job:detail:" + v_job_id)
		dao.CacheDelete("job:summary:" + v_job_id)
		dao.CacheDeletePattern("job:list:*")
		others.ResponseService(p_c, v_method, 200, "Success", nil)
	} else {
		others.ResponseService(p_c, v_method, 400, v_hasil["message"].(string), nil)
	}
}

func ApplyJob(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method        := others.GetStr(v_data, "method")
	v_job_id        := others.GetStr(v_data, "job_id")
	v_cover_letter  := others.GetStr(v_data, "cover_letter")
	v_session_key   := others.GetStr(v_data, "session_key")
	v_datetime      := others.GetStr(v_data, "datetime")
	v_checksum      := others.GetStr(v_data, "checksum")
	v_skills        := others.GetSliceStr(v_data, "skills")

	if v_job_id == "" || v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	v_app_payload   := v_method + "#" + v_job_id + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	v_user_id   := v_session_data["user_id"].(string)
	v_hasil, _  := dao.InsertApplication(v_job_id, v_user_id, v_cover_letter, v_skills)

	v_status_code := v_hasil["status_code"].(int)
	if v_status_code == 200 {
		dao.CacheDeletePattern("application:my:" + v_user_id + ":*")
		dao.CacheDelete("job:summary:" + v_job_id)
		v_result := v_hasil["result"].(map[string]interface{})
		others.ResponseService(p_c, v_method, 200, "Success", map[string]interface{}{
			"application_id": v_result["application_id"],
		})
	} else {
		others.ResponseService(p_c, v_method, v_status_code, v_hasil["message"].(string), nil)
	}
}

func UpdateApplicationStatus(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method            := others.GetStr(v_data, "method")
	v_application_id    := others.GetStr(v_data, "application_id")
	v_status            := others.GetStr(v_data, "status")
	v_note              := others.GetStr(v_data, "note")
	v_session_key       := others.GetStr(v_data, "session_key")
	v_datetime          := others.GetStr(v_data, "datetime")
	v_checksum          := others.GetStr(v_data, "checksum")

	if v_application_id == "" || v_status == "" || v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	if v_status != "reviewed" && v_status != "accepted" && v_status != "rejected" {
		others.ResponseService(p_c, v_method, 422, "Status harus reviewed, accepted, atau rejected", nil)
		return
	}

	v_app_payload   := v_method + "#" + v_application_id + "#" + v_status + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	v_user_id   := v_session_data["user_id"].(string)
	v_hasil, _  := dao.UpdateApplicationStatus(v_application_id, v_user_id, v_status, v_note)

	if v_hasil["status"] == "T" {
		others.ResponseService(p_c, v_method, 200, "Success", nil)
	} else {
		others.ResponseService(p_c, v_method, 400, v_hasil["message"].(string), nil)
	}
}

func GetMyApplication(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method        := others.GetStr(v_data, "method")
	v_status        := others.GetStr(v_data, "status")
	v_session_key   := others.GetStr(v_data, "session_key")
	v_datetime      := others.GetStr(v_data, "datetime")
	v_checksum      := others.GetStr(v_data, "checksum")

	if v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	v_app_payload   := v_method + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	v_user_id   := v_session_data["user_id"].(string)
	v_hasil, _  := dao.GetMyApplication(v_user_id, v_status)

	if v_hasil["status"] == "T" {
		others.ResponseService(p_c, v_method, 200, "Success", map[string]interface{}{
			"data": v_hasil["result"],
		})
	} else {
		others.ResponseService(p_c, v_method, 400, v_hasil["message"].(string), nil)
	}
}

func GetJobApplicants(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method        := others.GetStr(v_data, "method")
	v_job_id        := others.GetStr(v_data, "job_id")
	v_status        := others.GetStr(v_data, "status")
	v_session_key   := others.GetStr(v_data, "session_key")
	v_datetime      := others.GetStr(v_data, "datetime")
	v_checksum      := others.GetStr(v_data, "checksum")

	if v_job_id == "" || v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	v_app_payload   := v_method + "#" + v_job_id + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	v_user_id   := v_session_data["user_id"].(string)
	v_hasil, _  := dao.GetJobApplicants(v_job_id, v_user_id, v_status)

	if v_hasil["status"] == "T" {
		others.ResponseService(p_c, v_method, 200, "Success", map[string]interface{}{
			"data": v_hasil["result"],
		})
	} else {
		others.ResponseService(p_c, v_method, 400, v_hasil["message"].(string), nil)
	}
}

func GetMatchedJobs(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method        := others.GetStr(v_data, "method")
	v_session_key   := others.GetStr(v_data, "session_key")
	v_datetime      := others.GetStr(v_data, "datetime")
	v_checksum      := others.GetStr(v_data, "checksum")
	v_skills        := others.GetSliceStr(v_data, "skills")

	if len(v_skills) == 0 || v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	v_app_payload   := v_method + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	v_hasil, _ := dao.GetMatchedJobs(v_skills, 100, 0)

	if v_hasil["status"] == "T" {
		others.ResponseService(p_c, v_method, 200, "Success", map[string]interface{}{
			"data": v_hasil["result"],
		})
	} else {
		others.ResponseService(p_c, v_method, 400, v_hasil["message"].(string), nil)
	}
}

func GetJobSummary(p_c *gin.Context) {
	v_body, _   := p_c.Get("parsed_body")
	v_data       := v_body.(map[string]interface{})

	v_method        := others.GetStr(v_data, "method")
	v_job_id        := others.GetStr(v_data, "job_id")
	v_session_key   := others.GetStr(v_data, "session_key")
	v_datetime      := others.GetStr(v_data, "datetime")
	v_checksum      := others.GetStr(v_data, "checksum")

	if v_job_id == "" || v_session_key == "" || v_datetime == "" || v_checksum == "" {
		others.ResponseService(p_c, v_method, 422, "Invalid Request Data", nil)
		return
	}

	v_app_payload   := v_method + "#" + v_job_id + "#" + v_datetime + "#" + os.Getenv("SECRET_KEY_REQUEST")
	v_app_checksum  := others.GenerateChecksum(v_app_payload)
	if v_checksum != v_app_checksum {
		others.ResponseService(p_c, v_method, 406, "Invalid Key", nil)
		return
	}

	v_session_data, v_err := dao.CacheGet("session_key:" + v_session_key)
	if v_err != nil || v_session_data == nil {
		others.ResponseService(p_c, v_method, 401, "Session tidak valid atau sudah expired", nil)
		return
	}

	v_user_id   := v_session_data["user_id"].(string)
	v_hasil, _  := dao.GetJobSummary(v_job_id, v_user_id)

	if v_hasil["status_code"] == 204 {
		others.ResponseService(p_c, v_method, 204, "Data tidak ditemukan", nil)
	} else if v_hasil["status"] == "T" {
		others.ResponseService(p_c, v_method, 200, "Success", map[string]interface{}{
			"data": v_hasil["result"],
		})
	} else {
		others.ResponseService(p_c, v_method, 400, v_hasil["message"].(string), nil)
	}
}
