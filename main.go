package main

import (
	"embed"
	"pgsync/cmd"
)

//go:embed branding/logo.txt
var logoFS embed.FS

func main() {
	logoData, _ := logoFS.ReadFile("branding/logo.txt")
	cmd.Execute(string(logoData))
}
