# YAML Dashboard Example

This example loads `dashboard.yaml` through `require("yaml")` and renders configured SQL metrics with `ui.table.fromRows`.

```bash
DB=/tmp/db-browser-yaml-dashboard.sqlite
python3 - <<'PY'
import os, sqlite3
path = os.environ.get('DB', '/tmp/db-browser-yaml-dashboard.sqlite')
con = sqlite3.connect(path)
con.execute('create table if not exists people(id integer primary key, name text)')
con.executemany('insert into people(name) values (?)', [('Alice',), ('Bob',)])
con.commit()
con.close()
PY

go run ./cmd/db-browser serve \
  --db "$DB" \
  --scripts-dir examples/yaml-dashboard/scripts \
  --addr :8081 \
  --dev
```

Open <http://127.0.0.1:8081/>.
