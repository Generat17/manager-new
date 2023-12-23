package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"manager/internal/domain"
)

type repository interface {
	SetStorage(domain.Storage)
	Get(string) (domain.Service, bool)
	GetAll() domain.Storage
	UpdateLogin(string, string, domain.Element) bool
	AppendLogin(string, string, domain.Element) bool
	DeleteLogin(string, string) bool
	UpdateService(string, string, bool) bool
	AppendService(string, string, bool) bool
	DeleteService(string) bool
}

type Service struct {
	repo        repository
	filename    string
	recordTypes []string
}

func New(repo repository, filename string, recordTypes []string) *Service {
	return &Service{
		repo:        repo,
		filename:    filename,
		recordTypes: recordTypes,
	}
}

func (s *Service) readFile() (domain.Storage, error) {
	bytes, err := os.ReadFile(s.filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err.Error())
	}

	var data domain.Storage

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal file: %s", err.Error())
	}

	return data, nil
}

func (s *Service) WriteStorageFromFile() error {
	data, err := s.readFile()
	if err != nil {
		return fmt.Errorf("failed to write storage from file: %s", err.Error())
	}

	s.repo.SetStorage(data)

	return nil
}

func (s *Service) UpdateFile() error {
	file, err := os.OpenFile(s.filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open or create file: %s", err.Error())
	}
	defer file.Close()

	storage, err := json.Marshal(s.repo.GetAll())
	if err != nil {
		return fmt.Errorf("failed to marshal storage: %s", err.Error())
	}

	_, err = file.Write(storage)
	if err != nil {
		return fmt.Errorf("failed to write storage in file: %s", err.Error())
	}

	return nil
}

func (s *Service) GetAll() domain.Storage {
	return s.repo.GetAll()
}

func (s *Service) GetByType(recordType string) (domain.Storage, error) {
	if !contains(s.recordTypes, recordType) {
		return domain.Storage{}, fmt.Errorf("undefined record type")
	}

	storage := s.repo.GetAll()
	storageWithType := make(domain.Storage)

	for name, value := range storage {
		if value.Type != recordType {
			continue
		}

		storageWithType[name] = value
	}

	return storageWithType, nil
}

func (s *Service) AppendService(serviceName, serviceType string, favorite bool) error {
	validServiceName, err := validationServiceName(serviceName)
	if err != nil {
		return fmt.Errorf("validation name error: %s", err.Error())
	}

	validServiceType, err := validationServiceType(serviceType, s.recordTypes)
	if err != nil {
		return fmt.Errorf("validation elem error: %s", err.Error())
	}

	ok := s.repo.AppendService(validServiceName, validServiceType, favorite)
	if !ok {
		log.Print("failed to update file: element already exists")

		return fmt.Errorf("element already exists")
	}

	err = s.UpdateFile()
	if err != nil {
		log.Printf("failed to update file: %s", err.Error())

		return err
	}

	return nil
}

func (s *Service) UpdateService(serviceName, serviceType string, favorite bool) error {
	validServiceName, err := validationServiceName(serviceName)
	if err != nil {
		return fmt.Errorf("validation name error: %s", err.Error())
	}

	validServiceType, err := validationServiceType(serviceType, s.recordTypes)
	if err != nil {
		return fmt.Errorf("validation elem error: %s", err.Error())
	}

	ok := s.repo.UpdateService(validServiceName, validServiceType, favorite)
	if !ok {
		log.Print("failed to update file: element not found")

		return fmt.Errorf("element not found")
	}

	err = s.UpdateFile()
	if err != nil {
		log.Printf("failed to update file: %s", err.Error())
	}

	return nil
}

func (s *Service) DeleteService(serviceName, login string) error {
	validServiceName, err := validationServiceName(serviceName)
	if err != nil {
		return fmt.Errorf("validation name error: %s", err.Error())
	}

	validLogin, err := validationLogin(login)
	if err != nil {
		return fmt.Errorf("validation elem error: %s", err.Error())
	}

	ok := s.repo.DeleteLogin(validServiceName, validLogin)
	if !ok {
		log.Print("failed to update file: element not found")

		return fmt.Errorf("element not found")
	}

	err = s.UpdateFile()
	if err != nil {
		log.Printf("failed to update file: %s", err.Error())
	}

	return nil
}

func (s *Service) AppendLogin(serviceName, login string, elem domain.Element) error {
	validServiceName, err := validationServiceName(serviceName)
	if err != nil {
		return fmt.Errorf("validation name error: %s", err.Error())
	}

	validElem, err := validationElem(elem)
	if err != nil {
		return fmt.Errorf("validation elem error: %s", err.Error())
	}

	validLogin, err := validationLogin(login)
	if err != nil {
		return fmt.Errorf("validation elem error: %s", err.Error())
	}

	ok := s.repo.AppendLogin(validServiceName, validLogin, validElem)
	if !ok {
		log.Print("failed to update file: element already exists")

		return fmt.Errorf("element already exists")
	}

	err = s.UpdateFile()
	if err != nil {
		log.Printf("failed to update file: %s", err.Error())

		return err
	}

	return nil
}

func (s *Service) UpdateLogin(serviceName, login string, elem domain.Element) error {
	validServiceName, err := validationServiceName(serviceName)
	if err != nil {
		return fmt.Errorf("validation name error: %s", err.Error())
	}

	validElem, err := validationElem(elem)
	if err != nil {
		return fmt.Errorf("validation elem error: %s", err.Error())
	}

	validLogin, err := validationLogin(login)
	if err != nil {
		return fmt.Errorf("validation elem error: %s", err.Error())
	}

	ok := s.repo.UpdateLogin(validServiceName, validLogin, validElem)
	if !ok {
		log.Print("failed to update file: element not found")

		return fmt.Errorf("element not found")
	}

	err = s.UpdateFile()
	if err != nil {
		log.Printf("failed to update file: %s", err.Error())
	}

	return nil
}

func (s *Service) DeleteLogin(serviceName, login string) error {
	validServiceName, err := validationServiceName(serviceName)
	if err != nil {
		return fmt.Errorf("validation name error: %s", err.Error())
	}

	validLogin, err := validationLogin(login)
	if err != nil {
		return fmt.Errorf("validation elem error: %s", err.Error())
	}

	ok := s.repo.DeleteLogin(validServiceName, validLogin)
	if !ok {
		log.Print("failed to update file: element not found")

		return fmt.Errorf("element not found")
	}

	err = s.UpdateFile()
	if err != nil {
		log.Printf("failed to update file: %s", err.Error())
	}

	return nil
}

func validationServiceName(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("the name field cannot be empty")
	}

	lenName := len([]rune(name))

	if lenName > 100 {
		return "", fmt.Errorf("the name is too long")
	}

	if lenName < 4 {
		return "", fmt.Errorf("the name is too short")
	}

	name = strings.ToLower(name)

	name = strings.TrimSpace(name)

	return name, nil
}

func validationServiceType(serviceType string, recordTypes []string) (string, error) {
	if !contains(recordTypes, serviceType) {
		return "", fmt.Errorf("undefined record type")
	}

	return serviceType, nil
}

func validationElem(elem domain.Element) (domain.Element, error) {

	return elem, nil
}

func validationLogin(login string) (string, error) {
	if login == "" {
		return "", fmt.Errorf("login cannot be empty")
	}

	return login, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
