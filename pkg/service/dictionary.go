package service

import (
	"encoding/json"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/db"
	"github.com/hidevopsio/hicicd/pkg/admin"
)

const Dictionary  = "dictionary"

type DictionaryService struct {
	Repository db.KVRepository `inject:"repository,dataSourceType=bolt,namespace=dictionary"`
}

func (ds *DictionaryService) Add(dictionary *admin.Dictionary) error {
	d, err := json.Marshal(dictionary)
	if err == nil {
		ds.Repository.Put([]byte(Dictionary), []byte(dictionary.Id), d)
	}
	return nil
}

func (ds *DictionaryService) Get(id string) (*admin.Dictionary, error) {
	d, err := ds.Repository.Get([]byte(Dictionary), []byte(id))
	if err != nil {
		return nil, err
	}
	var dictionary = &admin.Dictionary{}
	err = json.Unmarshal(d, dictionary)
	return dictionary, err
}

func (ds *DictionaryService) Delete(id string) error {
	err := ds.Repository.Delete([]byte(Dictionary), []byte(id))
	log.Debug("Delete Dictionary Service:", err)
	return err
}
