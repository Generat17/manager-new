package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"manager/internal/domain"
	"net/http"
)

type service interface {
	WriteStorageFromFile() error
	UpdateFile() error
	GetAll() domain.Storage
	GetByType(recordType string) (domain.Storage, error)
	AppendService(serviceName string, serviceType string, favorite bool) error
	UpdateService(serviceName string, serviceType string, favorite bool) error
	DeleteService(serviceName string, login string) error
	AppendLogin(serviceName string, login string, elem domain.Element) error
	UpdateLogin(serviceName string, login string, elem domain.Element) error
	DeleteLogin(serviceName string, login string) error
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
	router.Handle("/add-login", h.addLogin())
	router.Handle("/update-login", h.updateLogin())
	router.Handle("/delete-login", h.deleteLogin())
	router.Handle("/add-service", h.addService())
	router.Handle("/update-service", h.updateService())
	router.Handle("/delete-service", h.deleteService())

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

func (h *Handler) addLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Print("failed to parse parameters")
			http.Error(w, "bad Request", http.StatusBadRequest)

			return
		}

		serviceName := r.Form.Get("name")

		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, "error reading request body", http.StatusInternalServerError)

			return
		}

		var requestBody domain.LoginBody
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)

			return
		}

		err = h.s.AppendLogin(serviceName, requestBody.Login, requestBody.Element)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) updateLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Print("failed to parse parameters")
			http.Error(w, "bad Request", http.StatusBadRequest)

			return
		}

		serviceName := r.Form.Get("name")

		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, "error reading request body", http.StatusInternalServerError)

			return
		}

		var requestBody domain.LoginBody
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)

			return
		}

		err = h.s.UpdateLogin(serviceName, requestBody.Login, requestBody.Element)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) deleteLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Print("failed to parse parameters")
			http.Error(w, "bad Request", http.StatusBadRequest)

			return
		}

		serviceName := r.Form.Get("name")
		login := r.Form.Get("login")

		err := h.s.DeleteLogin(serviceName, login)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) addService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Print("failed to parse parameters")
			http.Error(w, "bad Request", http.StatusBadRequest)

			return
		}

		serviceName := r.Form.Get("name")
		serviceType := r.Form.Get("type")
		serviceFavorite := r.Form.Get("favorite")

		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, "error reading request body", http.StatusInternalServerError)

			return
		}

		var elem domain.Service
		err = json.Unmarshal(body, &elem)
		if err != nil {
			http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)

			return
		}

		err = h.s.AppendService(serviceName, serviceType, serviceFavorite)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) updateService() http.HandlerFunc {
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

		var elem domain.Service
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

func (h *Handler) deleteService() http.HandlerFunc {
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
