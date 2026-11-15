# chitsime

A Go + SQLite API and dashboard over Malawian job adverts. The first five
records are ICT/data roles taken from `jobs.txt`.

## Structure

- `internal/database/schema.sql` - tables: `employer`, `skill`, `job`, `job_skill`
- `internal/database/seed.sql` - the five adverts, their skills, and the skill catalogue
- `internal/store/store.go` - SQLite access
- `internal/api/` - HTTP handlers and JSON shapes
- `cmd/scrape/main.go` - captures raw adverts from job sites into JSON
- `cmd/setup/main.go` - creates the database and inserts the seed rows
- `frontend/` - React + TypeScript dashboard (Vite)
- `main.go` - serves the API and the embedded frontend

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

## Capturing adverts

The scraper saves whole adverts as JSON. It never writes to the database;
structuring is a separate manual step.

```powershell
go run ./cmd/scrape -source careerad
go run ./cmd/scrape -source jobsearch
```

Flags:

- `-source` - `careerad` (default) or `jobsearch`
- `-out` - output file, default `data/<source>-jobs-raw.json`
- `-minutes` - stop after this long; rerun to resume
- `-delay` - pause between requests, default `300ms`
- `-limit` - fetch at most this many pages
- `-fresh` - discard the existing output and start over
- `-sitemap` - override the sitemap URL

Both sources read their sitemap, fetch every job page, and record title,
employer, dates, employment type, location, categories, description,
requirements, application text, and links. Reruns load the existing output and
skip URLs already saved, so `-minutes` works as a checkpoint.

CareerAd pages also expose JSON-LD `JobPosting` data. Some adverts are
subscriber-gated; those are saved with `"gated": true` and a public preview
only.

The working file used for structuring is `data/jobs-raw.json`. It holds the
adverts from every source in one list, newest first, where each advert has an
`id` and a `source`. Raw files are local working data and are not committed.

## Adding an advert

Append the advert's structured rows to `internal/database/seed.sql`, then rerun:

```powershell
go run ./cmd/setup
```

Setup is idempotent. `INSERT OR IGNORE` keeps existing rows, so use `-fresh` only
when you want to rebuild the database from scratch.

The database file lives in `data/` and is not committed.
