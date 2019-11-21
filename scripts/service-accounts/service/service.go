package service

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/ONSdigital/log.go/log"
)

var (
	serviceIDChars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

// Account represents data structure stored against a service
type Account struct {
	ID   string `json:"-"`
	Name string `json:"id,omitempty"`
}

// CreateDirectory represents the creation `services` directory in chosen path
func CreateDirectory(path string, replaceServiceDir bool) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Event(nil, "services dir does not exist, now creating")
		if err = os.MkdirAll(path, 0777); err != nil {
			return err
		}

		return nil
	}

	if replaceServiceDir {
		// Remove existing services directory before recreating
		log.Event(nil, "removing directory")
		os.RemoveAll(path + "/")
	}

	log.Event(nil, "services dir does not exist, now creating")
	if err := os.MkdirAll(path, 0777); err != nil {
		return err
	}

	return nil
}

// LoadExisting represents building a list of existing service accounts
func LoadExisting(servicePath string) (map[string]*Account, error) {
	results := make(map[string]*Account, 0)

	err := filepath.Walk(servicePath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".json" {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			var acc Account
			if err := json.Unmarshal(b, &acc); err != nil {
				return err
			}
			results[acc.Name] = &acc
		}
		return nil
	})

	return results, err
}

// CreateJSONAccount represents the creation of a service account (incl. token)
func CreateJSONAccount(r *rand.Rand, path, service string) (*Account, error) {
	serviceAcc := &Account{
		ID:   newRandomID(r, 64),
		Name: service,
	}

	logD := log.Data{
		"service": service,
	}

	jsonData, err := json.Marshal(serviceAcc)
	if err != nil {
		log.Event(nil, "error attempting to marshal service account to json", log.Error(err), logD)
		return nil, err
	}

	p := filepath.Join(path, serviceAcc.ID+".json")
	jsonFile, err := os.Create(p)
	if err != nil {
		log.Event(nil, "error creating service account json", log.Error(err), logD)
		return nil, err
	}
	defer func() {
		err := jsonFile.Close()
		if err != nil {
			log.Event(nil, "error closing service account json", log.Error(err), logD)
		}
	}()

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		log.Event(nil, "error writing service account to file", log.Error(err), logD)
		return nil, err
	}
	return serviceAcc, nil
}

func newRandomID(r *rand.Rand, size int) string {
	id := ""
	for len(id) < size {
		id += string(serviceIDChars[r.Intn(len(serviceIDChars))])
	}

	return id
}
