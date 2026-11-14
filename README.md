# chitsime

A Go + SQLite API over Malawian job adverts. The first five records are
ICT/data roles taken from `jobs.txt`.

## Structure

- `internal/database/schema.sql` - tables: `employer`, `skill`, `job`, `job_skill`
- `internal/database/seed.sql` - the five adverts, their skills, and the skill catalogue
- `internal/store/store.go` - SQLite access
- `internal/api/` - HTTP handlers and JSON shapes
- `cmd/setup/main.go` - creates the database and inserts the seed rows
- `main.go` - starts the API

Each `job` row has comparable fields (`role_family`, `seniority`,
`education_level`, `experience_min_years`, `posted_on`, ...), plus JSON arrays
for responsibilities and requirements. Each `job_skill` row records whether a
skill was `required`, `preferred`, or only `mentioned`, with the evidence text.

## Run it

Create the database and insert the structured rows:

```powershell
go run ./cmd/setup -fresh
```

Expected:

```text
database ready: data/chitsime.db (5 jobs)
```

Start the API:

```powershell
go run .
```

It listens on `:8080` by default. `DB_PATH` and `ADDR` can override that.

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

## Adding an advert

Append the advert's structured rows to `internal/database/seed.sql`, then rerun:

```powershell
go run ./cmd/setup
```

Setup is idempotent. `INSERT OR IGNORE` keeps existing rows, so use `-fresh` only
when you want to rebuild the database from scratch.

The database file lives in `data/` and is not committed.
