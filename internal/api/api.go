package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/darlingson/chitsime/internal/store"
)

type Server struct {
	store *store.Store
	log   *slog.Logger
}

func New(st *store.Store, log *slog.Logger) http.Handler {
	s := &Server{store: st, log: log}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", s.health)
	mux.HandleFunc("GET /api/jobs", s.listJobs)
	mux.HandleFunc("GET /api/jobs/{id}", s.getJob)
	mux.HandleFunc("GET /api/skills", s.listSkills)

	return s.middleware(mux)
}

func (s *Server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		start := time.Now()
		defer func() {
			if rec := recover(); rec != nil {
				s.log.Error("panic", "err", rec)
				writeError(w, http.StatusInternalServerError, "internal error")
			}
			s.log.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"query", r.URL.RawQuery,
				"duration", time.Since(start).String(),
			)
		}()

		next.ServeHTTP(w, r)
	})
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	if err := s.store.Ping(r.Context()); err != nil {
		writeError(w, http.StatusServiceUnavailable, "database unavailable")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "database": "ok"})
}

var roleFamilies = map[string]bool{
	"software":       true,
	"data":           true,
	"gis":            true,
	"infrastructure": true,
	"support":        true,
	"security":       true,
	"management":     true,
	"other":          true,
}

func (s *Server) listJobs(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	role := q.Get("role_family")
	if role != "" && !roleFamilies[role] {
		writeError(w, http.StatusBadRequest, "role_family must be one of software, data, gis, infrastructure, support, security, management, other")
		return
	}
	postedFrom := q.Get("posted_from")
	if postedFrom != "" && !validDate(postedFrom) {
		writeError(w, http.StatusBadRequest, "posted_from must be YYYY-MM-DD")
		return
	}
	postedTo := q.Get("posted_to")
	if postedTo != "" && !validDate(postedTo) {
		writeError(w, http.StatusBadRequest, "posted_to must be YYYY-MM-DD")
		return
	}

	filter := store.JobFilter{
		RoleFamily: role,
		Skill:      q.Get("skill"),
		PostedFrom: postedFrom,
		PostedTo:   postedTo,
		Limit:      intQuery(q.Get("limit"), 50, 1, 100),
		Offset:     intQuery(q.Get("offset"), 0, 0, 100000),
	}

	jobs, err := s.store.ListJobs(r.Context(), filter)
	if err != nil {
		s.log.Error("list jobs", "err", err)
		writeError(w, http.StatusInternalServerError, "could not list jobs")
		return
	}

	out := make([]jobDTO, 0, len(jobs))
	for _, job := range jobs {
		out = append(out, toJobDTO(job))
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"jobs":   out,
		"count":  len(out),
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

func (s *Server) getJob(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "job id must be a positive integer")
		return
	}

	job, err := s.store.GetJob(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	if err != nil {
		s.log.Error("get job", "err", err)
		writeError(w, http.StatusInternalServerError, "could not load job")
		return
	}
	writeJSON(w, http.StatusOK, toJobDTO(job))
}

func (s *Server) listSkills(w http.ResponseWriter, r *http.Request) {
	demand, err := s.store.SkillDemand(r.Context())
	if err != nil {
		s.log.Error("skill demand", "err", err)
		writeError(w, http.StatusInternalServerError, "could not list skills")
		return
	}

	out := make([]skillDemandDTO, 0, len(demand))
	for _, d := range demand {
		out = append(out, skillDemandDTO{
			Name:     d.Name,
			Category: d.Category,
			Jobs:     d.Jobs,
			Required: d.Required,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"skills": out,
		"count":  len(out),
	})
}

func intQuery(raw string, def, min, max int) int {
	if raw == "" {
		return def
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return def
	}
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}

func validDate(s string) bool {
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
