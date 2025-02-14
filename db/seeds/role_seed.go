package seeds

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/utils"
	"github.com/rs/zerolog/log"
)

// rolesSeed seeds the `roles` table.
func (s *Seed) rolesSeed() {
	roleMaps := []map[string]any{
		{
			"name":        "USER",
			"description": "Just a normal User",
		},
		{
			"name":        "SUPER_ADMIN_E_COMMERCE",
			"description": "User who manage a e-commerce",
		},
		{
			"name":        "SUPER_ADMIN_WEBSITE",
			"description": "User who manage a Website Company",
		},
		{
			"name":        "SUPER_ADMIN_POS",
			"description": "User who manage a Cashier App",
		},
		{
			"name":        "MONITOR",
			"description": "User just monitoring the website",
		},
	}

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer rollbackOrCommit(tx, &err)

	// Check if roles already exist
	var count int
	err = tx.Get(&count, `SELECT COUNT(id) FROM roles`)
	if err != nil {
		log.Error().Err(err).Msg("Error checking roles table")
		return
	}
	if count > 0 {
		log.Info().Msg("Roles already seeded")
		return
	}

	insertRoleQuery := `INSERT INTO roles (id, name, description) VALUES (:id, :name, :description)`
	for _, role := range roleMaps {
		uuid, _ := utils.GenerateUUIDv7String()
		role["id"] = uuid
		_, err = tx.NamedExec(insertRoleQuery, role)
		if err != nil {
			log.Error().Err(err).Msg("Error inserting roles")
			return
		}
	}

	log.Info().Msg("roles table seeded successfully")
}
