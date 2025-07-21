package storage

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/tarantool/go-tarantool/v2"
	"github.com/tarantool/go-tarantool/v2/pool"
	"kvStorage/types"
	"time"
)

type Storage struct {
	conn *pool.ConnectionPool
	log  *logrus.Logger
}

func New(dbUrl, user, psw string, logger *logrus.Logger) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	dialer := &tarantool.NetDialer{
		Address:  dbUrl,
		User:     user,
		Password: psw,
	}

	opts := tarantool.Opts{
		Timeout: time.Second,
	}

	inst := pool.Instance{
		Name:   "",
		Dialer: dialer,
		Opts:   opts,
	}

	connPool, err := pool.Connect(ctx, []pool.Instance{inst})
	if err != nil {
		logrus.Errorf("Connection refused %s", err)
		return nil, err
	}

	return &Storage{
		conn: connPool,
		log:  logger,
	}, nil
}

func (s *Storage) GetValue(key string) (value string, err error) {
	req := tarantool.NewSelectRequest("kv_storage").
		Index("primary").
		Key([]interface{}{key}).
		Limit(1).
		Iterator(tarantool.IterEq)

	data, err := s.conn.Do(req, pool.ANY).Get()
	if err != nil {
		s.log.WithError(err).Errorln("Get value error")
		return "", err
	}

	if len(data) == 0 {
		return "", types.ErrKeyNotFound
	}

	tuple, ok := data[0].([]interface{})
	if !ok || len(tuple) < 2 {
		return "", errors.New("Invalid tuple format")
	}

	val, ok := tuple[1].(string)
	if !ok {
		return "", errors.New("Value is not a string")
	}

	return val, nil
}

func (s *Storage) PutKeyValues(pair types.KeyValue) error {
	req := tarantool.NewInsertRequest("kv_storage").Tuple([]interface{}{pair.Key, pair.Value})
	data, err := s.conn.Do(req, pool.RW).Get()
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("empty response from Tarantool")
	}
	return nil
}

func (s *Storage) UpdateValue(pair types.KeyValue) error {
	req := tarantool.NewUpdateRequest("kv_storage").
		Index("primary").
		Key([]interface{}{pair.Key}).
		Operations(
			tarantool.NewOperations().Assign(1, pair.Value),
		)

	data, err := s.conn.Do(req, pool.RW).Get()
	if err != nil {
		s.log.WithError(err).Errorln("Update operation failed")
	}

	if len(data) == 0 {
		return types.ErrKeyNotFound
	}

	return err
}

func (s *Storage) RemoveKeyValue(key string) (err error) {
	req := tarantool.NewDeleteRequest("kv_storage").Key([]interface{}{key})
	data, err := s.conn.Do(req, pool.RW).Get()
	if err != nil {
		s.log.Error(err)
	}
	if len(data) == 0 {
		s.log.WithError(err).Errorln("Key does not exist")
		return types.ErrKeyNotFound
	}
	return err
}
