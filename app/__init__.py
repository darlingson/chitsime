from flask import Flask
from flask_sqlalchemy import SQLAlchemy
from flask_migrate import Migrate
from config import Config

db = SQLAlchemy()
migrate = Migrate()


def create_app():
    app = Flask(__name__)
    app.config.from_object(Config)

    db.init_app(app)
    migrate.init_app(app, db)

    # Import models so Alembic can detect them
    from . import models  # noqa: F401

    from .routes.jobs import jobs_bp
    from .routes.skills import skills_bp
    from .routes.admin import admin_bp

    app.register_blueprint(jobs_bp)
    app.register_blueprint(skills_bp)
    app.register_blueprint(admin_bp)

    return app
