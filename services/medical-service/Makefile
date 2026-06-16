.PHONY: setup wire run swag

setup:
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/swaggo/swag/cmd/swag@latest

wire:
	wire

run:
	go run cmd/api/*.go

swag:
	swag init -g cmd/api/main.go
