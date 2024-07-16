
run:
	go run cmd/main.go

SWAG_BIN=$(GOPATH)/bin/swag

swag-gen:
	$(SWAG_BIN) init -g ./api/api.go -o api/docs force 1