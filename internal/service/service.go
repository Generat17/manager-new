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
	Get(string) (domain.Element, bool)
	GetAll() domain.Storage
	Update(string, domain.Element) bool
	Append(string, domain.Element) bool
	DeleteByName(string) bool
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

func (s *Service) DeleteByName(name string) error {
	validName, err := validationName(name)
	if err != nil {
		return fmt.Errorf("validation name error: %s", err.Error())
	}

	ok := s.repo.DeleteByName(validName)
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

func (s *Service) Append(name string, elem domain.Element) error {
	validName, err := validationName(name)
	if err != nil {
		return fmt.Errorf("validation name error: %s", err.Error())
	}

	validElem, err := validationElem(elem, s.recordTypes)
	if err != nil {
		return fmt.Errorf("validation elem error: %s", err.Error())
	}

	ok := s.repo.Append(validName, validElem)
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

func (s *Service) UpdateByName(name string, elem domain.Element) error {
	validName, err := validationName(name)
	if err != nil {
		return fmt.Errorf("validation name error: %s", err.Error())
	}

	validElem, err := validationElem(elem, s.recordTypes)
	if err != nil {
		return fmt.Errorf("validation elem error: %s", err.Error())
	}

	ok := s.repo.Update(validName, validElem)
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

func validationName(name string) (string, error) {
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

func validationElem(elem domain.Element, recordTypes []string) (domain.Element, error) {
	if !contains(recordTypes, elem.Type) {
		return domain.Element{}, fmt.Errorf("undefined record type")
	}

	return elem, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
