from . import db


class Employer(db.Model):
    __tablename__ = "employer"

    id = db.Column(db.Integer, primary_key=True, autoincrement=True)
    name = db.Column(db.String(255), nullable=False, unique=True)

    jobs = db.relationship("Job", back_populates="employer", lazy="dynamic")

    def __repr__(self):
        return f"<Employer {self.name}>"


class Skill(db.Model):
    __tablename__ = "skill"

    id = db.Column(db.Integer, primary_key=True, autoincrement=True)
    name = db.Column(db.String(255), nullable=False, unique=True)
    category = db.Column(db.String(100), nullable=False)

    job_skills = db.relationship("JobSkill", back_populates="skill", lazy="dynamic")

    def __repr__(self):
        return f"<Skill {self.name}>"


EMPLOYMENT_TYPES = [
    "full_time",
    "part_time",
    "contract",
    "internship",
    "consultancy",
    "temporary",
    "volunteer",
]

ROLE_FAMILIES = [
    "software",
    "data",
    "gis",
    "infrastructure",
    "support",
    "security",
    "management",
    "other",
]

SENIORITY_LEVELS = [
    "intern",
    "junior",
    "mid",
    "senior",
    "manager",
    "director",
    "executive",
]

EDUCATION_LEVELS = [
    "secondary",
    "certificate",
    "diploma",
    "bachelor",
    "master",
    "phd",
]

SALARY_PERIODS = ["hour", "day", "week", "month", "year"]

REQUIREMENT_LEVELS = ["required", "preferred", "mentioned"]


class Job(db.Model):
    __tablename__ = "job"

    id = db.Column(db.Integer, primary_key=True, autoincrement=True)
    source_id = db.Column(db.Integer, nullable=True)

    title = db.Column(db.String(255), nullable=False)
    employer_id = db.Column(db.Integer, db.ForeignKey("employer.id"), nullable=False)
    location = db.Column(db.String(255), nullable=False)

    employment_type = db.Column(db.Enum(*EMPLOYMENT_TYPES), nullable=False)
    role_family = db.Column(db.Enum(*ROLE_FAMILIES), nullable=False)
    industry = db.Column(db.String(100), nullable=True)
    seniority = db.Column(db.Enum(*SENIORITY_LEVELS), nullable=False)
    is_ict = db.Column(db.Boolean, nullable=False, default=True)
    education_level = db.Column(db.Enum(*EDUCATION_LEVELS), nullable=True)

    experience_min_years = db.Column(db.Integer, nullable=True)
    experience_max_years = db.Column(db.Integer, nullable=True)

    salary_min = db.Column(db.Float, nullable=True)
    salary_max = db.Column(db.Float, nullable=True)
    salary_currency = db.Column(db.String(10), nullable=True)
    salary_period = db.Column(db.Enum(*SALARY_PERIODS), nullable=True)

    posted_on = db.Column(db.Date, nullable=False)
    deadline = db.Column(db.Date, nullable=True)

    summary = db.Column(db.Text, nullable=False)
    responsibilities = db.Column(db.JSON, nullable=False, default=list)
    requirements = db.Column(db.JSON, nullable=False, default=list)

    created_at = db.Column(db.DateTime, nullable=False, server_default=db.func.now())

    employer = db.relationship("Employer", back_populates="jobs")
    job_skills = db.relationship(
        "JobSkill", back_populates="job", cascade="all, delete-orphan", lazy="select"
    )

    __table_args__ = (
        db.UniqueConstraint("title", "employer_id", "posted_on", name="uq_job_natural"),
    )

    @property
    def salary_disclosed(self):
        return self.salary_min is not None or self.salary_max is not None

    def __repr__(self):
        return f"<Job {self.id}: {self.title}>"


class JobSkill(db.Model):
    __tablename__ = "job_skill"

    job_id = db.Column(db.Integer, db.ForeignKey("job.id", ondelete="CASCADE"), primary_key=True)
    skill_id = db.Column(db.Integer, db.ForeignKey("skill.id"), primary_key=True)
    requirement = db.Column(db.Enum(*REQUIREMENT_LEVELS), nullable=False)
    evidence = db.Column(db.Text, nullable=True)

    job = db.relationship("Job", back_populates="job_skills")
    skill = db.relationship("Skill", back_populates="job_skills")

    def __repr__(self):
        return f"<JobSkill job={self.job_id} skill={self.skill_id} req={self.requirement}>"
