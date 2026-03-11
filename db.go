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
	// 32 MB cache size, 128 MB mmap size, and foreign keys on for good measure, turn off foreign keys
	// as they are not needed for read-only access and can cause performance issues with large datasets,
	// and turn off synchronous for better performance as well as no disk syncs needed. Safe for read-only, slightly faster query execution.
	SqliteDSN = "?mode=ro&immutable=1&_foreign_keys=OFF&_synchronous=OFF&_cache_size=-32768&_mmap_size=134217728"

	PostgresDB = "postgres"
	MySQLDB    = "mysql"
	BlankUUID  = "00000000-0000-0000-0000-000000000000"
)
