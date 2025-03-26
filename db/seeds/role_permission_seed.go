package seeds

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/rs/zerolog/log"
)

// rolePermissionsSeed seeds the `role_permissions` table.
func (s *Seed) rolePermissionsSeed() {
	rolePermissionsMap := map[string][]map[string]string{
		"SUPER_ADMIN_E_COMMERCE": {
			{"action": "create", "resource": "notification"},
			{"action": "view_all", "resource": "notification"},
			{"action": "create", "resource": "voucher"},
			{"action": "update", "resource": "voucher"},
			{"action": "read", "resource": "voucher"},
			{"action": "delete", "resource": "voucher"},
			{"action": "create", "resource": "users"},
			{"action": "read", "resource": "users"},
			{"action": "update", "resource": "users"},
			{"action": "delete", "resource": "users"},
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
			{"action": "read", "resource": "roles"},
			{"action": "create", "resource": "roles"},
			{"action": "update", "resource": "roles"},
			{"action": "delete", "resource": "roles"},
			{"action": "read", "resource": "permissions"},
			{"action": "read", "resource": "banner"},
			{"action": "create", "resource": "banner"},
			{"action": "update", "resource": "banner"},
			{"action": "delete", "resource": "banner"},
			{"action": "read", "resource": "expenses"},
			{"action": "create", "resource": "expenses"},
			{"action": "update", "resource": "expenses"},
			{"action": "delete", "resource": "expenses"},
			{"action": "read", "resource": "order"},
			{"action": "update", "resource": "order"},
			{"action": "read", "resource": "dashboard_ecommerce"},
		},
		"SUPER_ADMIN_WEBSITE": {
			{"action": "create", "resource": "product_content"},
			{"action": "read", "resource": "product_content"},
			{"action": "update", "resource": "product_content"},
			{"action": "delete", "resource": "product_content"},
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
		},
	}

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction for role_permissions seeding")
		return
	}
	defer rollbackOrCommit(tx, &err)

	var count int
	err = tx.Get(&count, `SELECT COUNT(*) FROM role_permissions`)
	if err != nil {
		log.Error().Err(err).Msg("Error checking role_permissions table")
		return
	}
	if count > 0 {
		log.Info().Msg("Role permissions already seeded")
		return
	}

	getRoleIDQuery := `SELECT id FROM roles WHERE name = $1`
	getPermissionIDQuery := `SELECT id FROM permissions WHERE action = $1 AND resource = $2`
	insertRolePermissionQuery := `INSERT INTO role_permissions (id, role_id, permission_id) VALUES (:id, :role_id, :permission_id)`

	for roleName, permissions := range rolePermissionsMap {
		var roleID string
		err = tx.Get(&roleID, getRoleIDQuery, roleName)
		if err != nil {
			log.Error().Err(err).Msgf("Error getting role id for role %s", roleName)
			return
		}

		for _, perm := range permissions {
			var permissionID string
			err = tx.Get(&permissionID, getPermissionIDQuery, perm["action"], perm["resource"])
			if err != nil {
				log.Error().Err(err).Msgf("Error getting permission id for action %s resource %s", perm["action"], perm["resource"])
				return
			}

			uuid, _ := utils.GenerateUUIDv7String()
			params := map[string]interface{}{
				"id":            uuid,
				"role_id":       roleID,
				"permission_id": permissionID,
			}
			_, err = tx.NamedExec(insertRolePermissionQuery, params)
			if err != nil {
				log.Error().Err(err).Msg("Error inserting into role_permissions")
				return
			}
		}
	}

	log.Info().Msg("role_permissions table seeded successfully")
}
