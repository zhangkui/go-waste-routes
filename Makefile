.PHONY: test integration-test verify-go-size run

test:
\tGOCACHE=$(CURDIR)/.gocache go test -mod=mod ./...

integration-test:
\tdocker compose -f docker-compose.yml up -d mysql redis

verify-go-size:
\tpowershell -ExecutionPolicy Bypass -File scripts/verify-go-size.ps1

run:
\tGOCACHE=$(CURDIR)/.gocache go run ./cmd/api
