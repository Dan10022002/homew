package main

import (
	"context"
	"fmt"
	"homework-week-5/internal/pkg/db"
	ticket_repo "homework-week-5/internal/pkg/repository/postgresql/tickets"
	user_repo "homework-week-5/internal/pkg/repository/postgresql/users"
	server "homework-week-5/internal/pkg/server"
	"log"
	"net/http"
)

func serverMux(s *server.Server) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users", server.HandlerUsers(s))
	mux.HandleFunc("/tickets", server.HandlerTickets(s))

	if err := http.ListenAndServe(":9000", mux); err != nil {
		log.Fatal(err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	database, err := db.NewDB(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	userRepo := user_repo.NewUsers(database)
	ticketRepo := ticket_repo.NewTickets(database)

	defer database.GetPool(ctx).Close()

	implementation := server.Server{}
	implementation.InitServer(ctx, userRepo, ticketRepo)

	serverMux(&implementation)
}
