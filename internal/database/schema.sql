PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS employer (
  id   INTEGER PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS skill (
  id       INTEGER PRIMARY KEY,
  name     TEXT NOT NULL UNIQUE,
  category TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS job (
  id                   INTEGER PRIMARY KEY,
  source_id            INTEGER,
  title                TEXT NOT NULL,
  employer_id          INTEGER NOT NULL REFERENCES employer(id),
  location             TEXT NOT NULL,
  employment_type      TEXT NOT NULL CHECK (employment_type IN ('full_time','part_time','contract','internship','consultancy','temporary','volunteer')),
  role_family          TEXT NOT NULL CHECK (role_family IN ('software','data','gis','infrastructure','support','security','management','other')),
  industry             TEXT,
  seniority            TEXT NOT NULL CHECK (seniority IN ('intern','junior','mid','senior','manager','director','executive')),
  is_ict               INTEGER NOT NULL DEFAULT 1 CHECK (is_ict IN (0,1)),
  education_level      TEXT CHECK (education_level IS NULL OR education_level IN ('secondary','certificate','diploma','bachelor','master','phd')),
  experience_min_years INTEGER,
  experience_max_years INTEGER,
  salary_min           REAL,
  salary_max           REAL,
  salary_currency      TEXT,
  salary_period        TEXT CHECK (salary_period IN ('hour','day','week','month','year')),
  posted_on            TEXT NOT NULL,
  deadline             TEXT,
  summary              TEXT NOT NULL,
  responsibilities     TEXT NOT NULL DEFAULT '[]',
  requirements         TEXT NOT NULL DEFAULT '[]',
  created_at           TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ','now')),
  UNIQUE (title, employer_id, posted_on)
);

CREATE TABLE IF NOT EXISTS job_skill (
  job_id      INTEGER NOT NULL REFERENCES job(id) ON DELETE CASCADE,
  skill_id    INTEGER NOT NULL REFERENCES skill(id),
  requirement TEXT NOT NULL CHECK (requirement IN ('required','preferred','mentioned')),
  evidence    TEXT,
  PRIMARY KEY (job_id, skill_id)
);

CREATE INDEX IF NOT EXISTS idx_job_posted ON job(posted_on DESC);
CREATE INDEX IF NOT EXISTS idx_job_role ON job(role_family, posted_on DESC);
CREATE INDEX IF NOT EXISTS idx_job_skill_skill ON job_skill(skill_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_job_source_title ON job(source_id, title) WHERE source_id IS NOT NULL;
