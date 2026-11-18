package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/darlingson/chitsime/internal/database"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	if dir := filepath.Dir(path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create database directory: %w", err)
		}
	}

	dsn := "file:" + filepath.ToSlash(path) + "?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) Ping(ctx context.Context) error { return s.db.PingContext(ctx) }

func (s *Store) Setup(ctx context.Context) error {
	if err := s.execScript(ctx, database.Schema); err != nil {
		return fmt.Errorf("apply schema: %w", err)
	}
	if err := s.migrate(ctx); err != nil {
		return fmt.Errorf("migrate schema: %w", err)
	}
	return nil
}

func (s *Store) migrate(ctx context.Context) error {
	// Best-effort migrations for databases created before idempotent import support.
	if _, err := s.db.ExecContext(ctx, `ALTER TABLE job ADD COLUMN source_id INTEGER`); err != nil {
		msg := strings.ToLower(err.Error())
		if !strings.Contains(msg, "duplicate column") {
			return err
		}
	}
	if _, err := s.db.ExecContext(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS idx_job_source_title ON job(source_id, title) WHERE source_id IS NOT NULL`); err != nil {
		return err
	}
	return nil
}

func (s *Store) execScript(ctx context.Context, script string) error {
	for _, stmt := range strings.Split(script, ";") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := s.db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("%w\nstatement: %s", err, stmt)
		}
	}
	return nil
}

func (s *Store) CountJobs(ctx context.Context) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM job`).Scan(&n)
	return n, err
}

type Skill struct {
	Name        string
	Category    string
	Requirement string
	Evidence    string
}

type Job struct {
	ID               int64
	SourceID         *int64
	Title            string
	Employer         string
	Location         string
	EmploymentType   string
	RoleFamily       string
	Industry         *string
	Seniority        string
	IsICT            bool
	EducationLevel   *string
	ExperienceMin    *int
	ExperienceMax    *int
	SalaryMin        *float64
	SalaryMax        *float64
	SalaryCurrency   *string
	SalaryPeriod     *string
	PostedOn         string
	Deadline         *string
	Summary          string
	Responsibilities []string
	Requirements     []string
	Skills           []Skill
}

type JobFilter struct {
	RoleFamily string
	Skill      string
	PostedFrom string
	PostedTo   string
	Limit      int
	Offset     int
}

const jobSelect = `SELECT
	j.id, j.source_id, j.title, e.name, j.location, j.employment_type, j.role_family,
	COALESCE(j.industry, ''), j.seniority, j.is_ict, COALESCE(j.education_level, ''),
	j.experience_min_years, j.experience_max_years,
	j.salary_min, j.salary_max, COALESCE(j.salary_currency, ''), COALESCE(j.salary_period, ''),
	j.posted_on, COALESCE(j.deadline, ''), j.summary,
	j.responsibilities, j.requirements
FROM job j
JOIN employer e ON e.id = j.employer_id`

func (s *Store) GetJob(ctx context.Context, id int64) (Job, error) {
	row := s.db.QueryRowContext(ctx, jobSelect+` WHERE j.id = ?`, id)
	job, err := scanJob(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Job{}, ErrNotFound
	}
	if err != nil {
		return Job{}, err
	}
	job.Skills, err = s.skillsFor(ctx, id)
	if err != nil {
		return Job{}, err
	}
	return job, nil
}

func (s *Store) ListJobs(ctx context.Context, f JobFilter) ([]Job, error) {
	q := `SELECT j.id FROM job j WHERE 1 = 1`
	args := []any{}

	if f.RoleFamily != "" {
		q += ` AND j.role_family = ?`
		args = append(args, f.RoleFamily)
	}
	if f.Skill != "" {
		q += ` AND EXISTS (
			SELECT 1 FROM job_skill js
			JOIN skill s ON s.id = js.skill_id
			WHERE js.job_id = j.id AND s.name = ? COLLATE NOCASE
		)`
		args = append(args, f.Skill)
	}
	if f.PostedFrom != "" {
		q += ` AND j.posted_on >= ?`
		args = append(args, f.PostedFrom)
	}
	if f.PostedTo != "" {
		q += ` AND j.posted_on <= ?`
		args = append(args, f.PostedTo)
	}
	q += ` ORDER BY j.posted_on DESC, j.id LIMIT ? OFFSET ?`
	args = append(args, f.Limit, f.Offset)

	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()

	jobs := make([]Job, 0, len(ids))
	for _, id := range ids {
		job, err := s.GetJob(ctx, id)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (s *Store) skillsFor(ctx context.Context, jobID int64) ([]Skill, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT s.name, s.category, js.requirement, COALESCE(js.evidence, '')
		FROM job_skill js
		JOIN skill s ON s.id = js.skill_id
		WHERE js.job_id = ?
		ORDER BY CASE js.requirement
			WHEN 'required' THEN 0
			WHEN 'preferred' THEN 1
			ELSE 2
		END, s.category, s.name`, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []Skill{}
	for rows.Next() {
		var sk Skill
		if err := rows.Scan(&sk.Name, &sk.Category, &sk.Requirement, &sk.Evidence); err != nil {
			return nil, err
		}
		out = append(out, sk)
	}
	return out, rows.Err()
}

type SkillDemand struct {
	Name     string
	Category string
	Jobs     int
	Required int
}

func (s *Store) SkillDemand(ctx context.Context) ([]SkillDemand, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT s.name, s.category,
			COUNT(js.job_id) AS jobs,
			COALESCE(SUM(CASE WHEN js.requirement = 'required' THEN 1 ELSE 0 END), 0) AS required
		FROM skill s
		LEFT JOIN job_skill js ON js.skill_id = s.id
		GROUP BY s.id
		ORDER BY jobs DESC, s.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []SkillDemand{}
	for rows.Next() {
		var d SkillDemand
		if err := rows.Scan(&d.Name, &d.Category, &d.Jobs, &d.Required); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanJob(row rowScanner) (Job, error) {
	var (
		j                Job
		sourceID         sql.NullInt64
		isICT            int
		industry         string
		educationLevel   string
		salaryCurrency   string
		salaryPeriod     string
		deadline         string
		expMin, expMax   sql.NullInt64
		salaryMin        sql.NullFloat64
		salaryMax        sql.NullFloat64
		responsibilities string
		requirements     string
	)
	err := row.Scan(
		&j.ID, &sourceID, &j.Title, &j.Employer, &j.Location, &j.EmploymentType, &j.RoleFamily,
		&industry, &j.Seniority, &isICT, &educationLevel,
		&expMin, &expMax,
		&salaryMin, &salaryMax, &salaryCurrency, &salaryPeriod,
		&j.PostedOn, &deadline, &j.Summary,
		&responsibilities, &requirements,
	)
	if err != nil {
		return Job{}, err
	}

	if sourceID.Valid {
		v := sourceID.Int64
		j.SourceID = &v
	}
	j.IsICT = isICT == 1
	j.Industry = ptrString(industry)
	j.EducationLevel = ptrString(educationLevel)
	j.SalaryCurrency = ptrString(salaryCurrency)
	j.SalaryPeriod = ptrString(salaryPeriod)
	j.Deadline = ptrString(deadline)
	if expMin.Valid {
		v := int(expMin.Int64)
		j.ExperienceMin = &v
	}
	if expMax.Valid {
		v := int(expMax.Int64)
		j.ExperienceMax = &v
	}
	if salaryMin.Valid {
		v := salaryMin.Float64
		j.SalaryMin = &v
	}
	if salaryMax.Valid {
		v := salaryMax.Float64
		j.SalaryMax = &v
	}
	if err := json.Unmarshal([]byte(responsibilities), &j.Responsibilities); err != nil {
		return Job{}, fmt.Errorf("job %d responsibilities: %w", j.ID, err)
	}
	if err := json.Unmarshal([]byte(requirements), &j.Requirements); err != nil {
		return Job{}, fmt.Errorf("job %d requirements: %w", j.ID, err)
	}
	if j.Responsibilities == nil {
		j.Responsibilities = []string{}
	}
	if j.Requirements == nil {
		j.Requirements = []string{}
	}
	return j, nil
}

func ptrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

var ErrNotFound = errors.New("job not found")

type ImportFile struct {
	Jobs []ImportJob `json:"jobs"`
}

type ImportJob struct {
	SourceID       *int64        `json:"source_id"`
	Title          string        `json:"title"`
	Employer       string        `json:"employer"`
	Location       string        `json:"location"`
	EmploymentType string        `json:"employment_type"`
	RoleFamily     string        `json:"role_family"`
	Industry       *string       `json:"industry"`
	Seniority      string        `json:"seniority"`
	IsICT          bool          `json:"is_ict"`
	EducationLevel *string       `json:"education_level"`
	Experience     ImportRange   `json:"experience_years"`
	Salary         ImportSalary  `json:"salary"`
	PostedOn       string        `json:"posted_on"`
	Deadline       *string       `json:"deadline"`
	Summary        string        `json:"summary"`
	Responsibilities []string    `json:"responsibilities"`
	Requirements     []string    `json:"requirements"`
	Skills           []ImportSkill `json:"skills"`
}

type ImportRange struct {
	Min *int `json:"min"`
	Max *int `json:"max"`
}

type ImportSalary struct {
	Min      *float64 `json:"min"`
	Max      *float64 `json:"max"`
	Currency *string  `json:"currency"`
	Period   *string  `json:"period"`
}

type ImportSkill struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Requirement string `json:"requirement"`
	Evidence    string `json:"evidence"`
}

func (s *Store) ImportFromFile(ctx context.Context, path string) (int, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("read import file: %w", err)
	}
	var doc ImportFile
	if err := json.Unmarshal(b, &doc); err != nil {
		return 0, fmt.Errorf("parse import file: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	for i, job := range doc.Jobs {
		if err := importJob(ctx, tx, job); err != nil {
			return 0, fmt.Errorf("job %d (%s): %w", i+1, job.Title, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return len(doc.Jobs), nil
}

func importJob(ctx context.Context, tx *sql.Tx, j ImportJob) error {
	empID, err := ensureEmployer(ctx, tx, j.Employer)
	if err != nil {
		return err
	}

	resp, err := json.Marshal(nonNilStrings(j.Responsibilities))
	if err != nil {
		return err
	}
	reqs, err := json.Marshal(nonNilStrings(j.Requirements))
	if err != nil {
		return err
	}

	isICT := 0
	if j.IsICT {
		isICT = 1
	}

	var source any
	if j.SourceID != nil {
		source = *j.SourceID
	}

	// Upsert by (source_id, title) when source_id is known; otherwise fall back to
	// (title, employer_id, posted_on).
	var jobID int64
	if j.SourceID != nil {
		err = tx.QueryRowContext(ctx, `
			INSERT INTO job (
				source_id, title, employer_id, location, employment_type, role_family, industry,
				seniority, is_ict, education_level, experience_min_years, experience_max_years,
				salary_min, salary_max, salary_currency, salary_period,
				posted_on, deadline, summary, responsibilities, requirements
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(source_id, title) WHERE source_id IS NOT NULL DO UPDATE SET
				employer_id=excluded.employer_id,
				location=excluded.location,
				employment_type=excluded.employment_type,
				role_family=excluded.role_family,
				industry=excluded.industry,
				seniority=excluded.seniority,
				is_ict=excluded.is_ict,
				education_level=excluded.education_level,
				experience_min_years=excluded.experience_min_years,
				experience_max_years=excluded.experience_max_years,
				salary_min=excluded.salary_min,
				salary_max=excluded.salary_max,
				salary_currency=excluded.salary_currency,
				salary_period=excluded.salary_period,
				posted_on=excluded.posted_on,
				deadline=excluded.deadline,
				summary=excluded.summary,
				responsibilities=excluded.responsibilities,
				requirements=excluded.requirements
			RETURNING id`,
			source, j.Title, empID, j.Location, j.EmploymentType, j.RoleFamily, j.Industry,
			j.Seniority, isICT, j.EducationLevel, j.Experience.Min, j.Experience.Max,
			j.Salary.Min, j.Salary.Max, j.Salary.Currency, j.Salary.Period,
			j.PostedOn, j.Deadline, j.Summary, string(resp), string(reqs),
		).Scan(&jobID)
	} else {
		err = tx.QueryRowContext(ctx, `
			INSERT INTO job (
				source_id, title, employer_id, location, employment_type, role_family, industry,
				seniority, is_ict, education_level, experience_min_years, experience_max_years,
				salary_min, salary_max, salary_currency, salary_period,
				posted_on, deadline, summary, responsibilities, requirements
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(title, employer_id, posted_on) DO UPDATE SET
				location=excluded.location,
				employment_type=excluded.employment_type,
				role_family=excluded.role_family,
				industry=excluded.industry,
				seniority=excluded.seniority,
				is_ict=excluded.is_ict,
				education_level=excluded.education_level,
				experience_min_years=excluded.experience_min_years,
				experience_max_years=excluded.experience_max_years,
				salary_min=excluded.salary_min,
				salary_max=excluded.salary_max,
				salary_currency=excluded.salary_currency,
				salary_period=excluded.salary_period,
				deadline=excluded.deadline,
				summary=excluded.summary,
				responsibilities=excluded.responsibilities,
				requirements=excluded.requirements
			RETURNING id`,
			nil, j.Title, empID, j.Location, j.EmploymentType, j.RoleFamily, j.Industry,
			j.Seniority, isICT, j.EducationLevel, j.Experience.Min, j.Experience.Max,
			j.Salary.Min, j.Salary.Max, j.Salary.Currency, j.Salary.Period,
			j.PostedOn, j.Deadline, j.Summary, string(resp), string(reqs),
		).Scan(&jobID)
	}
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM job_skill WHERE job_id = ?`, jobID); err != nil {
		return err
	}
	for _, sk := range j.Skills {
		skillID, err := ensureSkill(ctx, tx, sk.Name, sk.Category)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO job_skill (job_id, skill_id, requirement, evidence)
			VALUES (?, ?, ?, ?)
			ON CONFLICT(job_id, skill_id) DO UPDATE SET
				requirement=excluded.requirement,
				evidence=excluded.evidence`,
			jobID, skillID, sk.Requirement, sk.Evidence,
		); err != nil {
			return err
		}
	}
	return nil
}

func ensureEmployer(ctx context.Context, tx *sql.Tx, name string) (int64, error) {
	if _, err := tx.ExecContext(ctx, `INSERT INTO employer(name) VALUES (?) ON CONFLICT(name) DO NOTHING`, name); err != nil {
		return 0, err
	}
	var id int64
	if err := tx.QueryRowContext(ctx, `SELECT id FROM employer WHERE name = ?`, name).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func ensureSkill(ctx context.Context, tx *sql.Tx, name, category string) (int64, error) {
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO skill(name, category) VALUES (?, ?)
		ON CONFLICT(name) DO UPDATE SET category=excluded.category`, name, category); err != nil {
		return 0, err
	}
	var id int64
	if err := tx.QueryRowContext(ctx, `SELECT id FROM skill WHERE name = ?`, name).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func nonNilStrings(in []string) []string {
	if in == nil {
		return []string{}
	}
	return in
}
