package main

import (
	"encoding/binary"
	"encoding/json"
	"fideliy/dins"
	"github.com/prologic/bitcask"
	"log"
)

type Store struct {
	db *bitcask.Bitcask
}

func NewStore(path string) *Store {
	db, err := bitcask.Open(path)
	if err != nil {
		log.Fatal("Read User failed: ", err)
	}
	return &Store{db}

}

func byteOf(key int64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(key))
	return bytes
}

func (s *Store) Put(chatId int64, user dins.User) error {
	err := s.db.Put(byteOf(chatId), user.GetBytes())
	return err
}

func (s *Store) Get(chatId int64) (dins.User, error) {
	parsed := dins.User{}
	bytes, err := s.db.Get(byteOf(chatId))
	if err != nil {
		log.Println("Read User failed: ", err)
		return dins.User{}, err
	}
	parseErr := json.Unmarshal(bytes, &parsed)
	if parseErr != nil {
		log.Println("Failed Parse user: ", parseErr)
		return dins.User{}, parseErr
	}

	return parsed, nil
}

func (s *Store) Has(chatId int64) bool {
	return s.db.Has(byteOf(chatId))
}

func (s *Store) Close() {
	if err := s.db.Close(); err != nil {
		log.Println("Close Store failed: ", err)
	}
}
