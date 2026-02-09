package sys

type (
	IdEntity struct {
		Id       int    `json:"-"`
		PublicId string `db:"public_id" json:"id"`
	}

	Entity struct {
		IdEntity
		Name string `json:"name"`
	}
)

const (
	SqliteDB             = "sqlite3"
	SqliteReadOnlySuffix = "?mode=ro"
	PostgresDB           = "postgres"
	MySQLDB              = "mysql"
	BlankUUID            = "00000000-0000-0000-0000-000000000000"
)
