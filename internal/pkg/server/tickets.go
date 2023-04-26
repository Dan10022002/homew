package server

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"homework-week-5/internal/pkg/repository"
)

type ticketsRequest struct {
	UserID int `json:"user_id"`
	Cost   int `json:"cost"`
	Place  int `json:"place"`
}

func HandlerTickets(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			s.GetTickets(w, s.ctx, s.userRepo, s.ticketRepo, req)
		case http.MethodPost:
			s.CreateTickets(w, s.ctx, s.userRepo, s.ticketRepo, req)
		case http.MethodPut:
			s.UpdateTickets(w, s.ctx, s.userRepo, s.ticketRepo, req)
		case http.MethodDelete:
			s.DeleteTickets(w, s.ctx, s.ticketRepo, req)
		}
	}
}

func (s *Server) CreateTickets(w http.ResponseWriter, ctx context.Context, userRepo repository.UsersRepo, ticketRepo repository.TicketRepo, req *http.Request) int {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error while reading body, err: [%s]\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError) //500
		return http.StatusInternalServerError
	}
	var unmarshalled ticketsRequest
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		log.Printf("Error while unmarshalling body, err: [%s]\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError) //500
		return http.StatusInternalServerError
	}

	ticket := &repository.Ticket{UserID: unmarshalled.UserID, Cost: unmarshalled.Cost, Place: unmarshalled.Place}

	_, err = userRepo.GetById(ctx, ticket.UserID)
	if err == repository.ErrObjectNotFound {
		log.Println("User with giving ID doesn't exist")
		w.WriteHeader(http.StatusConflict) //409
		return http.StatusConflict
	}

	id, err := ticketRepo.Add(ctx, ticket)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest) //400
		return http.StatusBadRequest
	}
	log.Println(id)
	w.WriteHeader(http.StatusOK) //200
	return http.StatusOK
}

func (s *Server) GetTickets(w http.ResponseWriter, ctx context.Context, userRepo repository.UsersRepo,
	ticketRepo repository.TicketRepo, req *http.Request) int {
	requestUserId := req.URL.Query().Get("user_id")
	if requestUserId != "" {
		userId, err := strconv.Atoi(requestUserId)
		if err != nil {
			log.Printf("\"user_id\" should have INT type, err: [%s]\n", err.Error())
			w.WriteHeader(http.StatusBadRequest) //400
			return http.StatusBadRequest
		}

		_, err = userRepo.GetById(ctx, userId)
		if err == repository.ErrObjectNotFound {
			log.Println("User with giving ID doesn't exist")
			w.WriteHeader(http.StatusConflict) //409
			return http.StatusConflict
		}

		ticket, err := ticketRepo.GetByUserId(ctx, userId)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest) //400
			return http.StatusBadRequest
		}

		jsonTicket, err := json.Marshal(ticket)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest) //400
			return http.StatusBadRequest
		}
		log.Println(string(jsonTicket))
		w.WriteHeader(http.StatusOK) //200
		return http.StatusOK
	}

	requestId := req.URL.Query().Get("id")
	if requestId == "" {
		log.Print("You should add \"id\" parameter\n")
		w.WriteHeader(http.StatusBadRequest) //400
		return http.StatusBadRequest
	}

	if requestId == "all" {
		tickets, err := ticketRepo.List(ctx)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest) //400
			return http.StatusBadRequest
		}

		jsonTickets, err := json.Marshal(tickets)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest) //400
			return http.StatusBadRequest
		}
		log.Println(string(jsonTickets))
		w.WriteHeader(http.StatusOK) //200
		return http.StatusOK
	}

	id, err := strconv.Atoi(requestId)
	if err != nil {
		log.Printf("\"id\" should have INT type, err: [%s]\n", err.Error())
		w.WriteHeader(http.StatusBadRequest) //400
		return http.StatusBadRequest
	}

	ticket, err := ticketRepo.GetById(ctx, id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest) //400
		return http.StatusBadRequest
	}

	jsonTicket, err := json.Marshal(ticket)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest) //400
		return http.StatusBadRequest
	}
	log.Println(string(jsonTicket))
	w.WriteHeader(http.StatusOK) //200
	return http.StatusOK
}

func (s *Server) UpdateTickets(w http.ResponseWriter, ctx context.Context, userRepo repository.UsersRepo,
	ticketRepo repository.TicketRepo, req *http.Request) int {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error while reading body, err: [%s]\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError) //500
		return http.StatusInternalServerError
	}
	var unmarshalled ticketsRequest
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		log.Printf("Error while unmarshalling body, err: [%s]\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError) //500
		return http.StatusInternalServerError
	}

	_, err = userRepo.GetById(ctx, unmarshalled.UserID)
	if err == repository.ErrObjectNotFound {
		log.Println("User with giving ID doesn't exist")
		w.WriteHeader(http.StatusConflict) //409
		return http.StatusConflict
	}

	requestId := req.URL.Query().Get("id")
	if requestId == "" {
		log.Print("You should add \"id\" parameter\n")
		w.WriteHeader(http.StatusBadRequest) //400
		return http.StatusBadRequest
	}

	id, err := strconv.Atoi(requestId)
	if err != nil {
		log.Printf("\"id\" should have INT type, err: [%s]\n", err.Error())
		w.WriteHeader(http.StatusBadRequest) //400
		return http.StatusBadRequest
	}

	ticket := &repository.Ticket{ID: id, UserID: unmarshalled.UserID, Cost: unmarshalled.Cost, Place: unmarshalled.Place}

	updated, err := ticketRepo.Update(ctx, ticket)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest) //400
		return http.StatusBadRequest
	}
	if updated {
		w.WriteHeader(http.StatusOK) //200
		return http.StatusOK
	} else {
		w.WriteHeader(http.StatusBadRequest) //400
		return http.StatusBadRequest
	}
}

func (s *Server) DeleteTickets(w http.ResponseWriter, ctx context.Context, ticketRepo repository.TicketRepo, req *http.Request) int {
	requestId := req.URL.Query().Get("id")
	if requestId == "" {
		log.Print("You should add \"id\" parameter\n")
		w.WriteHeader(http.StatusBadRequest) //400
		return http.StatusBadRequest
	}
	id, err := strconv.Atoi(requestId)
	if err != nil {
		log.Printf("\"id\" should have INT type, err: [%s]\n", err.Error())
		w.WriteHeader(http.StatusBadRequest) //400
		return http.StatusBadRequest
	}

	deleted, err := ticketRepo.Delete(ctx, id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest) //400
		return http.StatusBadRequest
	}
	if deleted {
		w.WriteHeader(http.StatusOK) //200
		return http.StatusOK
	} else {
		w.WriteHeader(http.StatusBadRequest) //400
		return http.StatusBadRequest
	}
}
