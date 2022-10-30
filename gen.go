//go:build ignore

package main

//go:generate easyjson -all --lower_camel_case ./pkg/infra/clients/naga/models.go

//go:generate swag init -g cmd/server/main.go
//go:generate wire ./pkg/app/rest
