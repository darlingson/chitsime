import type { SkillDemand } from '../api'

type SkillDemandPanelProps = {
  skills: SkillDemand[]
  error: string | null
}

function SkillDemandPanel({ skills, error }: SkillDemandPanelProps) {
  const covered = skills.filter((skill) => skill.jobs > 0).length
  const maxJobs = Math.max(1, ...skills.map((skill) => skill.jobs))

  return (
    <section className="section">
      <div className="intro">
        <h2>Skill demand</h2>
        <p>
          {covered} skills appear across the structured adverts. Bars compare how many adverts mention
          each skill; the required column counts adverts where the skill is essential.
        </p>
      </div>

      {error ? (
        <div className="banner banner-error" role="alert">
          {error}
        </div>
      ) : skills.length === 0 ? (
        <div className="state">
          <strong>No skills recorded</strong>
          Structure an advert and rerun setup to populate this view.
        </div>
      ) : (
        <div className="table-wrap">
          <table className="demand-table">
            <thead>
              <tr>
                <th scope="col">Skill</th>
                <th scope="col">Category</th>
                <th scope="col">Adverts</th>
                <th scope="col">Required</th>
              </tr>
            </thead>
            <tbody>
              {skills.map((skill) => (
                <tr key={skill.name} className={skill.jobs === 0 ? 'is-unused' : undefined}>
                  <td className="demand-name">{skill.name}</td>
                  <td>{skill.category}</td>
                  <td>
                    <div className="demand-count">
                      <div className="bar-track" aria-hidden="true">
                        <div
                          className="bar-fill"
                          style={{ width: `${(skill.jobs / maxJobs) * 100}%` }}
                        />
                      </div>
                      <span>{skill.jobs}</span>
                    </div>
                  </td>
                  <td className="demand-number">{skill.required}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  )
}

export default SkillDemandPanel
