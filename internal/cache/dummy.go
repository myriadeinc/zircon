package cache

type DummyService struct{}

func NewDummyClient() CacheService {
	return &DummyService{}
}

func (d *DummyService) SaveNewTemplate(template map[string]interface{}) error {
	return nil
}

func (d *DummyService) FetchTemplate() (*StrictTemplate, error) {

	return nil, nil
}
