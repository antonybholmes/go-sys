package sys

type (
	IdEntity struct {
		Id       int    `json:"-"`
		PublicId string `db:"public_id" json:"id,omitempty"`
	}

	Entity struct {
		IdEntity
		Name string `json:"name"`
	}
)

const (
	Sqlite3DB            = "sqlite3"
	SqliteReadOnlySuffix = "?mode=ro"
	PostgresDB           = "postgres"
	MySQLDB              = "mysql"
	BlankUUID            = "00000000-0000-0000-0000-000000000000"
)
