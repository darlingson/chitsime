import { useEffect, useRef } from 'react'
import type { Job } from '../api'
import { formatDate, formatExperience, formatSalary, titleCase } from '../format'

type JobDetailProps = {
  job: Job
  onClose: () => void
}

type MetaItemProps = {
  label: string
  value: string
}

function MetaItem({ label, value }: MetaItemProps) {
  return (
    <div>
      <dt>{label}</dt>
      <dd>{value}</dd>
    </div>
  )
}

function JobDetail({ job, onClose }: JobDetailProps) {
  const closeButton = useRef<HTMLButtonElement>(null)

  useEffect(() => {
    closeButton.current?.focus()

    function onKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape') {
        onClose()
      }
    }

    const previousOverflow = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    document.addEventListener('keydown', onKeyDown)

    return () => {
      document.body.style.overflow = previousOverflow
      document.removeEventListener('keydown', onKeyDown)
    }
  }, [onClose])

  return (
    <div className="modal-backdrop" onClick={onClose} role="presentation">
      <div
        className="modal"
        role="dialog"
        aria-modal="true"
        aria-labelledby="job-detail-title"
        onClick={(event) => event.stopPropagation()}
      >
        <header className="modal-head">
          <div>
            <h2 id="job-detail-title">{job.title}</h2>
            <p className="employer">
              {job.employer}
              {job.industry ? ` · ${job.industry}` : ''}
            </p>
          </div>
          <button
            ref={closeButton}
            type="button"
            className="modal-close"
            onClick={onClose}
            aria-label="Close job details"
          >
            ×
          </button>
        </header>

        <div className="modal-body">
          <dl className="detail-meta">
            <MetaItem label="Location" value={job.location || '—'} />
            <MetaItem label="Employment" value={titleCase(job.employment_type)} />
            <MetaItem label="Role family" value={titleCase(job.role_family)} />
            <MetaItem label="Seniority" value={titleCase(job.seniority)} />
            <MetaItem
              label="Education"
              value={job.education_level ? titleCase(job.education_level) : 'Not specified'}
            />
            <MetaItem label="Experience" value={formatExperience(job.experience_years)} />
            <MetaItem label="Salary" value={formatSalary(job.salary)} />
            <MetaItem label="Posted" value={formatDate(job.posted_on)} />
            <MetaItem label="Deadline" value={formatDate(job.deadline)} />
          </dl>

          {job.summary && (
            <section className="detail-section">
              <h3>Summary</h3>
              <p>{job.summary}</p>
            </section>
          )}

          {job.responsibilities.length > 0 && (
            <section className="detail-section">
              <h3>Responsibilities</h3>
              <ul className="bullet-list">
                {job.responsibilities.map((item) => (
                  <li key={item}>{item}</li>
                ))}
              </ul>
            </section>
          )}

          {job.requirements.length > 0 && (
            <section className="detail-section">
              <h3>Requirements</h3>
              <ul className="bullet-list">
                {job.requirements.map((item) => (
                  <li key={item}>{item}</li>
                ))}
              </ul>
            </section>
          )}

          {job.skills.length > 0 && (
            <section className="detail-section">
              <h3>Skills ({job.skills.length})</h3>
              <table className="skill-table">
                <thead>
                  <tr>
                    <th scope="col">Skill</th>
                    <th scope="col">Category</th>
                    <th scope="col">Requirement</th>
                    <th scope="col">Evidence</th>
                  </tr>
                </thead>
                <tbody>
                  {job.skills.map((skill) => (
                    <tr key={skill.name}>
                      <td>{skill.name}</td>
                      <td>{skill.category}</td>
                      <td>
                        <span className={`req req-${skill.requirement}`}>{skill.requirement}</span>
                      </td>
                      <td className="evidence">{skill.evidence || '—'}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </section>
          )}

          {job.source_url && (
            <a className="source-link" href={job.source_url} target="_blank" rel="noreferrer">
              View original advert
            </a>
          )}
        </div>
      </div>
    </div>
  )
}

export default JobDetail
