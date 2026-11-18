# chitsime

A Go + SQLite API and dashboard over manually structured Malawian ICT job
adverts.

## Structure

- `internal/database/schema.sql` - tables: `employer`, `skill`, `job`, `job_skill`
- `internal/store/store.go` - SQLite access and idempotent import
- `internal/api/` - HTTP handlers and JSON shapes
- `cmd/setup/main.go` - creates or migrates the database schema
- `cmd/import/main.go` - imports structured jobs from JSON
- `frontend/` - React + TypeScript dashboard (Vite)
- `main.go` - serves the API and the embedded frontend

Each `job` row has comparable fields (`source_id`, `role_family`, `seniority`,
`education_level`, `experience_min_years`, `posted_on`, ...), plus JSON arrays
for responsibilities and requirements. Each `job_skill` row records whether a
skill was `required`, `preferred`, or only `mentioned`, with the evidence text.

## Data

- `data/jobs-raw.json` - raw source (686 adverts). Local working data, not committed.
- `data/jobs-structured.json` - curated ICT dataset imported into the database.

## Run it

Create the database schema:

```powershell
go run ./cmd/setup -fresh
```

Import the structured ICT jobs:

```powershell
go run ./cmd/import -file data/jobs-structured.json
```

The import is idempotent. Re-running it upserts employers, skills, jobs, and
each job's skill links by `(source_id, title)`, so the same file can be imported
again after edits without duplicating rows.

Start the API:

```powershell
go run .
```

It listens on `:8080` by default. `DB_PATH` and `ADDR` can override that.
The same server also serves the dashboard at `http://localhost:8080/`.

## Frontend

The dashboard lives in `frontend/` and is built with Vite. The Go binary embeds
`frontend/dist`, so build it before building or running the server:

```powershell
cd frontend
npm install
npm run build
cd ..
go run .
```

`go build ./...` fails while `frontend/dist` is missing, because of the
`//go:embed frontend/dist` directive in `main.go`.

For live reload, keep the API running on `:8080` and start Vite in another
terminal:

```powershell
cd frontend
npm run dev
```

Vite serves on `:5173` and proxies `/api` to `http://localhost:8080`.

## Endpoints

```text
GET /api/health
GET /api/jobs
GET /api/jobs/{id}
GET /api/skills
```

Examples:

```powershell
curl.exe http://localhost:8080/api/jobs
curl.exe http://localhost:8080/api/jobs/1
curl.exe "http://localhost:8080/api/jobs?role_family=data"
curl.exe "http://localhost:8080/api/jobs?skill=DHIS2"
curl.exe "http://localhost:8080/api/jobs?posted_from=2026-11-01&posted_to=2026-11-30"
curl.exe http://localhost:8080/api/skills
```

## Adding or updating jobs

Edit `data/jobs-structured.json`, then rerun:

```powershell
go run ./cmd/import -file data/jobs-structured.json
```

Use `go run ./cmd/setup -fresh` only when you want to rebuild the database from
scratch. The database file lives in `data/` and is not committed.
