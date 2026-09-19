package loyalty

import "embed"

//go:embed web/templates
var TemplatesFS embed.FS

//go:embed migrations/*.sql
var MigrationsFS embed.FS
