package database

import _ "embed"

//go:embed schema.sql
var Schema string

//go:embed seed.sql
var Seed string
