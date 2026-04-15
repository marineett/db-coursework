package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
)

type ILessonService interface {
	CreateLesson(lesson types.Lesson) (int64, error)
	GetLessons(contractID int64, from int64, size int64) ([]types.Lesson, error)
}

type LessonService struct {
	lessonRepository data_base.ILessonRepository
}

func CreateLessonService(lessonRepository data_base.ILessonRepository) ILessonService {
	return &LessonService{lessonRepository: lessonRepository}
}

func (s *LessonService) CreateLesson(lesson types.Lesson) (int64, error) {
	return s.lessonRepository.InsertLesson(lesson)
}

func (s *LessonService) GetLessons(contractID int64, from int64, size int64) ([]types.Lesson, error) {
	return s.lessonRepository.GetLessons(contractID, from, size)
}
