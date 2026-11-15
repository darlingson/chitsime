import { useCallback, useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { listJobs, listSkills, roleFamilies } from './api'
import type { Job, SkillDemand } from './api'
import JobCard from './components/JobCard'
import JobDetail from './components/JobDetail'
import SkillDemandPanel from './components/SkillDemandPanel'
import './App.css'

const PAGE_SIZE = 24

type Filters = {
  roleFamily: string
  skill: string
  postedFrom: string
  postedTo: string
}

type View = 'jobs' | 'skills'

type JobsState = {
  request: string
  jobs: Job[]
  error: string | null
}

const emptyFilters: Filters = {
  roleFamily: '',
  skill: '',
  postedFrom: '',
  postedTo: '',
}

function App() {
  const [view, setView] = useState<View>('jobs')
  const [draft, setDraft] = useState<Filters>(emptyFilters)
  const [query, setQuery] = useState<Filters>(emptyFilters)
  const [offset, setOffset] = useState(0)
  const [reload, setReload] = useState(0)

  const [jobsState, setJobsState] = useState<JobsState>({
    request: '',
    jobs: [],
    error: null,
  })

  const [skills, setSkills] = useState<SkillDemand[]>([])
  const [skillsError, setSkillsError] = useState<string | null>(null)

  const [selected, setSelected] = useState<Job | null>(null)

  const requestKey = [
    query.roleFamily,
    query.skill,
    query.postedFrom,
    query.postedTo,
    String(offset),
    String(reload),
  ].join('|')

  useEffect(() => {
    const controller = new AbortController()

    listJobs(
      {
        roleFamily: query.roleFamily,
        skill: query.skill,
        postedFrom: query.postedFrom,
        postedTo: query.postedTo,
        limit: PAGE_SIZE,
        offset,
      },
      controller.signal,
    )
      .then((response) => {
        if (controller.signal.aborted) return
        setJobsState({ request: requestKey, jobs: response.jobs, error: null })
      })
      .catch((cause: unknown) => {
        if (controller.signal.aborted) return
        setJobsState({
          request: requestKey,
          jobs: [],
          error: cause instanceof Error ? cause.message : 'Could not load jobs',
        })
      })

    return () => controller.abort()
  }, [requestKey, query.roleFamily, query.skill, query.postedFrom, query.postedTo, offset])

  useEffect(() => {
    const controller = new AbortController()

    listSkills(controller.signal)
      .then((response) => {
        if (controller.signal.aborted) return
        setSkills(response.skills)
        setSkillsError(null)
      })
      .catch((cause: unknown) => {
        if (controller.signal.aborted) return
        setSkillsError(cause instanceof Error ? cause.message : 'Could not load skill demand')
      })

    return () => controller.abort()
  }, [])

  const closeDetail = useCallback(() => setSelected(null), [])

  const jobs = jobsState.jobs
  const jobsError = jobsState.error
  const jobsLoading = jobsState.request !== requestKey

  function applyFilters(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setOffset(0)
    setSelected(null)
    setQuery({
      roleFamily: draft.roleFamily,
      skill: draft.skill.trim(),
      postedFrom: draft.postedFrom,
      postedTo: draft.postedTo,
    })
  }

  function resetFilters() {
    setDraft(emptyFilters)
    setQuery(emptyFilters)
    setOffset(0)
    setSelected(null)
  }

  function goToOffset(next: number) {
    setOffset(next)
    setSelected(null)
  }

  const hasFilters =
    query.roleFamily !== '' || query.skill !== '' || query.postedFrom !== '' || query.postedTo !== ''
  const page = Math.floor(offset / PAGE_SIZE) + 1
  const canGoBack = offset > 0
  const canGoForward = jobs.length === PAGE_SIZE

  return (
    <div className="app">
      <header className="masthead">
        <div className="masthead-inner">
          <div className="brand">
            <span className="brand-mark" aria-hidden="true">
              C
            </span>
            <div>
              <h1>Chitsime</h1>
              <p>ICT job and skill demand in Malawi</p>
            </div>
          </div>

          <nav className="tabs" aria-label="Sections">
            <button type="button" aria-pressed={view === 'jobs'} onClick={() => setView('jobs')}>
              Jobs
            </button>
            <button type="button" aria-pressed={view === 'skills'} onClick={() => setView('skills')}>
              Skill demand
            </button>
          </nav>
        </div>
      </header>

      <main className="page">
        {view === 'jobs' ? (
          <>
            <section className="intro">
              <h2>Structured vacancies</h2>
              <p>
                Hand-structured Malawian job adverts with comparable fields and explicit skill
                requirements. Filter by role family, skill, or posting window.
              </p>
            </section>

            <form className="filters" onSubmit={applyFilters}>
              <label className="field">
                <span>Role family</span>
                <select
                  value={draft.roleFamily}
                  onChange={(event) => setDraft({ ...draft, roleFamily: event.target.value })}
                >
                  {roleFamilies.map((family) => (
                    <option key={family.value} value={family.value}>
                      {family.label}
                    </option>
                  ))}
                </select>
              </label>

              <label className="field">
                <span>Skill</span>
                <input
                  list="skill-options"
                  value={draft.skill}
                  placeholder="e.g. DHIS2"
                  onChange={(event) => setDraft({ ...draft, skill: event.target.value })}
                />
                <datalist id="skill-options">
                  {skills.map((skill) => (
                    <option key={skill.name} value={skill.name} />
                  ))}
                </datalist>
              </label>

              <label className="field">
                <span>Posted from</span>
                <input
                  type="date"
                  value={draft.postedFrom}
                  onChange={(event) => setDraft({ ...draft, postedFrom: event.target.value })}
                />
              </label>

              <label className="field">
                <span>Posted to</span>
                <input
                  type="date"
                  value={draft.postedTo}
                  onChange={(event) => setDraft({ ...draft, postedTo: event.target.value })}
                />
              </label>

              <div className="filter-actions">
                <button type="submit" className="primary">
                  Apply filters
                </button>
                <button type="button" className="ghost" onClick={resetFilters}>
                  Reset
                </button>
              </div>
            </form>

            {jobsError ? (
              <div className="banner banner-error" role="alert">
                <span>{jobsError}</span>
                <button type="button" onClick={() => setReload((value) => value + 1)}>
                  Try again
                </button>
              </div>
            ) : jobsLoading && jobs.length === 0 ? (
              <div className="state" aria-live="polite">
                <strong>Loading jobs…</strong>
                Fetching structured adverts from the API.
              </div>
            ) : jobs.length === 0 ? (
              <div className="state">
                <strong>No jobs match those filters</strong>
                {hasFilters
                  ? 'Try widening the date range or clearing the skill.'
                  : 'Run setup to seed the database.'}
              </div>
            ) : (
              <ul className={`job-grid${jobsLoading ? ' is-loading' : ''}`} aria-busy={jobsLoading}>
                {jobs.map((job) => (
                  <JobCard key={job.id} job={job} onSelect={setSelected} />
                ))}
              </ul>
            )}

            <div className="pager">
              <p>
                Page {page} · {jobs.length} {jobs.length === 1 ? 'advert' : 'adverts'}
                {hasFilters ? ' matching your filters' : ''}
              </p>
              <div className="pager-buttons">
                <button
                  type="button"
                  disabled={!canGoBack || jobsLoading}
                  onClick={() => goToOffset(Math.max(0, offset - PAGE_SIZE))}
                >
                  Previous
                </button>
                <button
                  type="button"
                  disabled={!canGoForward || jobsLoading}
                  onClick={() => goToOffset(offset + PAGE_SIZE)}
                >
                  Next
                </button>
              </div>
            </div>
          </>
        ) : (
          <SkillDemandPanel skills={skills} error={skillsError} />
        )}
      </main>

      <footer className="footer">
        Data is manually structured from CareerAd Malawi and Jobsearch Malawi adverts.
      </footer>

      {selected && <JobDetail job={selected} onClose={closeDetail} />}
    </div>
  )
}

export default App
