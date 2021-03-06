package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/kpacha/load-test/db"
)

func main() {
	storePath := flag.String("f", ".", "path to use as store")
	port := flag.Int("p", 7879, "port to expose the html ui")
	inMemory := flag.Bool("m", false, "use in-memory store instead of the fs persistent one")
	flag.Parse()

	var store db.DB
	if *inMemory {
		store = db.NewInMemory()
	} else {
		store = db.NewFS(*storePath)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-quit
		cancel()
	}()

	server, err := NewServer(gin.New(), store, NewExecutor(store))
	if err != nil {
		fmt.Println("error building the server:", err.Error())
		return
	}

	fmt.Println(server.Run(ctx, fmt.Sprintf(":%d", *port)))
}
