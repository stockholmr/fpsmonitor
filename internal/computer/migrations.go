package computer

import "github.com/jmoiron/sqlx"

func Migrate(db *sqlx.DB) error {
	migrations := []string{

		`CREATE TABLE computers (
            "id" INTEGER,
            "created" TEXT,
            "updated" TEXT,
            "deleted" TEXT,
            "name" TEXT,
            PRIMARY KEY("id" AUTOINCREMENT)
        )`,

		`CREATE TABLE computer_network_adapters (
            "id" INTEGER,
            "created" TEXT,
            "updated" TEXT,
            "deleted" TEXT,
            "computer_id" INTEGER NOT NULL,
			"name" TEXT,
            "mac_address" TEXT,
			"ip_address" TEXT,
			FOREIGN KEY("computer_id") REFERENCES "computers"("id") ON DELETE CASCADE ON UPDATE NO ACTION,
            PRIMARY KEY("id" AUTOINCREMENT)
        )`,

		`CREATE TABLE computer_users (
            "id" INTEGER,
            "created" TEXT,
            "updated" TEXT,
            "deleted" TEXT,
            "computer_id" INTEGER NOT NULL,
			"username" TEXT,
			FOREIGN KEY("computer_id") REFERENCES "computers"("id") ON DELETE CASCADE ON UPDATE NO ACTION,
            PRIMARY KEY("id" AUTOINCREMENT)
        )`,
	}

	for _, s := range migrations {
		_, err := db.Exec(s)
		if err != nil {
			return err
		}
	}

	return nil
}
