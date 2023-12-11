package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"manager/internal/domain"
	"net/http"
)

type service interface {
	GetAll() domain.Storage
	DeleteByName(string) error
	Append(string, domain.Element) error
	UpdateByName(string, domain.Element) error
	GetByType(string) (domain.Storage, error)
}

type Handler struct {
	s service
}

func New(s service) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) InitRouter() http.Handler {
	router := http.NewServeMux()

	router.Handle("/get-by-type", h.getByType())
	router.Handle("/get-all", h.getAll())
	router.Handle("/add", h.add())
	router.Handle("/update", h.update())
	router.Handle("/delete", h.delete())

	return router
}

func (h *Handler) getByType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Print("failed to parse parameters")
			http.Error(w, "bad Request", http.StatusBadRequest)

			return
		}

		recordType := r.Form.Get("type")

		storage, err := h.s.GetByType(recordType)
		if err != nil {
			log.Printf("bad request: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		storageJSON, err := json.Marshal(storage)
		if err != nil {
			log.Printf("failed to marshal storage: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(storageJSON)
	}
}

func (h *Handler) getAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		storage := h.s.GetAll()

		storageJSON, err := json.Marshal(storage)
		if err != nil {
			log.Printf("failed to marshal storage: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(storageJSON)
	}
}

func (h *Handler) add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Print("failed to parse parameters")
			http.Error(w, "bad Request", http.StatusBadRequest)

			return
		}

		name := r.Form.Get("name")

		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, "error reading request body", http.StatusInternalServerError)

			return
		}

		var elem domain.Element
		err = json.Unmarshal(body, &elem)
		if err != nil {
			http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)

			return
		}

		err = h.s.Append(name, elem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Print("failed to parse parameters")
			http.Error(w, "bad Request", http.StatusBadRequest)

			return
		}

		name := r.Form.Get("name")

		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, "error reading request body", http.StatusInternalServerError)

			return
		}

		var elem domain.Element
		err = json.Unmarshal(body, &elem)
		if err != nil {
			http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)

			return
		}

		err = h.s.UpdateByName(name, elem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Print("failed to parse parameters")
			http.Error(w, "bad Request", http.StatusBadRequest)

			return
		}

		name := r.Form.Get("name")

		err := h.s.DeleteByName(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
