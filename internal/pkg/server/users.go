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

type usersRequest struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Age     int    `json:"age"`
}

func HandlerUsers(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			s.GetUsers(w, s.ctx, s.userRepo, req)
		case http.MethodPost:
			s.CreateUsers(w, s.ctx, s.userRepo, req)
		case http.MethodPut:
			s.UpdateUsers(w, s.ctx, s.userRepo, req)
		case http.MethodDelete:
			s.DeleteUsers(w, s.ctx, s.userRepo, req)
		}
	}
}

func (s *Server) CreateUsers(w http.ResponseWriter, ctx context.Context, userRepo repository.UsersRepo, req *http.Request) int {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error while reading body, err: [%s]\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError) //500
		return http.StatusInternalServerError
	}
	var unmarshalled usersRequest
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		log.Printf("Error while unmarshalling body, err: [%s]\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError) //500
		return http.StatusInternalServerError
	}

	user := &repository.User{Name: unmarshalled.Name, Surname: unmarshalled.Surname, Age: unmarshalled.Age}

	id, err := userRepo.Add(ctx, user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest) //400
		return http.StatusBadRequest
	}
	log.Println(id)
	w.WriteHeader(http.StatusOK) //200
	return http.StatusOK
}

func (s *Server) GetUsers(w http.ResponseWriter, ctx context.Context, userRepo repository.UsersRepo, req *http.Request) int {
	requestId := req.URL.Query().Get("id")
	if requestId == "" {
		log.Print("You should add \"id\" parameter\n")
		w.WriteHeader(http.StatusBadRequest) //400
		return http.StatusBadRequest
	}

	if requestId == "all" {
		users, err := userRepo.List(ctx)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest) //400
			return http.StatusBadRequest
		}

		jsonUsers, err := json.Marshal(users)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest) //400
			return http.StatusBadRequest
		}
		log.Println(string(jsonUsers))
		w.WriteHeader(http.StatusOK) //200
		return http.StatusOK
	}

	id, err := strconv.Atoi(requestId)
	if err != nil {
		log.Printf("\"id\" should have INT type, err: [%s]\n", err.Error())
		w.WriteHeader(http.StatusBadRequest) //400
		return http.StatusBadRequest
	}

	user, err := userRepo.GetById(ctx, id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest) //400
		return http.StatusBadRequest
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest) //400
		return http.StatusBadRequest
	}
	log.Println(string(jsonUser))
	w.WriteHeader(http.StatusOK) //200
	return http.StatusOK
}

func (s *Server) UpdateUsers(w http.ResponseWriter, ctx context.Context, userRepo repository.UsersRepo, req *http.Request) int {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error while reading body, err: [%s]\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError) //500
		return http.StatusInternalServerError
	}
	var unmarshalled usersRequest
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		log.Printf("Error while unmarshalling body, err: [%s]\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError) //500
		return http.StatusInternalServerError
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

	user := &repository.User{ID: id, Name: unmarshalled.Name, Surname: unmarshalled.Surname, Age: unmarshalled.Age}

	updated, err := userRepo.Update(ctx, user)
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

func (s *Server) DeleteUsers(w http.ResponseWriter, ctx context.Context, userRepo repository.UsersRepo, req *http.Request) int {
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

	deleted, err := userRepo.Delete(ctx, id)
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
