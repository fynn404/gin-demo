export interface TodoItem {
    id: number;
    title: string;
    description: string;
    priority: 'low' | 'medium' | 'high';
    dueDate: string;
    completed: boolean;
    createdAt: string;
    updatedAt: string;
}

export interface ListTodosResponse {
    todos: TodoItem[];
    total: number;
}

export interface CreateTodoRequest {
    title: string;
    description?: string;
    priority?: string;
    dueDate?: string;
}

export interface UpdateTodoRequest {
    title?: string;
    description?: string;
    priority?: string;
    dueDate?: string;
}

export interface ChangeTodoStatusRequest {
    completed: boolean;
} 