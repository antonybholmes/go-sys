package sys

type (
	IdEntity struct {
		PublicId string `db:"public_id" json:"id,omitempty"`
		Id       int    `json:"-"`
	}

	Entity struct {
		Name string `json:"name"`
		IdEntity
	}
)

const (
	Sqlite3DB            = "sqlite3"
	SqliteReadOnlySuffix = "?mode=ro"
	PostgresDB           = "postgres"
	MySQLDB              = "mysql"
	BlankUUID            = "00000000-0000-0000-0000-000000000000"
)
