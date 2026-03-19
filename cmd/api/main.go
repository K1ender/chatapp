package main

import (
	"chatapp/pkg/api"
	"context"
)

func main() {
	ctx := context.Background()

	err := api.Run(ctx)
	if err != nil {
		panic(err)
	}
}
