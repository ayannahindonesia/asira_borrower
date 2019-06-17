package tests

import "asira/migration"

func RebuildData() {
	migration.Truncate([]string{"all"})
	migration.Seed()
}
