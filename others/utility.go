package others

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

var TZ_JAKARTA *time.Location

func init() {
	var v_err error
	TZ_JAKARTA, v_err = time.LoadLocation("Asia/Jakarta")
	if v_err != nil {
		TZ_JAKARTA = time.FixedZone("WIB", 7*3600)
	}
}

func GetTTLUntilMidnightJakarta() int {
	v_now      := time.Now().In(TZ_JAKARTA)
	v_midnight := time.Date(v_now.Year(), v_now.Month(), v_now.Day()+1, 0, 0, 0, 0, TZ_JAKARTA)
	return int(v_midnight.Sub(v_now).Seconds())
}

func ValidateRequestDatetime(p_datetime string) (bool, string) {
	v_layout         := "2006-01-02 15:04:05"
	v_request_dt, v_err := time.ParseInLocation(v_layout, p_datetime, TZ_JAKARTA)
	if v_err != nil {
		return false, "Format datetime tidak valid. Gunakan YYYY-MM-DD HH:MM:SS"
	}
	v_now_jakarta   := time.Now().In(TZ_JAKARTA)
	v_diff          := v_now_jakarta.Sub(v_request_dt)
	v_abs_diff      := time.Duration(math.Abs(float64(v_diff)))
	if v_abs_diff > 24*time.Hour {
		return false, "Datetime request harus dalam rentang ±24 jam dari waktu sekarang"
	}
	return true, "OK"
}

func NowJakartaStr() string {
	return time.Now().In(TZ_JAKARTA).Format("2006-01-02 15:04:05")
}

func ParsePagination(p_page, p_limit string) (int, int, int) {
	v_page, v_err := strconv.Atoi(p_page)
	if v_err != nil || v_page < 1 {
		v_page = 1
	}
	v_limit, v_err := strconv.Atoi(p_limit)
	if v_err != nil || v_limit < 1 {
		v_limit = 20
	}
	if v_limit > 100 {
		v_limit = 100
	}
	v_offset := (v_page - 1) * v_limit
	return v_page, v_limit, v_offset
}

func IntToStr(p_n int) string {
	return fmt.Sprintf("%d", p_n)
}

func CalcMatchScore(p_job_skills, p_applicant_skills []string) int {
	if len(p_job_skills) == 0 {
		return 0
	}
	v_job_map := map[string]bool{}
	for _, v_s := range p_job_skills {
		v_job_map[strings.ToLower(strings.TrimSpace(v_s))] = true
	}
	v_matched := 0
	for _, v_s := range p_applicant_skills {
		if v_job_map[strings.ToLower(strings.TrimSpace(v_s))] {
			v_matched++
		}
	}
	v_score := int(math.Round(float64(v_matched) / float64(len(p_job_skills)) * 100))
	return v_score
}
