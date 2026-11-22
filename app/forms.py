from flask_wtf import FlaskForm
from flask_wtf.file import FileAllowed, FileField, FileRequired
from wtforms import (
    StringField, TextAreaField, SelectField, BooleanField,
    IntegerField, FloatField, FieldList, FormField, DateField, PasswordField,
)
from wtforms.validators import DataRequired, Optional, Length, NumberRange
from app.models import (
    EMPLOYMENT_TYPES, ROLE_FAMILIES, SENIORITY_LEVELS,
    EDUCATION_LEVELS, SALARY_PERIODS, REQUIREMENT_LEVELS,
)


def _choices(values, include_blank=False, blank_label="— select —"):
    choices = [(v, v.replace("_", " ").title()) for v in values]
    if include_blank:
        choices.insert(0, ("", blank_label))
    return choices


class SkillEntryForm(FlaskForm):
    """Sub-form for a single skill row on the job form."""
    class Meta:
        csrf = False  # CSRF handled by parent form

    name = StringField("Skill name", validators=[Optional(), Length(max=255)])
    category = StringField("Category", validators=[Optional(), Length(max=100)])
    requirement = SelectField(
        "Requirement",
        choices=_choices(REQUIREMENT_LEVELS),
        default="mentioned",
    )
    evidence = TextAreaField("Evidence", validators=[Optional()])


class JobForm(FlaskForm):
    title = StringField("Job title", validators=[DataRequired(), Length(max=255)])
    employer = StringField("Employer", validators=[DataRequired(), Length(max=255)])
    location = StringField("Location", validators=[DataRequired(), Length(max=255)])

    employment_type = SelectField(
        "Employment type",
        choices=_choices(EMPLOYMENT_TYPES),
        validators=[DataRequired()],
    )
    role_family = SelectField(
        "Role family",
        choices=_choices(ROLE_FAMILIES),
        validators=[DataRequired()],
    )
    industry = StringField("Industry", validators=[Optional(), Length(max=100)])
    seniority = SelectField(
        "Seniority",
        choices=_choices(SENIORITY_LEVELS),
        validators=[DataRequired()],
    )
    is_ict = BooleanField("ICT role", default=True)
    education_level = SelectField(
        "Minimum education",
        choices=_choices(EDUCATION_LEVELS, include_blank=True, blank_label="Not specified"),
        validators=[Optional()],
    )

    experience_min_years = IntegerField(
        "Experience min (years)", validators=[Optional(), NumberRange(min=0, max=50)]
    )
    experience_max_years = IntegerField(
        "Experience max (years)", validators=[Optional(), NumberRange(min=0, max=50)]
    )

    salary_min = FloatField("Salary min", validators=[Optional()])
    salary_max = FloatField("Salary max", validators=[Optional()])
    salary_currency = StringField("Currency", validators=[Optional(), Length(max=10)])
    salary_period = SelectField(
        "Salary period",
        choices=_choices(SALARY_PERIODS, include_blank=True, blank_label="Not specified"),
        validators=[Optional()],
    )

    posted_on = DateField("Posted on", validators=[DataRequired()])
    deadline = DateField("Application deadline", validators=[Optional()])

    summary = TextAreaField("Summary", validators=[DataRequired()])
    responsibilities = TextAreaField(
        "Responsibilities (one per line)", validators=[Optional()]
    )
    requirements = TextAreaField(
        "Requirements (one per line)", validators=[Optional()]
    )

    skills = FieldList(FormField(SkillEntryForm), min_entries=1)


class LoginForm(FlaskForm):
    username = StringField("Username", validators=[DataRequired(), Length(max=80)])
    password = PasswordField("Password", validators=[DataRequired(), Length(max=255)])


class ImportJobsForm(FlaskForm):
    file = FileField(
        "JSON file",
        validators=[FileRequired(), FileAllowed(["json"], "JSON files only.")],
    )
