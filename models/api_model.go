package models

type PostJobRequest struct {
	Method       string   `json:"method"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Company      string   `json:"company"`
	Location     string   `json:"location"`
	JobType      string   `json:"job_type"`
	SalaryMin    string   `json:"salary_min"`
	SalaryMax    string   `json:"salary_max"`
	Skills       []string `json:"skills"`
	ExpiredDate  string   `json:"expired_date"`
	SessionKey   string   `json:"session_key"`
	Datetime     string   `json:"datetime"`
	Checksum     string   `json:"checksum"`
}

type GetJobListRequest struct {
	Method     string `json:"method"`
	Keyword    string `json:"keyword"`
	Location   string `json:"location"`
	JobType    string `json:"job_type"`
	Page       string `json:"page"`
	Limit      string `json:"limit"`
	SessionKey string `json:"session_key"`
	Datetime   string `json:"datetime"`
	Checksum   string `json:"checksum"`
}

type GetJobDetailRequest struct {
	Method     string `json:"method"`
	JobID      string `json:"job_id"`
	SessionKey string `json:"session_key"`
	Datetime   string `json:"datetime"`
	Checksum   string `json:"checksum"`
}

type UpdateJobRequest struct {
	Method      string `json:"method"`
	JobID       string `json:"job_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	SalaryMin   string `json:"salary_min"`
	SalaryMax   string `json:"salary_max"`
	ExpiredDate string `json:"expired_date"`
	SessionKey  string `json:"session_key"`
	Datetime    string `json:"datetime"`
	Checksum    string `json:"checksum"`
}

type CloseJobRequest struct {
	Method     string `json:"method"`
	JobID      string `json:"job_id"`
	SessionKey string `json:"session_key"`
	Datetime   string `json:"datetime"`
	Checksum   string `json:"checksum"`
}

type ApplyJobRequest struct {
	Method      string   `json:"method"`
	JobID       string   `json:"job_id"`
	CoverLetter string   `json:"cover_letter"`
	Skills      []string `json:"skills"`
	SessionKey  string   `json:"session_key"`
	Datetime    string   `json:"datetime"`
	Checksum    string   `json:"checksum"`
}

type UpdateApplicationStatusRequest struct {
	Method        string `json:"method"`
	ApplicationID string `json:"application_id"`
	Status        string `json:"status"`
	Note          string `json:"note"`
	SessionKey    string `json:"session_key"`
	Datetime      string `json:"datetime"`
	Checksum      string `json:"checksum"`
}

type GetMyApplicationRequest struct {
	Method     string `json:"method"`
	Status     string `json:"status"`
	SessionKey string `json:"session_key"`
	Datetime   string `json:"datetime"`
	Checksum   string `json:"checksum"`
}

type GetJobApplicantsRequest struct {
	Method     string `json:"method"`
	JobID      string `json:"job_id"`
	Status     string `json:"status"`
	SessionKey string `json:"session_key"`
	Datetime   string `json:"datetime"`
	Checksum   string `json:"checksum"`
}

type GetMatchedJobsRequest struct {
	Method     string   `json:"method"`
	Skills     []string `json:"skills"`
	SessionKey string   `json:"session_key"`
	Datetime   string   `json:"datetime"`
	Checksum   string   `json:"checksum"`
}

type GetJobSummaryRequest struct {
	Method     string `json:"method"`
	JobID      string `json:"job_id"`
	SessionKey string `json:"session_key"`
	Datetime   string `json:"datetime"`
	Checksum   string `json:"checksum"`
}
