package keyvalue

import (
	"fmt"

	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
)

type postgresDB struct {
	db          *pg.DB
	getQuery    Query
	setQuery    Query
	deleteQuery Query
}

/*
const (
	getQuery    = "SELECT value FROM ?0 WHERE key=?1 ;"
	setQuery    = "INSERT INTO ?0 (key, value) VALUES ('?1', '?2') ON CONFLICT (key) DO UPDATE SET value = ?2 ;"
	deleteQuery = "DELETE FROM ?0 WHERE key=?1 ;"
)
*/

// NewpostgresDB returns a DBConnector interface backed by a postgres DB
func NewpostgresDB(addr, user, password, dbname string) DBConnector {
	options := &pg.Options{
		Addr:     addr,
		User:     user,
		Password: password,
		Database: dbname,
	}

	db := pg.Connect(options)
	qBuilder := &pgQueryBuilder{db: db}

	return &postgresDB{
		db:          db,
		getQuery:    qBuilder.AddGet().MustBuild(),
		setQuery:    qBuilder.AddSet().MustBuild(),
		deleteQuery: qBuilder.AddDelete().MustBuild(),
	}
}

// NewQueryBuilder returns an instance of a new QueryBuilder
func (d *postgresDB) NewQueryBuilder() QueryBuilder {
	return &pgQueryBuilder{db: d.db}
}

// GetQuery returns a Query instance ready to be executed as a get statement
func (d *postgresDB) GetQuery() Query {
	return d.getQuery
}

// SetQuery returns a Query instance ready to be executed as a set statement
func (d *postgresDB) SetQuery() Query {
	return d.setQuery
}

// DeleteQuery returns a Query instance ready to be executed as a delete statement
func (d *postgresDB) DeleteQuery() Query {
	return d.deleteQuery
}

// MigrateUp performs all the steps required to bring the backing DB into an initialised state
func (d *postgresDB) MigrateUp() error {
	return d.migrate("up")
}

// MigrateDown is the inverse of MigrateUp and intended to be used in testing environment to achieve a "clean slate".
func (d *postgresDB) MigrateDown() error {
	return d.migrate("reset")
}

func (d *postgresDB) migrate(cmd string) error {
	// Migrations
	migrations.DefaultCollection.DiscoverSQLMigrations("migrations/")
	_, _, _ = migrations.Run(d.db, "init")
	oldVersion, newVersion, err := migrations.Run(d.db, cmd)
	if err != nil {
		return err
	}
	if newVersion != oldVersion {
		fmt.Printf("Migration %v: from version %d to %d\n", cmd, oldVersion, newVersion)
	} else {
		fmt.Printf("Migration %v: not needed. version is %d\n", cmd, oldVersion)
	}

	return nil
}
