package dashboard

type Repository interface {
	Create(dashboard *Dashboard) error
	Update(dashboard *Dashboard) error
	GetByID(id string) (*Dashboard, error)
}
