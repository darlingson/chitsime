import click
from flask import Flask
from flask_login import LoginManager
from flask_migrate import Migrate
from flask_sqlalchemy import SQLAlchemy
from flask_wtf import CSRFProtect

from config import Config

db = SQLAlchemy()
migrate = Migrate()
csrf = CSRFProtect()
login_manager = LoginManager()
login_manager.login_view = "admin.login"
login_manager.login_message = "Please sign in to access the admin area."
login_manager.login_message_category = "danger"


@login_manager.user_loader
def load_user(user_id: str):
    from .models import AdminUser

    return db.session.get(AdminUser, int(user_id))


def create_app():
    app = Flask(__name__)
    app.config.from_object(Config)

    db.init_app(app)
    migrate.init_app(app, db)
    csrf.init_app(app)
    login_manager.init_app(app)

    # Import models so Alembic can detect them
    from . import models  # noqa: F401

    from .routes.jobs import jobs_bp
    from .routes.skills import skills_bp
    from .routes.admin import admin_bp

    app.register_blueprint(jobs_bp)
    app.register_blueprint(skills_bp)
    app.register_blueprint(admin_bp)

    register_cli(app)

    return app


def register_cli(app: Flask) -> None:
    @app.cli.command("create-admin")
    @click.option("--username", prompt=True)
    @click.option("--password", prompt=True, hide_input=True, confirmation_prompt=True)
    def create_admin(username: str, password: str):
        """Create an admin user."""
        from .models import AdminUser

        username = username.strip()
        if not username:
            raise click.ClickException("Username cannot be empty.")
        if len(password) < 8:
            raise click.ClickException("Password must be at least 8 characters.")

        existing = db.session.query(AdminUser).filter_by(username=username).first()
        if existing:
            raise click.ClickException(f'Admin user "{username}" already exists.')

        user = AdminUser(username=username)
        user.set_password(password)
        db.session.add(user)
        db.session.commit()
        click.echo(f'Admin user "{username}" created.')
