package seeds

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/rs/zerolog/log"
)

// applicationPermissionWebSeed seeds the application_permissions table for the "WEB" application.
func (s *Seed) applicationPermissionWebSeed() {
	ctx := context.Background()

	// 1. Fetch the application ID for "WEB"
	var appID string
	err := s.db.GetContext(ctx, &appID, `SELECT id FROM applications WHERE name = 'WEB' LIMIT 1`)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching application ID for WEB")
		return
	}

	// 2. Define permission mappings untuk aplikasi WEB.
	// Data ini harus sesuai dengan entri yang sudah ada di tabel permissions.
	permissionMaps := []map[string]string{
		{"resource": "product_content", "action": "create"},
		{"resource": "category_article", "action": "create"},
		{"resource": "article", "action": "create"},
		{"resource": "category_article", "action": "read"},
		{"resource": "faq", "action": "update"},
		{"resource": "faq", "action": "delete"},
		{"resource": "dashboard", "action": "read"},
		{"resource": "article", "action": "read"},
		{"resource": "article", "action": "update"},
		{"resource": "faq", "action": "create"},
		{"resource": "banner", "action": "read"},
		{"resource": "category_article", "action": "delete"},
		{"resource": "faq", "action": "read"},
		{"resource": "category_article", "action": "update"},
		{"resource": "article", "action": "delete"},
		{"resource": "product_content", "action": "delete"},
		{"resource": "product_content", "action": "read"},
	}

	// 3. Mulai transaksi
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction for application_permissions seeding")
		return
	}
	defer rollbackOrCommit(tx, &err)

	// 4. Periksa apakah sudah ada data untuk aplikasi WEB
	var count int
	err = tx.GetContext(ctx, &count, `SELECT COUNT(*) FROM application_permissions WHERE application_id = $1`, appID)
	if err != nil {
		log.Error().Err(err).Msg("Error checking application_permissions table for WEB")
		return
	}
	if count > 0 {
		log.Info().Msg("Application permissions for WEB already seeded")
		return
	}

	// 5. Untuk setiap permission mapping, cari permission ID dari tabel permissions dan masukkan record.
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

	log.Info().Msg("Application permissions for WEB seeded successfully")
}
