package v1

type IRepository interface {
}

type repository struct {
}

func NewRepository() IRepository {
	return &repository{}
}
