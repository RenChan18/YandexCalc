package orchestration

type TaskQueue struct {
    tasks chan Task
}

func NewTaskQueue(bufferSize int) *TaskQueue {
    return &TaskQueue{
        tasks: make(chan Task, bufferSize),
    }
}

// Enqueue добавляет задачу в очередь.
func (q *TaskQueue) Enqueue(task Task) {
    q.tasks <- task
}

// Dequeue возвращает задачу из очереди, если она есть.
func (q *TaskQueue) Dequeue() (Task, bool) {
    select {
    case task := <-q.tasks:
        return task, true
    default:
        return Task{}, false
    }
}
