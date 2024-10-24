package assets

import "embed"

//go:embed data/*.json
var StaticData embed.FS
