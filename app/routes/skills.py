from flask import Blueprint, render_template
from sqlalchemy import func, case
from app import db
from app.models import Skill, JobSkill

skills_bp = Blueprint("skills", __name__)


@skills_bp.route("/skills")
def skill_demand():
    rows = (
        db.session.query(
            Skill.name,
            Skill.category,
            func.count(JobSkill.job_id).label("jobs"),
            func.sum(
                case((JobSkill.requirement == "required", 1), else_=0)
            ).label("required"),
        )
        .outerjoin(JobSkill, JobSkill.skill_id == Skill.id)
        .group_by(Skill.id, Skill.name, Skill.category)
        .order_by(func.count(JobSkill.job_id).desc(), Skill.name)
        .all()
    )

    max_jobs = rows[0].jobs if rows else 1
    covered = sum(1 for r in rows if r.jobs > 0)

    return render_template(
        "skills/index.html",
        skills=rows,
        max_jobs=max_jobs,
        covered=covered,
    )
