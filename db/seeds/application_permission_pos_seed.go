package seeds

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/rs/zerolog/log"
)

func (s *Seed) applicationPermissionPOSSeed() {
	ctx := context.Background()

	// 1. Fetch the application ID for "POS"
	var appID string
	err := s.db.GetContext(ctx, &appID, `SELECT id FROM applications WHERE name = 'POS' LIMIT 1`)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching application ID for POS")
		return
	}

	// 2. Define the permission mappings for the POS application.
	permissionMaps := []map[string]string{
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

	// 3. Start a transaction
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction for application_permissions seeding")
		return
	}
	defer rollbackOrCommit(tx, &err)

	// 4. Check if there are already permissions for this application.
	var count int
	err = tx.GetContext(ctx, &count, `SELECT COUNT(*) FROM application_permissions WHERE application_id = $1`, appID)
	if err != nil {
		log.Error().Err(err).Msg("Error checking application_permissions table for POS")
		return
	}
	if count > 0 {
		log.Info().Msg("Application permissions for POS already seeded")
		return
	}

	// 5. For each permission mapping, fetch the permission id from the permissions table and insert a record.
	//    Di sini, permission id diambil dari kolom "id" pada tabel permissions.
	queryPermissionID := `SELECT id FROM permissions WHERE resource = $1 AND action = $2 LIMIT 1`
	insertQuery := `INSERT INTO application_permissions (id, permission_id, application_id) VALUES ($1, $2, $3)`
	for _, perm := range permissionMaps {
		var permissionID string
		err = tx.GetContext(ctx, &permissionID, queryPermissionID, perm["resource"], perm["action"])
		if err != nil {
			log.Error().Err(err).Msgf("Error retrieving permission id for resource '%s', action '%s'", perm["resource"], perm["action"])
			return
		}

		// Generate UUID untuk kolom id secara manual
		newID, err := utils.GenerateUUIDv7String()
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate UUID for application_permissions record")
			return
		}

		_, err = tx.ExecContext(ctx, insertQuery, newID, permissionID, appID)
		if err != nil {
			log.Error().Err(err).Msgf("Error inserting permission (resource: %s, action: %s)", perm["resource"], perm["action"])
			return
		}
	}

	log.Info().Msg("Application permissions for POS seeded successfully")
}
