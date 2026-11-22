import json
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
if str(ROOT) not in sys.path:
    sys.path.insert(0, str(ROOT))

from app import create_app  # noqa: E402
from app.services.job_import import import_jobs  # noqa: E402


def main():
    path = Path(sys.argv[1]) if len(sys.argv) > 1 else ROOT / "data" / "jobs-structured.json"
    payload = json.loads(path.read_text(encoding="utf-8"))

    app = create_app()
    with app.app_context():
        result = import_jobs(payload)
        print(
            f"Imported {result['total']} jobs from {path} "
            f"({result['created']} created, {result['updated']} updated)"
        )


if __name__ == "__main__":
    main()
