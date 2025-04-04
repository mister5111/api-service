package show

import (
	"api-service/src/storage/sqlite"
)

type showFromAlias struct {
	Alias string `json:"alias" validate:"required,alias"`
}

type getRows interface {
	ShowAlias(alias string) (sqlite.AliasTableSqlite, error)
}
type getRowsAll interface {
	ShowAll() ([]sqlite.AliasTableSqlite, error)
}
