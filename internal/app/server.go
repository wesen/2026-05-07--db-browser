package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/go-go-golems/go-go-goja/engine"
	databasemod "github.com/go-go-golems/go-go-goja/modules/database"
	expressmod "github.com/go-go-golems/go-go-goja/modules/express"
	"github.com/go-go-golems/go-go-goja/modules/uidsl"
	"github.com/go-go-golems/go-go-goja/pkg/gojahttp"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Addr        string
	DBPath      string
	ScriptsDir  string
	Dev         bool
	ReadOnly    bool
	AllowWrites bool
}

type Server struct {
	cfg     Config
	db      *sql.DB
	runtime *engine.Runtime
	host    *gojahttp.Host
	httpSrv *http.Server
}

func NewServer(ctx context.Context, cfg Config) (*Server, error) {
	if cfg.Addr == "" {
		cfg.Addr = ":8080"
	}
	if cfg.DBPath == "" {
		cfg.DBPath = "./app.db"
	}
	if cfg.ScriptsDir == "" {
		cfg.ScriptsDir = "./scripts"
	}
	if !cfg.AllowWrites {
		cfg.ReadOnly = true
	}

	if dir := filepath.Dir(cfg.DBPath); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create db directory: %w", err)
		}
	}
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}

	host := gojahttp.NewHost(gojahttp.HostOptions{Dev: cfg.Dev, Renderer: uidsl.RenderAny})
	guarded := &guardedDB{db: db, allowWrites: cfg.AllowWrites && !cfg.ReadOnly}
	databaseModule := databasemod.New(databasemod.WithPreconfiguredDB(guarded), databasemod.WithConfigureEnabled(false))
	dbAliasModule := databasemod.New(databasemod.WithName("db"), databasemod.WithPreconfiguredDB(guarded), databasemod.WithConfigureEnabled(false))

	factory, err := engine.NewBuilder().
		WithModules(
			engine.DefaultRegistryModulesNamed("fs", "path", "time", "timer", "yaml"),
			engine.NativeModuleSpec{ModuleID: "database:configured", ModuleName: databaseModule.Name(), Loader: databaseModule.Loader},
			engine.NativeModuleSpec{ModuleID: "database:db-alias", ModuleName: dbAliasModule.Name(), Loader: dbAliasModule.Loader},
		).
		WithRuntimeModuleRegistrars(expressmod.NewRegistrar(host), uidsl.NewRegistrar()).
		Build()
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("build goja factory: %w", err)
	}

	rt, err := factory.NewRuntime(ctx)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create goja runtime: %w", err)
	}
	host.SetRuntime(rt.Owner)

	s := &Server{cfg: cfg, db: db, runtime: rt, host: host}
	if err := s.LoadScripts(ctx); err != nil {
		_ = s.Close(context.Background())
		return nil, err
	}
	return s, nil
}

func (s *Server) Handler() http.Handler { return s.host }

func (s *Server) Run(ctx context.Context) error {
	s.httpSrv = &http.Server{Addr: s.cfg.Addr, Handler: s.host, ReadHeaderTimeout: 5 * time.Second}
	errCh := make(chan error, 1)
	go func() {
		err := s.httpSrv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()
	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.httpSrv.Shutdown(shutdownCtx)
		return nil
	case err := <-errCh:
		return err
	}
}

func (s *Server) Close(ctx context.Context) error {
	var errs []error
	if s.httpSrv != nil {
		if err := s.httpSrv.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	if s.runtime != nil {
		if err := s.runtime.Close(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (s *Server) LoadScripts(ctx context.Context) error {
	files, err := scriptFiles(s.cfg.ScriptsDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read script %s: %w", file, err)
		}
		_, err = s.runtime.Owner.Call(ctx, "load-script", func(_ context.Context, vm *goja.Runtime) (any, error) {
			_, err := vm.RunScript(file, string(data))
			return nil, err
		})
		if err != nil {
			return fmt.Errorf("execute script %s: %w", file, err)
		}
	}
	return nil
}

func scriptFiles(dir string) ([]string, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("stat scripts directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("scripts path %s is not a directory", dir)
	}
	files := []string{}
	if err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".js") {
			return nil
		}
		files = append(files, path)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("walk scripts directory: %w", err)
	}
	sort.Strings(files)
	return files, nil
}

type guardedDB struct {
	db          *sql.DB
	allowWrites bool
}

func (g *guardedDB) Query(query string, args ...any) (*sql.Rows, error) {
	return g.db.Query(query, args...)
}

func (g *guardedDB) Exec(query string, args ...any) (sql.Result, error) {
	if !g.allowWrites {
		return nil, fmt.Errorf("database writes are disabled; restart with --readonly=false --allow-writes")
	}
	return g.db.Exec(query, args...)
}
