from flask import Blueprint, render_template, request, abort
from sqlalchemy import func
from app import db
from app.models import Job, Employer, Skill, JobSkill, ROLE_FAMILIES

jobs_bp = Blueprint("jobs", __name__)

PAGE_SIZE = 24

ROLE_FAMILY_LABELS = {
    "software": "Software",
    "data": "Data",
    "gis": "GIS",
    "infrastructure": "Infrastructure",
    "support": "Support",
    "security": "Security",
    "management": "Management",
    "other": "Other",
}


@jobs_bp.route("/")
def index():
    return _job_list()


@jobs_bp.route("/jobs")
def job_list():
    return _job_list()


def _job_list():
    role_family = request.args.get("role_family", "").strip()
    skill_name = request.args.get("skill", "").strip()
    posted_from = request.args.get("posted_from", "").strip()
    posted_to = request.args.get("posted_to", "").strip()

    try:
        offset = max(0, int(request.args.get("offset", 0)))
    except (ValueError, TypeError):
        offset = 0

    # Validate role_family
    if role_family and role_family not in ROLE_FAMILIES:
        role_family = ""

    query = (
        db.session.query(Job)
        .join(Employer)
        .order_by(Job.posted_on.desc(), Job.id.desc())
    )

    if role_family:
        query = query.filter(Job.role_family == role_family)

    if skill_name:
        query = query.filter(
            Job.job_skills.any(
                JobSkill.skill.has(func.lower(Skill.name) == skill_name.lower())
            )
        )

    if posted_from:
        query = query.filter(Job.posted_on >= posted_from)

    if posted_to:
        query = query.filter(Job.posted_on <= posted_to)

    jobs = query.offset(offset).limit(PAGE_SIZE).all()

    # Skills for datalist autocomplete
    all_skills = db.session.query(Skill.name).order_by(Skill.name).all()
    skill_names = [s.name for s in all_skills]

    has_filters = any([role_family, skill_name, posted_from, posted_to])
    page = (offset // PAGE_SIZE) + 1
    can_go_back = offset > 0
    can_go_forward = len(jobs) == PAGE_SIZE

    role_families = [{"value": "", "label": "All role families"}] + [
        {"value": rf, "label": ROLE_FAMILY_LABELS[rf]} for rf in ROLE_FAMILIES
    ]

    return render_template(
        "jobs/index.html",
        jobs=jobs,
        skill_names=skill_names,
        role_families=role_families,
        filters={
            "role_family": role_family,
            "skill": skill_name,
            "posted_from": posted_from,
            "posted_to": posted_to,
        },
        offset=offset,
        page=page,
        page_size=PAGE_SIZE,
        has_filters=has_filters,
        can_go_back=can_go_back,
        can_go_forward=can_go_forward,
    )


@jobs_bp.route("/jobs/<int:job_id>")
def job_detail(job_id):
    job = db.session.get(Job, job_id)
    if job is None:
        abort(404)

    # Skills ordered: required → preferred → mentioned, then by category and name
    skills = (
        db.session.query(JobSkill)
        .filter(JobSkill.job_id == job_id)
        .join(Skill)
        .order_by(
            db.case(
                (JobSkill.requirement == "required", 0),
                (JobSkill.requirement == "preferred", 1),
                else_=2,
            ),
            Skill.category,
            Skill.name,
        )
        .all()
    )

    return render_template("jobs/detail.html", job=job, skills=skills)
