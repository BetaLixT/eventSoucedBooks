package main

import (
	_ "eventSourcedBooks/docs"
	"eventSourcedBooks/pkg/app/rest"
)

func main() {
	rest.Start(false)
}
