# Generic SQLite Browser Example

This is a minimal `db-browser serve` app that lists user tables from a SQLite database using `express`, `db`, and `ui.dsl`.

```bash
DB=/tmp/db-browser-example.sqlite
python3 - <<'PY'
import os, sqlite3
path = os.environ.get('DB', '/tmp/db-browser-example.sqlite')
con = sqlite3.connect(path)
con.execute('create table if not exists people(id integer primary key, name text)')
con.commit()
con.close()
PY

go run ./cmd/db-browser serve \
  --db "$DB" \
  --scripts-dir examples/generic-browser/scripts \
  --addr :8080 \
  --dev
```

Open <http://127.0.0.1:8080/>.
