package repository

const (
	queryInsertNewExpenses = `
		INSERT INTO expenses 
		(
			name, 
			cost, 
			date
		) VALUES (?, ?, ?)
		RETURNING id, name, cost, date
	`

	queryFindListExpenses = `
		SELECT
			id,
			name,
			cost,
			date
		FROM expenses
		WHERE
			deleted_at IS NULL AND
			name ILIKE '%' || ? || '%'
		ORDER BY created_at DESC, id DESC
		LIMIT ? OFFSET ?
	`

	queryCountFindListExpenses = `
		SELECT COUNT(*)
		FROM expenses
		WHERE name ILIKE '%' || ? || '%' AND deleted_at IS NULL
	`

	queryUpdateExpenses = `
		UPDATE expenses
		SET
			name = ?,
			cost = ?,
			date = ?
		WHERE 
			id = ? 
			AND deleted_at IS NULL
		RETURNING id, name, cost, date
	`

	querySoftDeleteExpensesByID = `
		UPDATE expenses
		SET
			deleted_at = NOW()
		WHERE 
			id = ? 
			AND deleted_at IS NULL
	`

	queryFindTotalSumExpenses = `
		SELECT 
			COALESCE(SUM(cost),0) AS sum
		FROM expenses
		WHERE created_at BETWEEN ? AND ?
	`
)
