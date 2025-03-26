package seeds

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/rs/zerolog/log"
)

// rolesSeed seeds the `roles` table.
func (s *Seed) applicationSeed() {
	roleMaps := []map[string]any{
		{
			"name": "E-COMMERCE",
		},
		{
			"name": "WEB",
		},
		{
			"name": "POS",
		},
	}

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer rollbackOrCommit(tx, &err)

	// Check if applications already exist
	var count int
	err = tx.Get(&count, `SELECT COUNT(id) FROM applications`)
	if err != nil {
		log.Error().Err(err).Msg("Error checking applications table")
		return
	}
	if count > 0 {
		log.Info().Msg("Applications already seeded")
		return
	}

	insertRoleQuery := `INSERT INTO applications (id, name) VALUES (:id, :name)`
	for _, role := range roleMaps {
		uuid, _ := utils.GenerateUUIDv7String()
		role["id"] = uuid
		_, err = tx.NamedExec(insertRoleQuery, role)
		if err != nil {
			log.Error().Err(err).Msg("Error inserting applications")
			return
		}
	}

	log.Info().Msg("applications table seeded successfully")
}
