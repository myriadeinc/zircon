package cache

import (
	"github.com/myriadeinc/zircon/internal/models"
)

type DummyService struct{}

func NewDummyClient() CacheService {
	return &DummyService{}
}

func (d *DummyService) SaveNewTemplate(models.StrictTemplate) error {
	return nil
}

func (d *DummyService) FetchTemplate() (*models.StrictTemplate, error) {

	return nil, nil
}
