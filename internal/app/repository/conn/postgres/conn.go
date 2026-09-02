package rcpostgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"

	"github.com/kenyako/catalog-service/internal/app/config/section"
	"github.com/kenyako/catalog-service/migration"
)

type (
	Client struct {
		_bunDB
		rawBunDB *bun.DB

		cfg section.RepositoryPostgres
	}

	_bunDB = bun.IDB
)

func NewClient(ctx context.Context, cfg section.RepositoryPostgres) (*Client, error) {
	dsnURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.Username, cfg.Password),
		Host:   cfg.Address,
		Path:   cfg.Name,
	}

	query := dsnURL.Query()
	query.Set("sslmode", "disable")

	dsnURL.RawQuery = query.Encode()
	dsn := dsnURL.String()

	log.Info().
		Str("dial_timeout", cfg.DialTimeout.String()).
		Str("read_timeout", cfg.ReadTimeout.String()).
		Str("write_timeout", cfg.WriteTimeout.String()).
		Msg("Initializing PostgreSQL connection")

	opts := []pgdriver.Option{
		pgdriver.WithDSN(dsn),
	}
	if cfg.DialTimeout > 0 {
		opts = append(opts, pgdriver.WithDialTimeout(cfg.DialTimeout))
	}
	if cfg.ReadTimeout > 0 {
		opts = append(opts, pgdriver.WithReadTimeout(cfg.ReadTimeout))
	}
	if cfg.WriteTimeout > 0 {
		opts = append(opts, pgdriver.WithWriteTimeout(cfg.WriteTimeout))
	}

	connector := pgdriver.NewConnector(opts...)

	sqlDB := sql.OpenDB(connector)
	sqlDB.SetMaxOpenConns(10)

	bunDB := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := bunDB.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("ping postgres failed: %w", err)
	}

	client := &Client{
		rawBunDB: bunDB,
		_bunDB:   bunDB,
		cfg:      cfg,
	}

	return client, nil
}

func (c *Client) GetRawBunDB() *bun.DB {
	return c.rawBunDB
}

func (c *Client) Migrate(ctx context.Context) (oldVer, newVer int64, err error) {
	migrations := migrate.NewMigrations()

	if err := migrations.Discover(migration.Postgres); err != nil {
		return 0, 0, fmt.Errorf("failed to discover migrations: %w", err)
	}

	migrator := migrate.NewMigrator(
		c.rawBunDB,
		migrations,
		migrate.WithTableName(c.cfg.MigrationTable),
		migrate.WithLocksTableName(c.cfg.MigrationTable+"_lock"),
		migrate.WithMarkAppliedOnSuccess(true),
	)

	if err := migrator.Init(ctx); err != nil {
		return 0, 0, fmt.Errorf("failed to init migration table: %w", err)
	}

	applied, err := migrator.AppliedMigrations(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get applied migrations: %w", err)
	}

	oldVer = 0
	for _, m := range applied {
		v, _ := strconv.ParseInt(m.Name, 10, 64)

		if v > oldVer {
			oldVer = v
		}
	}

	group, err := migrator.Migrate(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to apply migrations: %w", err)
	}

	newVer = oldVer
	for _, m := range group.Migrations {
		v, _ := strconv.ParseInt(m.Name, 10, 64)

		if v > newVer {
			newVer = v
		}
	}

	return oldVer, newVer, nil
}
