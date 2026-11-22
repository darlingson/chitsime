import json

from flask import Blueprint, render_template, request, redirect, url_for, flash, abort
from flask_login import current_user, login_required, login_user, logout_user
from sqlalchemy import func as sqla_func

from app import db
from app.forms import ImportJobsForm, JobForm, LoginForm
from app.models import (
    AdminUser,
    Job, Employer, Skill, JobSkill,
    EMPLOYMENT_TYPES, ROLE_FAMILIES, SENIORITY_LEVELS,
    EDUCATION_LEVELS, SALARY_PERIODS, REQUIREMENT_LEVELS,
)
from app.services.job_import import JobImportError, import_jobs

admin_bp = Blueprint("admin", __name__, url_prefix="/admin")


def _get_or_create_employer(name: str) -> Employer:
    name = name.strip()
    employer = db.session.query(Employer).filter_by(name=name).first()
    if not employer:
        employer = Employer(name=name)
        db.session.add(employer)
        db.session.flush()
    return employer


def _get_or_create_skill(name: str, category: str) -> Skill:
    name = name.strip()
    category = category.strip()
    skill = db.session.query(Skill).filter(
        func_lower(Skill.name) == name.lower()
    ).first()
    if not skill:
        skill = Skill(name=name, category=category)
        db.session.add(skill)
        db.session.flush()
    else:
        if skill.category != category and category:
            skill.category = category
    return skill


def func_lower(col):
    return sqla_func.lower(col)


@admin_bp.route("/")
def index():
    if current_user.is_authenticated:
        return redirect(url_for("admin.job_list"))
    return redirect(url_for("admin.login"))


@admin_bp.route("/login", methods=["GET", "POST"])
def login():
    if current_user.is_authenticated:
        return redirect(url_for("admin.job_list"))

    form = LoginForm()
    if form.validate_on_submit():
        user = db.session.query(AdminUser).filter_by(username=form.username.data.strip()).first()
        if user and user.check_password(form.password.data):
            login_user(user)
            next_url = request.args.get("next", "")
            if not next_url.startswith("/"):
                next_url = url_for("admin.job_list")
            return redirect(next_url)
        flash("Invalid username or password.", "danger")

    return render_template("admin/login.html", form=form)


@admin_bp.route("/logout", methods=["POST"])
@login_required
def logout():
    logout_user()
    flash("Signed out.", "success")
    return redirect(url_for("admin.login"))


@admin_bp.route("/jobs")
@login_required
def job_list():
    page = request.args.get("page", 1, type=int)
    pagination = (
        db.session.query(Job)
        .join(Employer)
        .order_by(Job.posted_on.desc(), Job.id.desc())
        .paginate(page=page, per_page=30, error_out=False)
    )
    return render_template("admin/job_list.html", pagination=pagination)


@admin_bp.route("/jobs/new", methods=["GET", "POST"])
@login_required
def new_job():
    form = JobForm()
    if form.validate_on_submit():
        try:
            employer = _get_or_create_employer(form.employer.data)
            job = Job(
                title=form.title.data,
                employer=employer,
                location=form.location.data,
                employment_type=form.employment_type.data,
                role_family=form.role_family.data,
                industry=form.industry.data or None,
                seniority=form.seniority.data,
                is_ict=form.is_ict.data,
                education_level=form.education_level.data or None,
                experience_min_years=form.experience_min_years.data or None,
                experience_max_years=form.experience_max_years.data or None,
                salary_min=form.salary_min.data or None,
                salary_max=form.salary_max.data or None,
                salary_currency=form.salary_currency.data or None,
                salary_period=form.salary_period.data or None,
                posted_on=form.posted_on.data,
                deadline=form.deadline.data or None,
                summary=form.summary.data,
                responsibilities=_parse_lines(form.responsibilities.data),
                requirements=_parse_lines(form.requirements.data),
            )
            db.session.add(job)
            db.session.flush()

            _save_skills(job, form)

            db.session.commit()
            flash(f'Job "{job.title}" added successfully.', "success")
            return redirect(url_for("admin.job_list"))
        except Exception as exc:
            db.session.rollback()
            flash(f"Error saving job: {exc}", "danger")

    return render_template("admin/job_form.html", form=form, job=None)


@admin_bp.route("/jobs/<int:job_id>/edit", methods=["GET", "POST"])
@login_required
def edit_job(job_id):
    job = db.session.get(Job, job_id)
    if job is None:
        abort(404)

    form = JobForm(obj=job)

    if request.method == "GET":
        # Pre-fill employer text field and text areas
        form.employer.data = job.employer.name
        form.responsibilities.data = "\n".join(job.responsibilities or [])
        form.requirements.data = "\n".join(job.requirements or [])

    if form.validate_on_submit():
        try:
            employer = _get_or_create_employer(form.employer.data)
            job.title = form.title.data
            job.employer = employer
            job.location = form.location.data
            job.employment_type = form.employment_type.data
            job.role_family = form.role_family.data
            job.industry = form.industry.data or None
            job.seniority = form.seniority.data
            job.is_ict = form.is_ict.data
            job.education_level = form.education_level.data or None
            job.experience_min_years = form.experience_min_years.data or None
            job.experience_max_years = form.experience_max_years.data or None
            job.salary_min = form.salary_min.data or None
            job.salary_max = form.salary_max.data or None
            job.salary_currency = form.salary_currency.data or None
            job.salary_period = form.salary_period.data or None
            job.posted_on = form.posted_on.data
            job.deadline = form.deadline.data or None
            job.summary = form.summary.data
            job.responsibilities = _parse_lines(form.responsibilities.data)
            job.requirements = _parse_lines(form.requirements.data)

            # Replace skills
            db.session.query(JobSkill).filter_by(job_id=job.id).delete()
            db.session.flush()
            _save_skills(job, form)

            db.session.commit()
            flash(f'Job "{job.title}" updated.', "success")
            return redirect(url_for("admin.job_list"))
        except Exception as exc:
            db.session.rollback()
            flash(f"Error updating job: {exc}", "danger")

    return render_template("admin/job_form.html", form=form, job=job)


@admin_bp.route("/jobs/<int:job_id>/delete", methods=["POST"])
@login_required
def delete_job(job_id):
    job = db.session.get(Job, job_id)
    if job is None:
        abort(404)
    title = job.title
    db.session.delete(job)
    db.session.commit()
    flash(f'Job "{title}" deleted.', "success")
    return redirect(url_for("admin.job_list"))


@admin_bp.route("/import", methods=["GET", "POST"])
@login_required
def import_jobs_view():
    form = ImportJobsForm()
    if form.validate_on_submit():
        try:
            payload = json.load(form.file.data.stream)
            result = import_jobs(payload)
            flash(
                f"Imported {result['total']} jobs ({result['created']} created, {result['updated']} updated).",
                "success",
            )
            return redirect(url_for("admin.job_list"))
        except json.JSONDecodeError as exc:
            db.session.rollback()
            flash(f"Invalid JSON file: {exc.msg}.", "danger")
        except JobImportError as exc:
            db.session.rollback()
            flash(str(exc), "danger")
        except Exception as exc:
            db.session.rollback()
            flash(f"Import failed: {exc}", "danger")

    return render_template("admin/import_jobs.html", form=form)


# ── helpers ──────────────────────────────────────────────────────────────────

def _parse_lines(text: str) -> list[str]:
    """Split a textarea value on newlines, strip blanks."""
    if not text:
        return []
    return [line.strip() for line in text.splitlines() if line.strip()]


def _save_skills(job: Job, form: JobForm):
    """Read skill rows from the dynamic skill subform and upsert them."""
    for skill_entry in form.skills.entries:
        skill_name = (skill_entry.data.get("name") or "").strip()
        if not skill_name:
            continue
        skill_category = (skill_entry.data.get("category") or "").strip() or "general"
        requirement = skill_entry.data.get("requirement") or "mentioned"
        evidence = (skill_entry.data.get("evidence") or "").strip() or None

        skill = _get_or_create_skill(skill_name, skill_category)
        js = JobSkill(
            job_id=job.id,
            skill_id=skill.id,
            requirement=requirement,
            evidence=evidence,
        )
        db.session.add(js)
