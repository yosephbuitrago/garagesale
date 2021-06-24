SHELL=/bin/bash

run:
	go run app/sales-api/main.go

runa:
	go run app/admin/admin.go

tidy:
	go mod tidy
	go mod vendor

test:
	go test -v ./...
	# statickcheck