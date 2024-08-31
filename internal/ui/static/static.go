package static

import "embed"

//go:embed css/*.css
var StylesheetFiles embed.FS
