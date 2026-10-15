// Package webui embeds the built frontend (web/ → `npm run build` writes dist/ here).
package webui

import "embed"

//go:embed all:dist
var Dist embed.FS
