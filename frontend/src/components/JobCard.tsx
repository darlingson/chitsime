import type { Job } from '../api'
import { formatDate, titleCase } from '../format'

type JobCardProps = {
  job: Job
  onSelect: (job: Job) => void
}

function JobCard({ job, onSelect }: JobCardProps) {
  const visibleSkills = job.skills.slice(0, 4)
  const hiddenSkills = job.skills.length - visibleSkills.length

  return (
    <li className="job-card">
      <div className="job-card-top">
        <h3>{job.title}</h3>
        {job.is_ict && <span className="badge-ict">ICT</span>}
      </div>

      <p className="employer">{job.employer}</p>

      <dl className="job-meta">
        <div>
          <dt>Location</dt>
          <dd>{job.location || '—'}</dd>
        </div>
        <div>
          <dt>Type</dt>
          <dd>{titleCase(job.employment_type)}</dd>
        </div>
        <div>
          <dt>Posted</dt>
          <dd>{formatDate(job.posted_on)}</dd>
        </div>
      </dl>

      {visibleSkills.length > 0 && (
        <ul className="chips" aria-label="Skills">
          {visibleSkills.map((skill) => (
            <li key={skill.name} className={`chip chip-${skill.requirement}`}>
              {skill.name}
            </li>
          ))}
          {hiddenSkills > 0 && <li className="chip chip-more">+{hiddenSkills} more</li>}
        </ul>
      )}

      <button type="button" className="card-action" onClick={() => onSelect(job)}>
        View details
      </button>
    </li>
  )
}

export default JobCard
