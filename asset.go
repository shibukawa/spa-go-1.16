package main

import (
	"embed"
)

//go:embed frontend/dist/*
var assets embed.FS
