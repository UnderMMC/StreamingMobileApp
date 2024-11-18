package service

import "log"

type Repository interface {
	FrameSaving(frame []byte) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GettingStreamFrames(frame []byte) error {
	err := s.repo.FrameSaving(frame)
	if err != nil {
		log.Println("Frame saving error:", err)
		return err
	}
	return nil
}
