package main

import "github.com/benjamingeer/accicalc/cmd"

var (
	// ldflags set by GoReleaser
	version = "unknown"
	commit  = "unknown"
)

func main() {
	cmd.Version = version
	cmd.Commit = commit

	cmd.Execute()
}
