package api

import "github.com/darlingson/chitsime/internal/store"

type jobDTO struct {
	ID               int64         `json:"id"`
	Title            string        `json:"title"`
	Employer         string        `json:"employer"`
	Location         string        `json:"location"`
	EmploymentType   string        `json:"employment_type"`
	RoleFamily       string        `json:"role_family"`
	Industry         *string       `json:"industry"`
	Seniority        string        `json:"seniority"`
	IsICT            bool          `json:"is_ict"`
	EducationLevel   *string       `json:"education_level"`
	Experience       rangeDTO      `json:"experience_years"`
	Salary           salaryDTO     `json:"salary"`
	PostedOn         string        `json:"posted_on"`
	Deadline         *string       `json:"deadline"`
	Summary          string        `json:"summary"`
	Responsibilities []string      `json:"responsibilities"`
	Requirements     []string      `json:"requirements"`
	Skills           []jobSkillDTO `json:"skills"`
}

type rangeDTO struct {
	Min *int `json:"min"`
	Max *int `json:"max"`
}

type salaryDTO struct {
	Disclosed bool     `json:"disclosed"`
	Min       *float64 `json:"min"`
	Max       *float64 `json:"max"`
	Currency  *string  `json:"currency"`
	Period    *string  `json:"period"`
}

type jobSkillDTO struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Requirement string `json:"requirement"`
	Evidence    string `json:"evidence"`
}

type skillDemandDTO struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Jobs     int    `json:"jobs"`
	Required int    `json:"required"`
}

func toJobDTO(j store.Job) jobDTO {
	skills := make([]jobSkillDTO, 0, len(j.Skills))
	for _, sk := range j.Skills {
		skills = append(skills, jobSkillDTO{
			Name:        sk.Name,
			Category:    sk.Category,
			Requirement: sk.Requirement,
			Evidence:    sk.Evidence,
		})
	}

	return jobDTO{
		ID:             j.ID,
		Title:          j.Title,
		Employer:       j.Employer,
		Location:       j.Location,
		EmploymentType: j.EmploymentType,
		RoleFamily:     j.RoleFamily,
		Industry:       j.Industry,
		Seniority:      j.Seniority,
		IsICT:          j.IsICT,
		EducationLevel: j.EducationLevel,
		Experience:     rangeDTO{Min: j.ExperienceMin, Max: j.ExperienceMax},
		Salary: salaryDTO{
			Disclosed: j.SalaryMin != nil || j.SalaryMax != nil,
			Min:       j.SalaryMin,
			Max:       j.SalaryMax,
			Currency:  j.SalaryCurrency,
			Period:    j.SalaryPeriod,
		},
		PostedOn:         j.PostedOn,
		Deadline:         j.Deadline,
		Summary:          j.Summary,
		Responsibilities: j.Responsibilities,
		Requirements:     j.Requirements,
		Skills:           skills,
	}
}
