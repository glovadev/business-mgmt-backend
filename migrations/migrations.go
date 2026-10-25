package migrations

import "embed"

// FS embeds all SQL migration files (*.sql) directly into the compiled Go binary.
//
//go:embed *.sql
var FS embed.FS
