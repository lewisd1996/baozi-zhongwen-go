package embeds

import "embed"

//go:embed internal/sql/migrations/*.sql
var Migrations embed.FS
