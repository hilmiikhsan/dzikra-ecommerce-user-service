package seeds

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/utils"
	"github.com/rs/zerolog/log"
)

// permissionPosSeed ecommerce permission seeds the `permissions` table.
func (s *Seed) permissionPosSeed() {
	permissionMaps := []map[string]any{
		{"resource": "recipe", "action": "create"},
		{"resource": "product_pos", "action": "create"},
		{"resource": "product_pos", "action": "delete"},
		{"resource": "expenses_pos", "action": "read"},
		{"resource": "product_pos", "action": "update"},
		{"resource": "ingredient", "action": "create"},
		{"resource": "dashboard_pos", "action": "read"},
		{"resource": "product_category_pos", "action": "delete"},
		{"resource": "expenses_pos", "action": "update"},
		{"resource": "member", "action": "create"},
		{"resource": "ingredient", "action": "delete"},
		{"resource": "ingredient", "action": "update"},
		{"resource": "product_category_pos", "action": "read"},
		{"resource": "member", "action": "update"},
		{"resource": "recipe", "action": "update"},
		{"resource": "ingredient", "action": "read"},
		{"resource": "member", "action": "delete"},
		{"resource": "product_pos", "action": "read"},
		{"resource": "order_pos", "action": "create"},
		{"resource": "expenses_pos", "action": "delete"},
		{"resource": "product_category_pos", "action": "update"},
		{"resource": "transaction_history_pos", "action": "read"},
		{"resource": "product_category_pos", "action": "create"},
		{"resource": "expenses_pos", "action": "create"},
		{"resource": "recipe", "action": "read"},
		{"resource": "member", "action": "read"},
		{"resource": "recipe", "action": "delete"},
		{"resource": "transaction_history_pos", "action": "update"},
	}

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer rollbackOrCommit(tx, &err)

	// Check if permissions already exist
	var count int
	err = tx.Get(&count, `SELECT COUNT(id) FROM permissions`)
	if err != nil {
		log.Error().Err(err).Msg("Error checking permissions table")
		return
	}
	if count > 0 {
		log.Info().Msg("Permissions already seeded")
		return
	}

	insertPermissionQuery := `INSERT INTO permissions (id, action, resource) VALUES (:id, :action, :resource)`
	for _, permission := range permissionMaps {
		uuid, _ := utils.GenerateUUIDv7String()
		permission["id"] = uuid
		_, err = tx.NamedExec(insertPermissionQuery, permission)
		if err != nil {
			log.Error().Err(err).Msg("Error inserting permissions")
			return
		}
	}

	log.Info().Msg("permissions table seeded successfully")
}
