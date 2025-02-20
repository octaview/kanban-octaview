package repository

import (
    "context"
    "github.com/octaview/kanban-backend/internal/board" // убедитесь, что путь правильный
)

// TaskRepository описывает операции для работы с задачами.
type TaskRepository interface {
    // GetByID получает задачу по её идентификатору.
    GetByID(ctx context.Context, taskID int64) (*board.Task, error)
    // Update обновляет данные задачи.
    Update(ctx context.Context, task *board.Task) error
    // Можно добавить дополнительные методы: Create, Delete, List и т.д.
}

type BoardRepository interface {
    // Create создаёт новую доску и возвращает её идентификатор.
    Create(ctx context.Context, b *board.Board) (int64, error)
    // GetByID получает доску по её идентификатору.
    GetByID(ctx context.Context, id int64) (*board.Board, error)
    // Добавьте дополнительные методы по необходимости.
}