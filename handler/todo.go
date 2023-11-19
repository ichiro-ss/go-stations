package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// ServeHTTP implements http.Handler interface.
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// todo_response := &model.TODO{}
	// POST
	if r.Method == http.MethodPost {
		var createTodoReq model.CreateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&createTodoReq); err != nil {
			log.Println(err)
			return
		}
		if createTodoReq.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			createTodoRes, err := h.Create(r.Context(), &createTodoReq)
			if err != nil {
				log.Println(err)
				return
			}
			if err := json.NewEncoder(w).Encode(*createTodoRes); err != nil {
				log.Println(err)
				return
			}
		}
	}
	// PUT
	if r.Method == http.MethodPut {
		var updateTodoReq model.UpdateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&updateTodoReq); err != nil {
			log.Println(err)
			return
		}
		if updateTodoReq.Subject == "" || updateTodoReq.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			updateTodoRes, err := h.Update(r.Context(), &updateTodoReq)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if err := json.NewEncoder(w).Encode(*updateTodoRes); err != nil {
				log.Println(err)
				return
			}
		}
	}
	// GET
	if r.Method == http.MethodGet {
		var readTodoReq model.ReadTODORequest
		query := r.URL.Query()
		prevID, present := query["prev_id"]
		if present && len(prevID) != 0 {
			readTodoReq.PrevID, _ = strconv.ParseInt(prevID[0], 10, 64)
		}

		Size, present := query["size"]
		if present && len(Size) != 0 {
			readTodoReq.Size, _ = strconv.ParseInt(Size[0], 10, 64)
		} else {
			readTodoReq.Size = 5
		}

		readTodoRes, err := h.Read(r.Context(), &readTodoReq)
		if err != nil {
			log.Println(err)
			return
		}
		if err := json.NewEncoder(w).Encode(*readTodoRes); err != nil {
			log.Println(err)
			return
		}
	}
	// DELETE
	if r.Method == http.MethodDelete {
		var deleteTodoReq model.DeleteTODORequest
		if err := json.NewDecoder(r.Body).Decode(&deleteTodoReq); err != nil {
			log.Println(err)
			return
		}
		if len(deleteTodoReq.IDs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
		}
		_, err := h.Delete(r.Context(), &deleteTodoReq)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	var createTodoRes model.CreateTODOResponse
	createTodoRes.TODO = *todo
	return &createTodoRes, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}
	return &model.ReadTODOResponse{TODOs: todos}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	var updateTodoRes model.UpdateTODOResponse
	updateTodoRes.TODO = *todo
	return &updateTodoRes, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	err := h.svc.DeleteTODO(ctx, req.IDs)
	if err != nil {
		return nil, err
	}
	return &model.DeleteTODOResponse{}, nil
}
