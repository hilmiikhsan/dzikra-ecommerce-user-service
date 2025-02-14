package seeds

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/utils"
	"github.com/rs/zerolog/log"
)

// permissionAllSeed seeds all permissions into the `permissions` table.
func (s *Seed) permissionAllSeed() {
	permissionMaps := []map[string]any{
		{"action": "read", "resource": "profile"},
		{"action": "update", "resource": "profile"},
		{"action": "delete", "resource": "profile"},

		{"action": "create", "resource": "address"},
		{"action": "read", "resource": "address"},
		{"action": "update", "resource": "address"},
		{"action": "delete", "resource": "address"},

		{"action": "create", "resource": "notification"},
		{"action": "view_on", "resource": "notification"},
		{"action": "read", "resource": "read_notification"},
		{"action": "view_all", "resource": "notification"},

		{"action": "create", "resource": "voucher"},
		{"action": "update", "resource": "voucher"},
		{"action": "read", "resource": "voucher"},
		{"action": "delete", "resource": "voucher"},
		{"action": "use", "resource": "voucher"},

		{"action": "create", "resource": "users"},
		{"action": "read", "resource": "users"},
		{"action": "update", "resource": "users"},
		{"action": "delete", "resource": "users"},

		{"action": "create", "resource": "roles"},
		{"action": "read", "resource": "roles"},
		{"action": "update", "resource": "roles"},
		{"action": "delete", "resource": "roles"},

		{"action": "read", "resource": "banner"},
		{"action": "create", "resource": "banner"},
		{"action": "update", "resource": "banner"},
		{"action": "delete", "resource": "banner"},

		{"action": "create", "resource": "subcategory"},
		{"action": "read", "resource": "subcategory"},
		{"action": "update", "resource": "subcategory"},
		{"action": "delete", "resource": "subcategory"},

		{"action": "create", "resource": "category"},
		{"action": "read", "resource": "category"},
		{"action": "update", "resource": "category"},
		{"action": "delete", "resource": "category"},

		{"action": "create", "resource": "product"},
		{"action": "read", "resource": "product"},
		{"action": "update", "resource": "product"},
		{"action": "delete", "resource": "product"},

		{"action": "create", "resource": "product_content"},
		{"action": "read", "resource": "product_content"},
		{"action": "update", "resource": "product_content"},
		{"action": "delete", "resource": "product_content"},

		{"action": "read", "resource": "permissions"},

		{"action": "create", "resource": "article"},
		{"action": "read", "resource": "article"},
		{"action": "update", "resource": "article"},
		{"action": "delete", "resource": "article"},

		{"action": "create", "resource": "faq"},
		{"action": "read", "resource": "faq"},
		{"action": "update", "resource": "faq"},
		{"action": "delete", "resource": "faq"},

		{"action": "create", "resource": "category_article"},
		{"action": "read", "resource": "category_article"},
		{"action": "update", "resource": "category_article"},
		{"action": "delete", "resource": "category_article"},

		{"action": "read", "resource": "dashboard"},

		{"action": "read", "resource": "expenses"},
		{"action": "create", "resource": "expenses"},
		{"action": "update", "resource": "expenses"},
		{"action": "delete", "resource": "expenses"},

		{"action": "read", "resource": "order"},
		{"action": "update", "resource": "order"},

		{"action": "read", "resource": "dashboard_ecommerce"},
	}

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer rollbackOrCommit(tx, &err)

	// Cek apakah data permission sudah ada
	var count int
	err = tx.Get(&count, `SELECT COUNT(id) FROM permissions`)
	if err != nil {
		log.Error().Err(err).Msg("Error checking permissions table")
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
