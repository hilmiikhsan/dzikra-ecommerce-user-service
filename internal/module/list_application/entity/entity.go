package entity

type Application struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type ApplicationPermission struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	Permissions string `db:"permissions"`
}
