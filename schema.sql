-- ============================================================
-- JOB BOARD API - DATABASE SCHEMA
-- ============================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    user_id     UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username    VARCHAR(100) NOT NULL UNIQUE,
    email       VARCHAR(150) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'),
    is_deleted  BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE jobs (
    job_id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    poster_id       UUID NOT NULL REFERENCES users(user_id),
    title           VARCHAR(200) NOT NULL,
    description     TEXT,
    company         VARCHAR(150) NOT NULL,
    location        VARCHAR(150) NOT NULL,
    job_type        VARCHAR(20) NOT NULL CHECK (job_type IN ('full_time', 'part_time', 'contract', 'internship', 'remote')),
    salary_min      NUMERIC(15, 2),
    salary_max      NUMERIC(15, 2),
    status          VARCHAR(10) NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'closed')),
    expired_date    DATE NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'),
    updated_at      TIMESTAMP,
    is_deleted      BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE job_skills (
    skill_id    UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_id      UUID NOT NULL REFERENCES jobs(job_id),
    skill       VARCHAR(100) NOT NULL
);

CREATE TABLE applications (
    application_id  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_id          UUID NOT NULL REFERENCES jobs(job_id),
    applicant_id    UUID NOT NULL REFERENCES users(user_id),
    cover_letter    TEXT,
    skills          TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'reviewed', 'accepted', 'rejected')),
    note            TEXT,
    created_at      TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'),
    updated_at      TIMESTAMP,
    UNIQUE (job_id, applicant_id)
);

CREATE INDEX idx_jobs_status          ON jobs(status, expired_date) WHERE is_deleted = false;
CREATE INDEX idx_jobs_poster          ON jobs(poster_id) WHERE is_deleted = false;
CREATE INDEX idx_job_skills_job       ON job_skills(job_id);
CREATE INDEX idx_applications_job     ON applications(job_id);
CREATE INDEX idx_applications_user    ON applications(applicant_id);
