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
	// 32 MB cache size, 128 MB mmap size, and foreign keys on for good measure
	SqliteDSN = "?mode=ro&immutable=1&_cache_size=-32768&_foreign_keys=on&_mmap_size=134217728"

	PostgresDB = "postgres"
	MySQLDB    = "mysql"
	BlankUUID  = "00000000-0000-0000-0000-000000000000"
)
