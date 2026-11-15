# Chitsime frontend

React + TypeScript dashboard for the structured job data. Vite builds it into
`dist/`, which `main.go` embeds and serves from the Go binary.

## Commands

```powershell
npm install
npm run dev     # dev server on :5173, proxies /api to :8080
npm run build   # type-check and emit dist/
npm run lint    # oxlint
```

## Layout

- `src/App.tsx` - filters, pagination, and switching between the jobs and skill views
- `src/api.ts` - typed client for `/api/jobs` and `/api/skills`
- `src/format.ts` - date, salary, experience, and label formatting
- `src/components/JobCard.tsx` - summary card in the jobs grid
- `src/components/JobDetail.tsx` - modal with the full advert and skill table
- `src/components/SkillDemandPanel.tsx` - skill demand table

Set `VITE_API_BASE` to point at a different API host; it defaults to the same
origin, which is what the Go server provides.
