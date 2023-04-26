package server

import (
	"context"
	"homework-week-5/internal/pkg/repository"

	ticket_repo "homework-week-5/internal/pkg/repository/postgresql/tickets"
	user_repo "homework-week-5/internal/pkg/repository/postgresql/users"
)

type Server struct {
	ctx        context.Context
	userRepo   repository.UsersRepo
	ticketRepo repository.TicketRepo
}

func (s *Server) InitServer(ctx context.Context, userRepo *user_repo.UsersRepo, ticketRepo *ticket_repo.TicketRepo) {
	s.ctx = ctx
	s.userRepo = userRepo
	s.ticketRepo = ticketRepo
}
