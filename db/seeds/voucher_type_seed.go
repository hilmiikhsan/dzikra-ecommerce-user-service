package seeds

import (
	"context"

	"github.com/rs/zerolog/log"
)

// voucherTypesSeed seeds the `voucher_types` table.
func (s *Seed) voucherTypesSeed() {
	voucherTypes := []map[string]any{
		{
			"name": "Offline Shop",
			"type": "OFFLINE_SHOP",
		},
		{
			"name": "Online Shop",
			"type": "ONLINE_SHOP",
		},
	}

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("voucherTypesSeed: Error starting transaction")
		return
	}
	defer rollbackOrCommit(tx, &err)

	var count int
	err = tx.Get(&count, `SELECT COUNT(id) FROM voucher_types`)
	if err != nil {
		log.Error().Err(err).Msg("voucherTypesSeed: Error checking voucher_types table")
		return
	}
	if count > 0 {
		log.Info().Msg("Voucher types already seeded")
		return
	}

	insertVoucherTypeQuery := `INSERT INTO voucher_types (name, type) VALUES (:name, :type)`

	for _, vt := range voucherTypes {
		_, err = tx.NamedExec(insertVoucherTypeQuery, vt)
		if err != nil {
			log.Error().Err(err).Msg("voucherTypesSeed: Error inserting voucher type")
			return
		}
	}

	log.Info().Msg("voucher_types table seeded successfully")
}
