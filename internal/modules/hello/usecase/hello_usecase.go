package usecase

import "github.com/Jexim/HelloGo/internal/modules/hello/model"

type Usecase struct {
	model.Datastore
}

func New(ds model.Datastore) *Usecase {
	return &Usecase{
		Datastore: ds,
	}
}
