package dao

import (
	"context"
	"fmt"
	"strings"

	"go-job-board-api/config"
	"go-job-board-api/others"
)

func ResponseJson(p_status_code int, p_flag string, p_message string, p_result interface{}) map[string]interface{} {
	return map[string]interface{}{
		"status_code": p_status_code,
		"status":      p_flag,
		"message":     p_message,
		"result":      p_result,
	}
}

func InsertJob(p_poster_id, p_title, p_description, p_company, p_location, p_job_type, p_expired_date string, p_salary_min, p_salary_max float64, p_skills []string) (map[string]interface{}, error) {
	v_pool  := config.GetPgPool()
	v_ctx   := context.Background()

	v_tx, v_err := v_pool.Begin(v_ctx)
	if v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}
	defer v_tx.Rollback(v_ctx)

	v_query := `
		INSERT INTO jobs
		(poster_id, title, description, company, location, job_type, salary_min, salary_max, expired_date, created_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')
		RETURNING job_id
	`
	v_row := v_tx.QueryRow(v_ctx, v_query, p_poster_id, p_title, p_description, p_company, p_location, p_job_type, p_salary_min, p_salary_max, p_expired_date)
	var v_job_id string
	if v_err := v_row.Scan(&v_job_id); v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}

	for _, v_skill := range p_skills {
		v_skill_query := `
			INSERT INTO job_skills (job_id, skill)
			VALUES($1, $2)
		`
		v_tx.Exec(v_ctx, v_skill_query, v_job_id, strings.TrimSpace(v_skill))
	}

	if v_err := v_tx.Commit(v_ctx); v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}
	return ResponseJson(200, "T", "Success", map[string]interface{}{"job_id": v_job_id}), nil
}

func GetJobList(p_keyword, p_location, p_job_type string, p_limit, p_offset int) (map[string]interface{}, error) {
	v_cache_key     := fmt.Sprintf("job:list:%s:%s:%s:%d:%d", p_keyword, p_location, p_job_type, p_offset, p_limit)
	v_cached, v_err := CacheGetList(v_cache_key)
	if v_err == nil {
		return ResponseJson(200, "T", "Success", v_cached), nil
	}

	v_pool      := config.GetPgPool()
	v_ctx       := context.Background()
	v_where     := "WHERE j.status = 'open' AND j.expired_date >= CURRENT_DATE AND j.is_deleted = false"
	v_args      := []interface{}{}
	v_idx       := 1

	if p_keyword != "" {
		v_where += fmt.Sprintf(" AND (j.title ILIKE $%d OR j.description ILIKE $%d OR j.company ILIKE $%d)", v_idx, v_idx, v_idx)
		v_args  = append(v_args, "%"+p_keyword+"%")
		v_idx++
	}
	if p_location != "" {
		v_where += fmt.Sprintf(" AND j.location ILIKE $%d", v_idx)
		v_args  = append(v_args, "%"+p_location+"%")
		v_idx++
	}
	if p_job_type != "" {
		v_where += fmt.Sprintf(" AND j.job_type = $%d", v_idx)
		v_args  = append(v_args, p_job_type)
		v_idx++
	}

	v_args  = append(v_args, p_limit, p_offset)
	v_query := fmt.Sprintf(`
		SELECT
			j.job_id::varchar,
			j.title,
			j.company,
			j.location,
			j.job_type,
			coalesce(j.salary_min::varchar, '') salary_min,
			coalesce(j.salary_max::varchar, '') salary_max,
			j.status,
			to_char(j.expired_date, 'YYYY-MM-DD') expired_date,
			to_char(j.created_at, 'YYYY-MM-DD HH24:MI:SS') created_at,
			COUNT(a.application_id)::varchar total_applicant
		FROM jobs j
		LEFT JOIN applications a ON a.job_id = j.job_id
		%s
		GROUP BY j.job_id, j.title, j.company, j.location, j.job_type, j.salary_min, j.salary_max, j.status, j.expired_date, j.created_at
		ORDER BY j.created_at DESC
		LIMIT $%d OFFSET $%d
	`, v_where, v_idx, v_idx+1)

	v_rows, v_err := v_pool.Query(v_ctx, v_query, v_args...)
	if v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}
	defer v_rows.Close()

	v_result := []map[string]interface{}{}
	for v_rows.Next() {
		var (
			v_job_id, v_title, v_company, v_location, v_job_type string
			v_salary_min, v_salary_max, v_status                 string
			v_expired_date, v_created_at, v_total_applicant      string
		)
		if v_err := v_rows.Scan(&v_job_id, &v_title, &v_company, &v_location, &v_job_type, &v_salary_min, &v_salary_max, &v_status, &v_expired_date, &v_created_at, &v_total_applicant); v_err != nil {
			continue
		}
		v_result = append(v_result, map[string]interface{}{
			"job_id":          v_job_id,
			"title":           v_title,
			"company":         v_company,
			"location":        v_location,
			"job_type":        v_job_type,
			"salary_min":      v_salary_min,
			"salary_max":      v_salary_max,
			"status":          v_status,
			"expired_date":    v_expired_date,
			"created_at":      v_created_at,
			"total_applicant": v_total_applicant,
		})
	}

	if len(v_result) > 0 {
		CacheSetList(v_cache_key, v_result, 120)
	}
	return ResponseJson(200, "T", "Success", v_result), nil
}

func GetJobDetail(p_job_id string) (map[string]interface{}, error) {
	v_cache_key     := fmt.Sprintf("job:detail:%s", p_job_id)
	v_cached, v_err := CacheGet(v_cache_key)
	if v_err == nil {
		return ResponseJson(200, "T", "Success", v_cached), nil
	}

	v_pool  := config.GetPgPool()
	v_ctx   := context.Background()
	v_query := `
		SELECT
			j.job_id::varchar,
			j.poster_id::varchar,
			j.title,
			coalesce(j.description, '') description,
			j.company,
			j.location,
			j.job_type,
			coalesce(j.salary_min::varchar, '') salary_min,
			coalesce(j.salary_max::varchar, '') salary_max,
			j.status,
			to_char(j.expired_date, 'YYYY-MM-DD') expired_date,
			to_char(j.created_at, 'YYYY-MM-DD HH24:MI:SS') created_at,
			to_char(j.updated_at, 'YYYY-MM-DD HH24:MI:SS') updated_at
		FROM jobs j
		WHERE j.job_id = $1
		AND j.is_deleted = false
	`
	v_row := v_pool.QueryRow(v_ctx, v_query, p_job_id)
	var (
		v_job_id, v_poster_id, v_title, v_description, v_company, v_location, v_job_type string
		v_salary_min, v_salary_max, v_status, v_expired_date, v_created_at, v_updated_at string
	)
	if v_err := v_row.Scan(&v_job_id, &v_poster_id, &v_title, &v_description, &v_company, &v_location, &v_job_type, &v_salary_min, &v_salary_max, &v_status, &v_expired_date, &v_created_at, &v_updated_at); v_err != nil {
		return ResponseJson(204, "T", "Data tidak ditemukan", nil), nil
	}

	v_skill_query := `
		SELECT skill FROM job_skills WHERE job_id = $1 ORDER BY skill ASC
	`
	v_skill_rows, _ := v_pool.Query(v_ctx, v_skill_query, p_job_id)
	defer v_skill_rows.Close()
	v_skills := []string{}
	for v_skill_rows.Next() {
		var v_skill string
		if v_err := v_skill_rows.Scan(&v_skill); v_err == nil {
			v_skills = append(v_skills, v_skill)
		}
	}

	v_data := map[string]interface{}{
		"job_id":       v_job_id,
		"poster_id":    v_poster_id,
		"title":        v_title,
		"description":  v_description,
		"company":      v_company,
		"location":     v_location,
		"job_type":     v_job_type,
		"salary_min":   v_salary_min,
		"salary_max":   v_salary_max,
		"status":       v_status,
		"skills":       v_skills,
		"expired_date": v_expired_date,
		"created_at":   v_created_at,
		"updated_at":   v_updated_at,
	}
	CacheSet(v_cache_key, v_data, 120)
	return ResponseJson(200, "T", "Success", v_data), nil
}

func UpdateJob(p_job_id, p_poster_id string, p_data map[string]interface{}) (map[string]interface{}, error) {
	v_pool  := config.GetPgPool()
	v_ctx   := context.Background()
	v_set   := ""
	v_args  := []interface{}{}
	v_idx   := 1

	for v_k, v_v := range p_data {
		if v_set != "" {
			v_set += ", "
		}
		v_set += fmt.Sprintf("%s = $%d", v_k, v_idx)
		v_args = append(v_args, v_v)
		v_idx++
	}
	v_set  += ", updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'"
	v_args  = append(v_args, p_job_id, p_poster_id)

	v_query := fmt.Sprintf(`
		UPDATE jobs SET %s
		WHERE job_id = $%d AND poster_id = $%d AND is_deleted = false
	`, v_set, v_idx, v_idx+1)

	_, v_err := v_pool.Exec(v_ctx, v_query, v_args...)
	if v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}
	return ResponseJson(200, "T", "Success", nil), nil
}

func CloseJob(p_job_id, p_poster_id string) (map[string]interface{}, error) {
	v_pool  := config.GetPgPool()
	v_ctx   := context.Background()
	v_query := `
		UPDATE jobs
		SET status = 'closed', updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'
		WHERE job_id = $1 AND poster_id = $2 AND status = 'open' AND is_deleted = false
	`
	_, v_err := v_pool.Exec(v_ctx, v_query, p_job_id, p_poster_id)
	if v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}
	return ResponseJson(200, "T", "Success", nil), nil
}

func InsertApplication(p_job_id, p_applicant_id, p_cover_letter string, p_skills []string) (map[string]interface{}, error) {
	v_pool  := config.GetPgPool()
	v_ctx   := context.Background()

	v_check_query := `
		SELECT COUNT(*) FROM applications
		WHERE job_id = $1 AND applicant_id = $2
	`
	v_row := v_pool.QueryRow(v_ctx, v_check_query, p_job_id, p_applicant_id)
	var v_count int
	v_row.Scan(&v_count)
	if v_count > 0 {
		return ResponseJson(409, "F", "Kamu sudah melamar pekerjaan ini", nil), fmt.Errorf("already applied")
	}

	v_job_query := `
		SELECT status, expired_date FROM jobs
		WHERE job_id = $1 AND is_deleted = false
	`
	v_job_row := v_pool.QueryRow(v_ctx, v_job_query, p_job_id)
	var v_job_status, v_expired_date string
	if v_err := v_job_row.Scan(&v_job_status, &v_expired_date); v_err != nil {
		return ResponseJson(404, "F", "Lowongan tidak ditemukan", nil), v_err
	}
	if v_job_status != "open" {
		return ResponseJson(400, "F", "Lowongan sudah ditutup", nil), fmt.Errorf("job closed")
	}

	v_query := `
		INSERT INTO applications (job_id, applicant_id, cover_letter, skills, created_at)
		VALUES($1, $2, $3, $4, CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')
		RETURNING application_id
	`
	v_app_row := v_pool.QueryRow(v_ctx, v_query, p_job_id, p_applicant_id, p_cover_letter, strings.Join(p_skills, ","))
	var v_application_id string
	if v_err := v_app_row.Scan(&v_application_id); v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}
	return ResponseJson(200, "T", "Success", map[string]interface{}{"application_id": v_application_id}), nil
}

func UpdateApplicationStatus(p_application_id, p_poster_id, p_status, p_note string) (map[string]interface{}, error) {
	v_pool  := config.GetPgPool()
	v_ctx   := context.Background()
	v_query := `
		UPDATE applications a
		SET status = $1, note = $2,
		updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'
		FROM jobs j
		WHERE a.application_id = $3
		AND a.job_id = j.job_id
		AND j.poster_id = $4
	`
	_, v_err := v_pool.Exec(v_ctx, v_query, p_status, p_note, p_application_id, p_poster_id)
	if v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}
	return ResponseJson(200, "T", "Success", nil), nil
}

func GetMyApplication(p_applicant_id, p_status string) (map[string]interface{}, error) {
	v_cache_key     := fmt.Sprintf("application:my:%s:%s", p_applicant_id, p_status)
	v_cached, v_err := CacheGetList(v_cache_key)
	if v_err == nil {
		return ResponseJson(200, "T", "Success", v_cached), nil
	}

	v_pool      := config.GetPgPool()
	v_ctx       := context.Background()
	v_where     := "WHERE a.applicant_id = $1"
	v_args      := []interface{}{p_applicant_id}
	v_idx       := 2

	if p_status != "" {
		v_where += fmt.Sprintf(" AND a.status = $%d", v_idx)
		v_args  = append(v_args, p_status)
	}

	v_query := fmt.Sprintf(`
		SELECT
			a.application_id::varchar,
			j.job_id::varchar,
			j.title,
			j.company,
			j.location,
			j.job_type,
			a.status,
			coalesce(a.note, '') note,
			to_char(a.created_at, 'YYYY-MM-DD HH24:MI:SS') created_at,
			to_char(a.updated_at, 'YYYY-MM-DD HH24:MI:SS') updated_at
		FROM applications a
		JOIN jobs j ON j.job_id = a.job_id
		%s
		ORDER BY a.created_at DESC
	`, v_where)

	v_rows, v_err := v_pool.Query(v_ctx, v_query, v_args...)
	if v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}
	defer v_rows.Close()

	v_result := []map[string]interface{}{}
	for v_rows.Next() {
		var (
			v_application_id, v_job_id, v_title, v_company, v_location, v_job_type string
			v_status, v_note, v_created_at, v_updated_at                            string
		)
		if v_err := v_rows.Scan(&v_application_id, &v_job_id, &v_title, &v_company, &v_location, &v_job_type, &v_status, &v_note, &v_created_at, &v_updated_at); v_err != nil {
			continue
		}
		v_result = append(v_result, map[string]interface{}{
			"application_id": v_application_id,
			"job_id":         v_job_id,
			"title":          v_title,
			"company":        v_company,
			"location":       v_location,
			"job_type":       v_job_type,
			"status":         v_status,
			"note":           v_note,
			"created_at":     v_created_at,
			"updated_at":     v_updated_at,
		})
	}

	if len(v_result) > 0 {
		CacheSetList(v_cache_key, v_result, 60)
	}
	return ResponseJson(200, "T", "Success", v_result), nil
}

func GetJobApplicants(p_job_id, p_poster_id, p_status string) (map[string]interface{}, error) {
	v_pool      := config.GetPgPool()
	v_ctx       := context.Background()
	v_where     := "WHERE a.job_id = $1 AND j.poster_id = $2"
	v_args      := []interface{}{p_job_id, p_poster_id}
	v_idx       := 3

	if p_status != "" {
		v_where += fmt.Sprintf(" AND a.status = $%d", v_idx)
		v_args  = append(v_args, p_status)
	}

	v_query := fmt.Sprintf(`
		SELECT
			a.application_id::varchar,
			a.applicant_id::varchar,
			a.skills,
			a.status,
			coalesce(a.note, '') note,
			coalesce(a.cover_letter, '') cover_letter,
			to_char(a.created_at, 'YYYY-MM-DD HH24:MI:SS') created_at
		FROM applications a
		JOIN jobs j ON j.job_id = a.job_id
		%s
		ORDER BY a.created_at ASC
	`, v_where)

	v_rows, v_err := v_pool.Query(v_ctx, v_query, v_args...)
	if v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}
	defer v_rows.Close()

	v_job_skills_query := `SELECT skill FROM job_skills WHERE job_id = $1`
	v_skill_rows, _    := v_pool.Query(v_ctx, v_job_skills_query, p_job_id)
	defer v_skill_rows.Close()
	v_job_skills := []string{}
	for v_skill_rows.Next() {
		var v_skill string
		if v_err := v_skill_rows.Scan(&v_skill); v_err == nil {
			v_job_skills = append(v_job_skills, v_skill)
		}
	}

	v_result := []map[string]interface{}{}
	for v_rows.Next() {
		var (
			v_application_id, v_applicant_id, v_skills_str string
			v_status, v_note, v_cover_letter, v_created_at string
		)
		if v_err := v_rows.Scan(&v_application_id, &v_applicant_id, &v_skills_str, &v_status, &v_note, &v_cover_letter, &v_created_at); v_err != nil {
			continue
		}

		v_applicant_skills  := strings.Split(v_skills_str, ",")
		v_match_score       := others.CalcMatchScore(v_job_skills, v_applicant_skills)

		v_result = append(v_result, map[string]interface{}{
			"application_id": v_application_id,
			"applicant_id":   v_applicant_id,
			"skills":         v_applicant_skills,
			"match_score":    others.IntToStr(v_match_score),
			"status":         v_status,
			"note":           v_note,
			"cover_letter":   v_cover_letter,
			"created_at":     v_created_at,
		})
	}

	return ResponseJson(200, "T", "Success", v_result), nil
}

func GetMatchedJobs(p_applicant_skills []string, p_limit, p_offset int) (map[string]interface{}, error) {
	v_pool  := config.GetPgPool()
	v_ctx   := context.Background()
	v_query := `
		SELECT
			j.job_id::varchar,
			j.title,
			j.company,
			j.location,
			j.job_type,
			coalesce(j.salary_min::varchar, '') salary_min,
			coalesce(j.salary_max::varchar, '') salary_max,
			to_char(j.expired_date, 'YYYY-MM-DD') expired_date,
			ARRAY_AGG(js.skill ORDER BY js.skill) skills
		FROM jobs j
		JOIN job_skills js ON js.job_id = j.job_id
		WHERE j.status = 'open'
		AND j.expired_date >= CURRENT_DATE
		AND j.is_deleted = false
		GROUP BY j.job_id, j.title, j.company, j.location, j.job_type, j.salary_min, j.salary_max, j.expired_date
		ORDER BY j.created_at DESC
		LIMIT $1 OFFSET $2
	`
	v_rows, v_err := v_pool.Query(v_ctx, v_query, p_limit, p_offset)
	if v_err != nil {
		return ResponseJson(400, "F", v_err.Error(), nil), v_err
	}
	defer v_rows.Close()

	v_result := []map[string]interface{}{}
	for v_rows.Next() {
		var (
			v_job_id, v_title, v_company, v_location, v_job_type string
			v_salary_min, v_salary_max, v_expired_date            string
			v_skills                                              []string
		)
		if v_err := v_rows.Scan(&v_job_id, &v_title, &v_company, &v_location, &v_job_type, &v_salary_min, &v_salary_max, &v_expired_date, &v_skills); v_err != nil {
			continue
		}

		v_match_score := others.CalcMatchScore(v_skills, p_applicant_skills)
		if v_match_score == 0 {
			continue
		}

		v_result = append(v_result, map[string]interface{}{
			"job_id":       v_job_id,
			"title":        v_title,
			"company":      v_company,
			"location":     v_location,
			"job_type":     v_job_type,
			"salary_min":   v_salary_min,
			"salary_max":   v_salary_max,
			"expired_date": v_expired_date,
			"skills":       v_skills,
			"match_score":  others.IntToStr(v_match_score),
		})
	}

	v_sorted := sortByMatchScore(v_result)
	return ResponseJson(200, "T", "Success", v_sorted), nil
}

func sortByMatchScore(p_data []map[string]interface{}) []map[string]interface{} {
	for v_i := 0; v_i < len(p_data)-1; v_i++ {
		for v_j := v_i + 1; v_j < len(p_data); v_j++ {
			v_score_i := 0
			v_score_j := 0
			fmt.Sscanf(p_data[v_i]["match_score"].(string), "%d", &v_score_i)
			fmt.Sscanf(p_data[v_j]["match_score"].(string), "%d", &v_score_j)
			if v_score_j > v_score_i {
				p_data[v_i], p_data[v_j] = p_data[v_j], p_data[v_i]
			}
		}
	}
	return p_data
}

func GetJobSummary(p_job_id, p_poster_id string) (map[string]interface{}, error) {
	v_cache_key     := fmt.Sprintf("job:summary:%s", p_job_id)
	v_cached, v_err := CacheGet(v_cache_key)
	if v_err == nil {
		return ResponseJson(200, "T", "Success", v_cached), nil
	}

	v_pool  := config.GetPgPool()
	v_ctx   := context.Background()
	v_query := `
		SELECT
			j.title,
			j.status,
			to_char(j.expired_date, 'YYYY-MM-DD') expired_date,
			COUNT(a.application_id)::varchar                                                    total_applicant,
			COUNT(CASE WHEN a.status = 'pending' THEN 1 END)::varchar                          total_pending,
			COUNT(CASE WHEN a.status = 'reviewed' THEN 1 END)::varchar                         total_reviewed,
			COUNT(CASE WHEN a.status = 'accepted' THEN 1 END)::varchar                         total_accepted,
			COUNT(CASE WHEN a.status = 'rejected' THEN 1 END)::varchar                         total_rejected,
			to_char(j.created_at, 'YYYY-MM-DD HH24:MI:SS') created_at
		FROM jobs j
		LEFT JOIN applications a ON a.job_id = j.job_id
		WHERE j.job_id = $1
		AND j.poster_id = $2
		AND j.is_deleted = false
		GROUP BY j.title, j.status, j.expired_date, j.created_at
	`
	v_row := v_pool.QueryRow(v_ctx, v_query, p_job_id, p_poster_id)
	var (
		v_title, v_status, v_expired_date, v_created_at                                 string
		v_total_applicant, v_total_pending, v_total_reviewed, v_total_accepted, v_total_rejected string
	)
	if v_err := v_row.Scan(&v_title, &v_status, &v_expired_date, &v_total_applicant, &v_total_pending, &v_total_reviewed, &v_total_accepted, &v_total_rejected, &v_created_at); v_err != nil {
		return ResponseJson(204, "T", "Data tidak ditemukan", nil), nil
	}

	v_data := map[string]interface{}{
		"title":            v_title,
		"status":           v_status,
		"expired_date":     v_expired_date,
		"total_applicant":  v_total_applicant,
		"total_pending":    v_total_pending,
		"total_reviewed":   v_total_reviewed,
		"total_accepted":   v_total_accepted,
		"total_rejected":   v_total_rejected,
		"created_at":       v_created_at,
	}
	CacheSet(v_cache_key, v_data, 60)
	return ResponseJson(200, "T", "Success", v_data), nil
}
