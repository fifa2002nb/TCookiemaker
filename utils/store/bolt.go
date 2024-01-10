package store

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"io"
	"sync"
	"time"
)

var DefaultStore = &Store{}

func Load(storeFile string) error {
	return DefaultStore.Init(storeFile)
}

func Run() error {
	return DefaultStore.Run()
}

type Storetask struct {
	Key  string
	Tag  string
	Data interface{}
}

type Store struct {
	mutex         sync.Mutex
	mutexHolder   string
	mutexAquired  time.Time
	mutexWaitTime int64
	keytags       map[string]time.Time
	db            *bolt.DB
	Taskqueue     chan Storetask
}

// async saver
func PushTaskQueue(t Storetask) {
	select {
	case DefaultStore.Taskqueue <- t:
	default:
		log.Fatal("something wrong with taskqueue.")
	}
}

// async saver
func (s *Store) Run() error {
	for task := range s.Taskqueue {
		s.Save(task.Key, task.Tag, task.Data)
	}
	return nil
}

func (s *Store) Init(storeFile string) error {
	var err error
	s.keytags = make(map[string]time.Time)
	s.Taskqueue = make(chan Storetask, 100)
	if storeFile != "" {
		s.db, err = bolt.Open(storeFile, 0600, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) Lock(method string) {
	start := time.Now()
	s.mutex.Lock()
	s.mutexAquired = time.Now()
	s.mutexHolder = method
	s.mutexWaitTime = int64(s.mutexAquired.Sub(start) / time.Millisecond) // remember this so we don't have to call put until we leave the critical section.
}

func (s *Store) Unlock() {
	holder := s.mutexHolder
	start := s.mutexAquired
	waitTime := s.mutexWaitTime
	s.mutex.Unlock()
	log.Debugf("holder:%v wait:%v hold:%v", holder, waitTime, int64(time.Since(start)/time.Millisecond))
}

type counterWriter struct {
	written int
	w       io.Writer
}

func (c *counterWriter) Write(p []byte) (n int, err error) {
	n, err = c.w.Write(p)
	c.written += n
	return n, err
}

// sync saver
func SaveSync(key, tag string, data interface{}) error {
	DefaultStore.Save(key, tag, data)
	return nil
}

func (s *Store) Save(key, tag string, data interface{}) {
	if s.db == nil {
		return
	}
	s.Lock("Save")
	s.keytags[key+tag] = time.Now()
	tostore := make(map[string][]byte)
	f := new(bytes.Buffer)
	gz := gzip.NewWriter(f)
	cw := &counterWriter{w: gz}
	enc := gob.NewEncoder(cw)
	if err := enc.Encode(data); err != nil {
		log.Error("error saving %s: %v", key+tag, err)
		s.Unlock()
		return
	}
	if err := gz.Flush(); err != nil {
		log.Error("gzip flush error saving %s: %v", key+tag, err)
	}
	if err := gz.Close(); err != nil {
		log.Error("gzip close error saving %s: %v", key+tag, err)
	}
	tostore[key+tag] = f.Bytes()
	s.Unlock()
	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(key + tag))
		if err != nil {
			return err
		}
		for name, data := range tostore {
			if err := b.Put([]byte(name), data); err != nil {
				log.Error(err.Error())
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Error("save db update error: %v", err)
		return
	}
	log.Debug("save to db complete")
}

// sync restorer
func Restore(key, tag string, ret interface{}) error {
	return DefaultStore.RestoreState(key, tag, ret)
}

// RestoreState restores records from the file on disk.
func (s *Store) RestoreState(key, tag string, ret interface{}) error {
	log.Debug("RestoreState")
	start := time.Now()
	s.Lock("RestoreState")
	defer s.Unlock()
	decode := func(name string, dst interface{}) error {
		var data []byte
		err := s.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(key + tag))
			if b == nil {
				return fmt.Errorf("unknown bucket: %v", key+tag)
			}
			data = b.Get([]byte(name))
			return nil
		})
		if err != nil {
			return err
		}
		gr, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return err
		}
		defer gr.Close()
		return gob.NewDecoder(gr).Decode(dst)
	}
	if err := decode(key+tag, ret); err != nil {
		return err
	} else {
		log.Debug("RestoreState done in", time.Since(start))
		delete(s.keytags, key+tag)
		return nil
	}
}
