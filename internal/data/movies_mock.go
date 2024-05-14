package data

import (
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
)

type MockMovieStorage struct {
	DB         *pgxpool.Pool // Optional if you need structural compatibility
	InsertFunc func(movie *Movie) error
	GetFunc    func(id int64) (*Movie, error)
	GetAllFunc func(title string, genres []string, filters Filters) ([]*Movie, Metadata, error)
	UpdateFunc func(movie *Movie) error
	DeleteFunc func(id int64) error
}

func (m *MockMovieStorage) Insert(movie *Movie) error {
	if m.InsertFunc != nil {
		return m.InsertFunc(movie)
	}
	return nil // or simulate a realistic response
}

func (m *MockMovieStorage) Get(id int64) (*Movie, error) {
	if m.GetFunc != nil {
		return m.GetFunc(id)
	}
	return nil, errors.New("function not implemented")
}

func (m *MockMovieStorage) GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(title, genres, filters)
	}
	return nil, Metadata{}, errors.New("function not implemented")
}

func (m *MockMovieStorage) Update(movie *Movie) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(movie)
	}
	return errors.New("function not implemented")
}

func (m *MockMovieStorage) Delete(id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return errors.New("function not implemented")
}
