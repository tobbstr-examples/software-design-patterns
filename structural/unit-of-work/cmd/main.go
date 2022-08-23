package main

import (
	"database/sql"

	"github.com/tobbstr-examples/business-logic-patterns/structural/unit-of-work/editor"
	"github.com/tobbstr-examples/business-logic-patterns/structural/unit-of-work/store"
)

func main() {
	// Wires the components together
	db := sql.DB{}
	uowDoer := store.NewUoWDoer(db)
	editor.NewApplicationService(uowDoer)
}
