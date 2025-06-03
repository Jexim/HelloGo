package usecase

import "github.com/Jexim/HelloGo/internal/svc/hello/model"

type Usecase struct {
	model.Datastore
}

func New(ds model.Datastore) *Usecase {
	return &Usecase{
		Datastore: ds,
	}
}
