package service

import (
	"context"
	"fmt"

	"github.com/octaview/kanban-backend/pkg/repository"
)

// TaskService описывает бизнес‑логику для работы с задачами.
type TaskService interface {
	// MoveTask перемещает задачу с заданным taskID в колонку targetColumnID.
	MoveTask(ctx context.Context, taskID, targetColumnID int64) error
}

type taskService struct {
	taskRepo repository.TaskRepository
}

// NewTaskService создает новый экземпляр TaskService.
func NewTaskService(tr repository.TaskRepository) TaskService {
	return &taskService{taskRepo: tr}
}

// MoveTask реализует перемещение задачи.
// Теперь ошибки обрабатываются отдельно, что позволяет точнее понять, на каком этапе произошла ошибка.
func (s *taskService) MoveTask(ctx context.Context, taskID, targetColumnID int64) error {
	// Получаем задачу по ID
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to get task with id %d: %w", taskID, err)
	}
	if task == nil {
		return fmt.Errorf("task with id %d not found", taskID)
	}

	// Обновляем поле ColumnID
	task.ColumnID = targetColumnID

	// Сохраняем изменения
	if err := s.taskRepo.Update(ctx, task); err != nil {
		return fmt.Errorf("failed to update task with id %d: %w", taskID, err)
	}
	return nil
}
