import json
import sys
from datetime import date
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
if str(ROOT) not in sys.path:
    sys.path.insert(0, str(ROOT))

from app import create_app, db  # noqa: E402
from app.models import Employer, Job, JobSkill, Skill  # noqa: E402


def parse_date(value):
    if not value:
        return None
    return date.fromisoformat(value)


def get_or_create_employer(name):
    employer = Employer.query.filter_by(name=name).first()
    if employer:
        return employer
    employer = Employer(name=name)
    db.session.add(employer)
    db.session.flush()
    return employer


def get_or_create_skill(name, category):
    skill = Skill.query.filter_by(name=name).first()
    if skill:
        if category and skill.category != category:
            skill.category = category
        return skill
    skill = Skill(name=name, category=category)
    db.session.add(skill)
    db.session.flush()
    return skill


def upsert_job(data):
    employer = get_or_create_employer(data["employer"])
    posted_on = parse_date(data["posted_on"])
    if posted_on is None:
        raise ValueError(f"Job {data.get('title')!r} is missing posted_on")

    job = Job.query.filter_by(
        title=data["title"],
        employer_id=employer.id,
        posted_on=posted_on,
    ).first()
    if job is None:
        job = Job(
            title=data["title"],
            employer_id=employer.id,
            posted_on=posted_on,
        )
        db.session.add(job)

    experience = data.get("experience_years") or {}
    salary = data.get("salary") or {}

    job.source_id = data.get("source_id")
    job.location = data.get("location") or ""
    job.employment_type = data["employment_type"]
    job.role_family = data["role_family"]
    job.industry = data.get("industry")
    job.seniority = data["seniority"]
    job.is_ict = bool(data.get("is_ict", True))
    job.education_level = data.get("education_level")
    job.experience_min_years = experience.get("min")
    job.experience_max_years = experience.get("max")
    job.salary_min = salary.get("min")
    job.salary_max = salary.get("max")
    job.salary_currency = salary.get("currency")
    job.salary_period = salary.get("period")
    job.deadline = parse_date(data.get("deadline"))
    job.summary = data.get("summary") or ""
    job.responsibilities = data.get("responsibilities") or []
    job.requirements = data.get("requirements") or []

    db.session.flush()

    JobSkill.query.filter_by(job_id=job.id).delete()
    for item in data.get("skills") or []:
        skill = get_or_create_skill(item["name"], item.get("category") or "other")
        db.session.add(
            JobSkill(
                job_id=job.id,
                skill_id=skill.id,
                requirement=item.get("requirement") or "mentioned",
                evidence=item.get("evidence"),
            )
        )

    return job


def main():
    path = Path(sys.argv[1]) if len(sys.argv) > 1 else ROOT / "data" / "jobs-structured.json"
    payload = json.loads(path.read_text(encoding="utf-8"))
    jobs = payload["jobs"] if isinstance(payload, dict) else payload

    app = create_app()
    with app.app_context():
        for item in jobs:
            upsert_job(item)
        db.session.commit()
        print(f"Imported {len(jobs)} jobs from {path}")


if __name__ == "__main__":
    main()
