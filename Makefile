-include .env

.PHONY:
	build run test cover swagger

swagger:
	@echo "  >  Generating swagger documentation..."
	go get -u github.com/swaggo/swag/cmd/swag
	swag i -o swagger

build:
	@echo "  >  Building package..."
	go build -o app github.com/f4nt0md3v/gymshark-packs

run:
	@echo "  >  Running package..."
	go run github.com/f4nt0md3v/gymshark-packs
