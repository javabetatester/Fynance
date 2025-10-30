package dashboard

type Service struct {
	Repository Repository
}

func (s *Service) Create(dashboard *Dashboard) error {
	return s.Repository.Create(dashboard)
}

func (s *Service) Update(dashboard *Dashboard) error {
	return s.Repository.Update(dashboard)
}

func (s *Service) GetByID(id string) (*Dashboard, error) {
	return s.Repository.GetByID(id)
}
