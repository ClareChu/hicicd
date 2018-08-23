package service

import (
	"github.com/hidevopsio/hicicd/pkg/entity"
		"github.com/hidevopsio/hiboot/pkg/log"
)

type SelectorService struct {
	repository BoltRepository
}

func (ss *SelectorService) Init(repository BoltRepository)  {
	log.Debug(ss)
	log.Debug(repository)
	log.Debug(ss.repository)
	ss.repository = repository
}

func (ss *SelectorService) Add(selectors *entity.Selector) error {
	err := ss.repository.Put(selectors)
	return err
}

func (ss *SelectorService) Get(id string) (*entity.Selector, error) {
	var selectors = entity.Selector{}
	err := ss.repository.Get(id, &selectors)
	if err != nil {
		return nil, err
	}

	return &selectors, err
}


func (ss *SelectorService) Delete(id string) error  {
	err := ss.repository.Delete([]byte(id))
	log.Debug("Delete Selector Service:", err)
	return err
}





