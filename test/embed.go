package test

import "embed"

//go:embed data/*.json
var TestData embed.FS
