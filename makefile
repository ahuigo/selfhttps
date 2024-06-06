install:
	go install -ldflags="-s -w -X main.BuildDate=$$(date -Iseconds)" .
