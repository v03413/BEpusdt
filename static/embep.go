package static

import (
	"embed"
)

//go:embed secure/*
var Secure embed.FS

//go:embed payment/*
var Payment embed.FS
