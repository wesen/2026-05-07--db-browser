package verbcli

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/dop251/goja"
	noderequire "github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/go-go-goja/engine"
	"github.com/go-go-golems/go-go-goja/pkg/jsverbs"
	_ "github.com/mattn/go-sqlite3"
)

type RuntimeSettings struct {
	DBPath      string
	ReadOnly    bool
	AllowWrites bool
}

func runtimeInvokerFactory(settings *RuntimeSettings) InvokerFactory {
	return func(repo ScannedRepository, _ *jsverbs.VerbSpec) jsverbs.VerbInvoker {
		return func(ctx context.Context, registry *jsverbs.Registry, verb *jsverbs.VerbSpec, parsedValues *values.Values) (interface{}, error) {
			factory, cleanup, err := newRuntimeFactory(repo, settings)
			if err != nil {
				return nil, err
			}
			defer cleanup()

			rt, err := factory.NewRuntime(ctx)
			if err != nil {
				return nil, err
			}
			defer func() { _ = rt.Close(context.Background()) }()

			return registry.InvokeInRuntime(ctx, rt, verb, parsedValues)
		}
	}
}

func newRuntimeFactory(repo ScannedRepository, settings *RuntimeSettings) (*engine.Factory, func(), error) {
	if repo.Registry == nil {
		return nil, nil, fmt.Errorf("repository %s has no jsverbs registry", describeRepository(repo))
	}
	if settings == nil {
		settings = &RuntimeSettings{ReadOnly: true}
	}

	cleanup := func() {}
	moduleSpecs := []engine.ModuleSpec{
		engine.DefaultRegistryModulesNamed("fs", "path", "time", "timer", "yaml"),
	}
	if settings.DBPath != "" {
		db, err := sql.Open("sqlite3", settings.DBPath)
		if err != nil {
			return nil, nil, fmt.Errorf("open sqlite database %s: %w", settings.DBPath, err)
		}
		if err := db.Ping(); err != nil {
			_ = db.Close()
			return nil, nil, fmt.Errorf("ping sqlite database %s: %w", settings.DBPath, err)
		}
		guarded := &guardedDB{db: db, allowWrites: settings.AllowWrites && !settings.ReadOnly}
		moduleSpecs = append(moduleSpecs,
			engine.NativeModuleSpec{ModuleID: "database:configured", ModuleName: "database", Loader: dbModuleLoader(guarded)},
			engine.NativeModuleSpec{ModuleID: "database:db-alias", ModuleName: "db", Loader: dbModuleLoader(guarded)},
		)
		cleanup = func() { _ = db.Close() }
	}

	builder := engine.NewBuilder(runtimeOptions(repo)...).
		WithRequireOptions(noderequire.WithLoader(repo.Registry.RequireLoader())).
		WithModules(moduleSpecs...)
	factory, err := builder.Build()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	return factory, cleanup, nil
}

func runtimeOptions(repo ScannedRepository) []engine.Option {
	if repo.Repository.Embedded {
		return nil
	}
	folders := []string{repo.Repository.RootDir, filepath.Join(repo.Repository.RootDir, "node_modules")}
	parent := filepath.Dir(repo.Repository.RootDir)
	if parent != repo.Repository.RootDir {
		folders = append(folders, parent, filepath.Join(parent, "node_modules"))
	}
	return []engine.Option{engine.WithRequireOptions(noderequire.WithGlobalFolders(folders...))}
}

type guardedDB struct {
	db          *sql.DB
	allowWrites bool
}

func dbModuleLoader(db *guardedDB) noderequire.ModuleLoader {
	return func(vm *goja.Runtime, moduleObj *goja.Object) {
		exports := moduleObj.Get("exports").(*goja.Object)
		_ = exports.Set("query", func(query string, args ...any) ([]map[string]any, error) {
			return db.query(query, args...)
		})
		_ = exports.Set("exec", func(query string, args ...any) (map[string]any, error) {
			return db.exec(query, args...)
		})
	}
}

func (g *guardedDB) query(query string, args ...any) ([]map[string]any, error) {
	rows, err := g.db.Query(query, flattenArgs(args)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	ret := []map[string]any{}
	for rows.Next() {
		vals := make([]any, len(cols))
		scan := make([]any, len(cols))
		for i := range vals {
			scan[i] = &vals[i]
		}
		if err := rows.Scan(scan...); err != nil {
			return nil, err
		}
		row := map[string]any{}
		for i, col := range cols {
			row[col] = vals[i]
		}
		ret = append(ret, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (g *guardedDB) exec(query string, args ...any) (map[string]any, error) {
	if !g.allowWrites {
		return nil, fmt.Errorf("database writes are disabled; rerun with --readonly=false --allow-writes")
	}
	result, err := g.db.Exec(query, flattenArgs(args)...)
	if err != nil {
		return nil, err
	}
	ret := map[string]any{}
	if n, err := result.RowsAffected(); err == nil {
		ret["rowsAffected"] = n
	}
	if n, err := result.LastInsertId(); err == nil {
		ret["lastInsertId"] = n
	}
	return ret, nil
}

func flattenArgs(args []any) []any {
	ret := make([]any, 0, len(args))
	for _, arg := range args {
		if slice, ok := arg.([]any); ok {
			ret = append(ret, slice...)
			continue
		}
		ret = append(ret, arg)
	}
	return ret
}
