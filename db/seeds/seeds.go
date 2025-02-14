package seeds

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type Seed struct {
	db *sqlx.DB
}

// NewSeed initializes a new Seed instance with a database connection.
func newSeed(db *sqlx.DB) Seed {
	return Seed{db: db}
}

// Execute runs the seeder for the specified table with the given number of entries.
func Execute(db *sqlx.DB, table string, total int) {
	seed := newSeed(db)
	seed.run(table, total)
}

// run handles seeding based on the table name.
func (s *Seed) run(table string, total int) {
	switch table {
	case "roles":
		s.rolesSeed()
	case "permissions-ecommerce":
		s.permissionEcommerceSeed()
	case "permissions-web":
		s.permissionWebSeed()
	case "permissions-all":
		s.permissionAllSeed()
	case "role-permissions":
		s.rolePermissionsSeed()
	case "all":
		s.rolesSeed()
		s.permissionEcommerceSeed()
		s.permissionWebSeed()
		s.permissionAllSeed()
		s.rolePermissionsSeed()
	case "delete-all":
		s.deleteAll()
	default:
		log.Warn().Msg("No seed to run")
	}
}

// rollbackOrCommit handles transaction rollback or commit based on error state.
func rollbackOrCommit(tx *sqlx.Tx, err *error) {
	if *err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			log.Error().Err(rbErr).Msg("Error rolling back transaction")
		}
	} else {
		cmErr := tx.Commit()
		if cmErr != nil {
			log.Error().Err(cmErr).Msg("Error committing transaction")
		}
	}
}
