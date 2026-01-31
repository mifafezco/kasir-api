package services

import (
	"kasir-api/model"
	"kasir-api/repositories"
)

type CategoriesService struct {
	repo *repositories.CategoriesRepository
}

func NewCategoriesService(repo *repositories.CategoriesRepository) *CategoriesService {
	return &CategoriesService{repo: repo}
}

func (s *CategoriesService) GetAll() ([]model.Categories, error) {
	return s.repo.GetAll()
}

func (s *CategoriesService) Create(data *model.Categories) error {
	return s.repo.Create(data)
}

func (s *CategoriesService) GetByID(id int) (*model.Categories, error) {
	return s.repo.GetByID(id)
}

func (s *CategoriesService) Update(categories *model.Categories) error {
	return s.repo.Update(categories)
}

func (s *CategoriesService) Delete(id int) error {
	return s.repo.Delete(id)
}