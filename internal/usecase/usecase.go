package usecase

import (
	"github.com/sirupsen/logrus"
	"kvStorage/internal/repository/storage"
	"kvStorage/types"
)

type UseCase struct {
	storage *storage.Storage
	log     *logrus.Logger
}

func New(storage *storage.Storage, log *logrus.Logger) *UseCase {
	return &UseCase{
		storage: storage,
		log:     log,
	}
}

func (s *UseCase) GetValue(key string) (string, error) {
	val, err := s.storage.GetValue(key)
	if err != nil {
		s.log.WithError(err).Errorln("Can`t get value")
		return "", err
	}
	return val, nil
}

func (s *UseCase) PutValue(pair types.KeyValue) error {
	err := s.storage.PutKeyValues(pair)
	if err != nil {
		s.log.Error(err)
	}
	return err
}

func (s *UseCase) SetValue(pair types.KeyValue) error {
	err := s.storage.UpdateValue(pair)
	if err != nil {
		s.log.Error(err)
	}
	return err
}

func (s *UseCase) DeleteValue(key string) error {
	err := s.storage.RemoveKeyValue(key)
	if err != nil {
		s.log.Error(err)
	}
	return err
}
