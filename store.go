package main

import (
	"encoding/binary"
	"encoding/json"
	"fideliy/dins"
	"log"
	"time"

	"github.com/prologic/bitcask"
)

//Store it's a wraper for Bitcask store
type Store struct {
	db *bitcask.Bitcask
}

//NewStore returns new Store instance
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

//Put user into store by telegram chatID
func (s *Store) Put(chatID int64, user dins.User) error {
	err := s.db.Put(byteOf(chatID), user.GetBytes())
	return err
}

//Get user from store by telegram chatID
func (s *Store) Get(chatID int64) (dins.User, error) {
	parsed := dins.User{Subs: map[string]time.Time{}}
	bytes, err := s.db.Get(byteOf(chatID))

	if err != nil {
		log.Println("Read User failed: ", err)
		return dins.User{Subs: map[string]time.Time{}}, err
	}
	parseErr := json.Unmarshal(bytes, &parsed)
	if parseErr != nil {
		log.Println("Failed Parse user: ", parseErr)
		return dins.User{Subs: map[string]time.Time{}}, parseErr
	}

	return parsed, nil
}

//Has return true if Store has a user
func (s *Store) Has(chatID int64) bool {
	return s.db.Has(byteOf(chatID))
}

//Keys return channel with all telegramm chatIds
func (s *Store) Keys() chan []byte {
	return s.db.Keys()
}

//Close gracefully close store
func (s *Store) Close() {
	if err := s.db.Close(); err != nil {
		log.Println("Close Store failed: ", err)
	}
}
