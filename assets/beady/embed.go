package beady

import "embed"

//go:embed templates/*.html static/*.css static/*.js
var FS embed.FS
