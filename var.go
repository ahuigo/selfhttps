package main

// go build -ldflags="-s -w -X ginapp/conf.BuildDate=$(date -Iseconds) -o selfhttps .
var (
	BuildCommitId = "000"
	BuildBranch   = ""
	BuildDate     = ""
)

