package store

import "uroborus/model"

// UserStore -
type UserStore struct {
	db *DB
}

// NewUserStore -
func NewUserStore(db *DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) Save(user *model.User) error {
	return s.db.Create(user).Error
}

func (s UserStore) Get(user *model.User) error {
	return s.db.First(user, user).Error
}
