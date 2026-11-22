from datetime import date

from app import db
from app.models import Employer, Job, JobSkill, Skill


class JobImportError(ValueError):
    pass


def parse_date(value):
    if not value:
        return None
    if isinstance(value, date):
        return value
    try:
        return date.fromisoformat(value)
    except Exception as exc:
        raise JobImportError(f"Invalid date {value!r}") from exc


def _get_or_create_employer(name):
    employer = Employer.query.filter_by(name=name).first()
    if employer:
        return employer
    employer = Employer(name=name)
    db.session.add(employer)
    db.session.flush()
    return employer


def _get_or_create_skill(name, category):
    skill = Skill.query.filter_by(name=name).first()
    if skill:
        if category and skill.category != category:
            skill.category = category
        return skill
    skill = Skill(name=name, category=category or "other")
    db.session.add(skill)
    db.session.flush()
    return skill


def upsert_job(data):
    if not isinstance(data, dict):
        raise JobImportError("Each job must be an object")

    title = (data.get("title") or "").strip()
    employer_name = (data.get("employer") or "").strip()
    if not title:
        raise JobImportError("Job is missing title")
    if not employer_name:
        raise JobImportError(f"Job {title!r} is missing employer")

    posted_on = parse_date(data.get("posted_on"))
    if posted_on is None:
        raise JobImportError(f"Job {title!r} is missing posted_on")

    employer = _get_or_create_employer(employer_name)
    job = Job.query.filter_by(
        title=title,
        employer_id=employer.id,
        posted_on=posted_on,
    ).first()
    created = job is None
    if job is None:
        job = Job(title=title, employer_id=employer.id, posted_on=posted_on)
        db.session.add(job)

    experience = data.get("experience_years") or {}
    salary = data.get("salary") or {}

    job.source_id = data.get("source_id")
    job.location = (data.get("location") or "").strip()
    job.employment_type = data.get("employment_type")
    job.role_family = data.get("role_family")
    job.industry = data.get("industry")
    job.seniority = data.get("seniority")
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
        if not isinstance(item, dict):
            raise JobImportError(f"Job {title!r} has an invalid skill entry")
        skill_name = (item.get("name") or "").strip()
        if not skill_name:
            continue
        skill = _get_or_create_skill(skill_name, (item.get("category") or "other").strip())
        db.session.add(
            JobSkill(
                job_id=job.id,
                skill_id=skill.id,
                requirement=item.get("requirement") or "mentioned",
                evidence=item.get("evidence"),
            )
        )

    return job, created


def import_jobs(payload):
    if isinstance(payload, dict):
        jobs = payload.get("jobs")
    elif isinstance(payload, list):
        jobs = payload
    else:
        raise JobImportError("JSON payload must be an object with a 'jobs' list or a list of jobs")

    if not isinstance(jobs, list):
        raise JobImportError("'jobs' must be a list")

    created = 0
    updated = 0
    for index, item in enumerate(jobs, start=1):
        try:
            _, was_created = upsert_job(item)
        except JobImportError:
            raise
        except Exception as exc:
            title = item.get("title") if isinstance(item, dict) else None
            label = f"job {index}" + (f" ({title!r})" if title else "")
            raise JobImportError(f"Failed to import {label}: {exc}") from exc
        if was_created:
            created += 1
        else:
            updated += 1

    db.session.commit()
    return {"total": len(jobs), "created": created, "updated": updated}
