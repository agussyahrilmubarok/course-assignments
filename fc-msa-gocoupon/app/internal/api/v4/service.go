package v1

type IService interface {
}

type service struct {
}

func NewService() IService {
	return &service{}
}
