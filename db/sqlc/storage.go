package db

type Store interface {
	Querier
}

type Storage struct {
	*Queries
}

func NewStorage(db DBTX) Store {
	return &Storage{
		Queries: New(db),
	}
}
