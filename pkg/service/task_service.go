package service

import (
    "context"
    "errors"

    "github.com/octaview/kanban-backend/pkg/repository"
)

type TaskService interface {
    MoveTask(ctx context.Context, taskID, targetColumnID int64) error
}

type taskService struct {
    taskRepo repository.TaskRepository
}

func NewTaskService(tr repository.TaskRepository) TaskService {
    return &taskService{taskRepo: tr}
}

func (s *taskService) MoveTask(ctx context.Context, taskID, targetColumnID int64) error {
    task, err := s.taskRepo.GetByID(ctx, taskID)
    if err != nil || task == nil {
        return errors.New("task not found")
    }
    task.ColumnID = targetColumnID
    return s.taskRepo.Update(ctx, task)
}
