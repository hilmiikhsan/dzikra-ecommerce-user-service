package seeds

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/utils"
	"github.com/rs/zerolog/log"
)

// ApplicationPermission represents a row in the application_permissions table.
type ApplicationPermission struct {
	ID            string `db:"id"`
	PermissionID  string `db:"permission_id"`
	ApplicationID string `db:"application_id"`
}

func (s *Seed) applicationPermissionEcommerceSeed() {
	ctx := context.Background()

	// 1. Fetch the application ID for "E-COMMERCE"
	var appID string
	err := s.db.GetContext(ctx, &appID, `SELECT id FROM applications WHERE name = 'E-COMMERCE' LIMIT 1`)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching application ID for E-COMMERCE")
		return
	}

	// 2. Define the permission mappings (resource dan action) yang akan dihubungkan dengan E-COMMERCE.
	permissionMaps := []map[string]string{
		{"resource": "product", "action": "delete"},
		{"resource": "product", "action": "update"},
		{"resource": "subcategory", "action": "create"},
		{"resource": "voucher", "action": "update"},
		{"resource": "category", "action": "delete"},
		{"resource": "banner", "action": "create"},
		{"resource": "users", "action": "update"},
		{"resource": "dashboard_ecommerce", "action": "read"},
		{"resource": "product", "action": "read"},
		{"resource": "order", "action": "read"},
		{"resource": "expenses", "action": "read"},
		{"resource": "expenses", "action": "update"},
		{"resource": "users", "action": "create"},
		{"resource": "banner", "action": "read"},
		{"resource": "banner", "action": "delete"},
		{"resource": "subcategory", "action": "update"},
		{"resource": "roles", "action": "update"},
		{"resource": "category", "action": "update"},
		{"resource": "voucher", "action": "read"},
		{"resource": "order", "action": "update"},
		{"resource": "notification", "action": "view_all"},
		{"resource": "roles", "action": "read"},
		{"resource": "expenses", "action": "delete"},
		{"resource": "expenses", "action": "create"},
		{"resource": "users", "action": "read"},
		{"resource": "category", "action": "read"},
		{"resource": "banner", "action": "update"},
		{"resource": "category", "action": "create"},
		{"resource": "voucher", "action": "create"},
		{"resource": "users", "action": "delete"},
	}

	// 3. Mulai transaksi
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction for application_permissions seeding")
		return
	}
	defer rollbackOrCommit(tx, &err)

	// 4. Periksa apakah data sudah ada untuk aplikasi ini
	var count int
	err = tx.GetContext(ctx, &count, `SELECT COUNT(*) FROM application_permissions WHERE application_id = $1`, appID)
	if err != nil {
		log.Error().Err(err).Msg("Error checking application_permissions table for E-COMMERCE")
		return
	}
	if count > 0 {
		log.Info().Msg("Application permissions for E-COMMERCE already seeded")
		return
	}

	// 5. Untuk setiap mapping, cari permission ID dari tabel permissions dan masukkan record ke application_permissions.
	// Perbaiki query di sini: ambil kolom "id" karena tabel permissions memiliki kolom "id" sebagai primary key.
	queryPermissionID := `SELECT id FROM permissions WHERE resource = $1 AND action = $2 LIMIT 1`
	insertQuery := `INSERT INTO application_permissions (id, permission_id, application_id) VALUES ($1, $2, $3)`
	for _, perm := range permissionMaps {
		var permissionID string
		err = tx.GetContext(ctx, &permissionID, queryPermissionID, perm["resource"], perm["action"])
		if err != nil {
			log.Error().Err(err).Msgf("Error retrieving permission ID for resource '%s', action '%s'", perm["resource"], perm["action"])
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

	log.Info().Msg("Application permissions for E-COMMERCE seeded successfully")
}
