package seeds

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/rs/zerolog/log"
)

// permissionEcommerceSeed ecommerce permission seeds the `permissions` table.
func (s *Seed) permissionEcommerceSeed() {
	permissionMaps := []map[string]any{
		{
			"action":   "view_all",
			"resource": "notification",
		},
		{
			"action":   "create",
			"resource": "notification",
		},
		{
			"action":   "view_on",
			"resource": "notification",
		},
		{
			"action":   "read",
			"resource": "read_notification",
		},
		{
			"action":   "view_all",
			"resource": "notification",
		},
		{
			"action":   "create",
			"resource": "voucher",
		},
		{
			"action":   "update",
			"resource": "voucher",
		},
		{
			"action":   "read",
			"resource": "voucher",
		},
		{
			"action":   "delete",
			"resource": "voucher",
		},
		{
			"action":   "create",
			"resource": "users",
		},
		{
			"action":   "read",
			"resource": "users",
		},
		{
			"action":   "update",
			"resource": "users",
		},
		{
			"action":   "delete",
			"resource": "users",
		},
		{
			"action":   "create",
			"resource": "subcategory",
		},
		{
			"action":   "read",
			"resource": "subcategory",
		},
		{
			"action":   "update",
			"resource": "subcategory",
		},
		{
			"action":   "delete",
			"resource": "subcategory",
		},
		{
			"action":   "create",
			"resource": "category",
		},
		{
			"action":   "read",
			"resource": "category",
		},
		{
			"action":   "update",
			"resource": "category",
		},
		{
			"action":   "delete",
			"resource": "category",
		},
		{
			"action":   "create",
			"resource": "product",
		},
		{
			"action":   "read",
			"resource": "product",
		},
		{
			"action":   "update",
			"resource": "product",
		},
		{
			"action":   "delete",
			"resource": "product",
		},
		{
			"action":   "read",
			"resource": "roles",
		},
		{
			"action":   "create",
			"resource": "roles",
		},
		{
			"action":   "update",
			"resource": "roles",
		},
		{
			"action":   "delete",
			"resource": "roles",
		},
		{
			"action":   "read",
			"resource": "permissions",
		},
		{
			"action":   "read",
			"resource": "banner",
		},
		{
			"action":   "create",
			"resource": "banner",
		},
		{
			"action":   "update",
			"resource": "banner",
		},
		{
			"action":   "delete",
			"resource": "banner",
		},
		{
			"action":   "read",
			"resource": "expenses",
		},
		{
			"action":   "create",
			"resource": "expenses",
		},
		{
			"action":   "update",
			"resource": "expenses",
		},
		{
			"action":   "delete",
			"resource": "expenses",
		},
		{
			"action":   "read",
			"resource": "order",
		},
		{
			"action":   "update",
			"resource": "order",
		},
		{
			"action":   "read",
			"resource": "dashboard_ecommerce",
		},
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
