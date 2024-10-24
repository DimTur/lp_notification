package main

import (
	"context"
	"log"

	"github.com/DimTur/lp_notification/cmd/notification"
)

func main() {
	ctx := context.Background()

	cmd := notification.NewServeCmd()
	if err := cmd.ExecuteContext(ctx); err != nil {
		log.Fatalf("smth went wrong: %s", err)
	}
}
