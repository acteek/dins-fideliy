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
	s.db.Sync()
	return err
}

//Get user from store by telegram chatID
func (s *Store) Get(chatID int64) (dins.User, error) {
	parsed := dins.User{}
	bytes, err := s.db.Get(byteOf(chatID))

	if err != nil {
		log.Println("Read User failed: ", err)
		return dins.User{}, err
	}
	parseErr := json.Unmarshal(bytes, &parsed)
	if parseErr != nil {
		log.Println("Failed Parse user: ", parseErr)
		return dins.User{}, parseErr
	}
	if parsed.Subs == nil {
		parsed.Subs = map[string]time.Time{}
	}

	return parsed, nil
}

//Has return true if Store has a user
func (s *Store) Has(chatID int64) bool {
	return s.db.Has(byteOf(chatID))
}

//ChatIDs return channel with all telegramm chatIds
func (s *Store) ChatIDs() chan int64 {
	byteCh := s.db.Keys()
	chatCh := make(chan int64)

	go func() {
		for key := range byteCh {
			if len(key) == 0 {
				break
			}
			chatCh <- int64(binary.LittleEndian.Uint64(key))
		}
		defer close(chatCh)
	}()

	return chatCh
}

//Close gracefully close store
func (s *Store) Close() {
	if err := s.db.Sync(); err != nil {
		log.Println("Sync Store failed: ", err)
	}

	if err := s.db.Close(); err != nil {
		log.Println("Close Store failed: ", err)
	}
}
