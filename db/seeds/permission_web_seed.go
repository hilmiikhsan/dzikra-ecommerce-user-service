package seeds

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/utils"
	"github.com/rs/zerolog/log"
)

// permissionWebSeed permission web seeds the `permissions` table.
func (s *Seed) permissionWebSeed() {
	permissionMaps := []map[string]any{
		{
			"action":   "create",
			"resource": "product_content",
		},
		{
			"action":   "read",
			"resource": "product_content",
		},
		{
			"action":   "update",
			"resource": "product_content",
		},
		{
			"action":   "delete",
			"resource": "product_content",
		},
		{
			"action":   "create",
			"resource": "article",
		},
		{
			"action":   "read",
			"resource": "article",
		},
		{
			"action":   "update",
			"resource": "article",
		},
		{
			"action":   "delete",
			"resource": "article",
		},
		{
			"action":   "create",
			"resource": "faq",
		},
		{
			"action":   "read",
			"resource": "faq",
		},
		{
			"action":   "update",
			"resource": "faq",
		},
		{
			"action":   "delete",
			"resource": "faq",
		},
		{
			"action":   "create",
			"resource": "category_article",
		},
		{
			"action":   "read",
			"resource": "category_article",
		},
		{
			"action":   "update",
			"resource": "category_article",
		},
		{
			"action":   "delete",
			"resource": "category_article",
		},
		{
			"action":   "read",
			"resource": "dashboard",
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
