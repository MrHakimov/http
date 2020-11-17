package banners

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"sync"
)

// STORAGE is place where to store
const STORAGE = "./web/banners/"

var starID int64 = 0

// Banner type
type Banner struct {
	ID      int64
	Title   string
	Content string
	Button  string
	Link    string
	Image   string
}

// Service type
type Service struct {
	mu    sync.RWMutex
	items []*Banner
}

// NewService construct
func NewService() *Service {
	return &Service{items: make([]*Banner, 0)}
}

// All simple implementation
func (s *Service) All(ctx context.Context) ([]*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.items, nil
}

// ByID function to look for banner by id
func (s *Service) ByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, banner := range s.items {
		if banner.ID == id {
			return banner, nil
		}
	}
	return nil, errors.New("banner by id not found")
}

// Save banner
func (s *Service) Save(ctx context.Context, item *Banner, image multipart.File) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if item.ID == 0 {
		starID++
		item.ID = starID
		if item.Image != "" {
			item.Image = fmt.Sprint(item.ID) + "." + item.Image
			err := uploadFile(image, item)
			if err != nil {
				return nil, err
			}
		}
		s.items = append(s.items, item)
		return item, nil
	}
	for i, banner := range s.items {
		if banner.ID == item.ID {
			if item.Image != "" {
				item.Image = fmt.Sprint(item.ID) + "." + item.Image
				err := uploadFile(image, item)
				if err != nil {
					return nil, err
				}
			} else {
				item.Image = s.items[i].Image
			}
			s.items[i] = item
			return item, nil
		}
	}

	return nil, errors.New("banner save error")
}

func (s *Service) RemoveByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i, banner := range s.items {
		if banner.ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return banner, nil
		}
	}
	return nil, errors.New("banner remove by id not found")
}
func uploadFile(file multipart.File, banner *Banner) error {
	var data, err = ioutil.ReadAll(file)
	if err != nil {
		return errors.New("error read file")
	}

	err = ioutil.WriteFile(STORAGE+banner.Image, data, 0666)
	if err != nil {
		return errors.New("error to write file")
	}

	return nil
}
