package beady

import "embed"

// FS embeds Beady’s HTML templates and static assets (CSS/JS) for runtime use.
//
//go:embed templates/*.html static/*.css static/*.js
var FS embed.FS
