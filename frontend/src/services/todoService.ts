import axios from 'axios';
import {
    TodoItem,
    ListTodosResponse,
    CreateTodoRequest,
    UpdateTodoRequest,
    ChangeTodoStatusRequest
} from '../types/todo';

const API_BASE_URL = '/api/v1';

export const todoService = {
    // 获取待办事项列表
    async listTodos(page: number = 1, size: number = 10): Promise<ListTodosResponse> {
        const response = await axios.get(`${API_BASE_URL}/todos`, {
            params: { page, size }
        });
        return response.data.data;
    },

    // 创建待办事项
    async createTodo(todo: CreateTodoRequest): Promise<TodoItem> {
        const response = await axios.post(`${API_BASE_URL}/todos`, todo);
        return response.data.data;
    },

    // 获取待办事项详情
    async getTodoDetail(id: number): Promise<TodoItem> {
        const response = await axios.get(`${API_BASE_URL}/todos/${id}`);
        return response.data.data;
    },

    // 更新待办事项
    async updateTodo(id: number, todo: UpdateTodoRequest): Promise<TodoItem> {
        const response = await axios.put(`${API_BASE_URL}/todos/${id}`, todo);
        return response.data.data;
    },

    // 删除待办事项
    async deleteTodo(id: number): Promise<void> {
        await axios.delete(`${API_BASE_URL}/todos/${id}`);
    },

    // 更改待办事项状态
    async changeTodoStatus(id: number, status: ChangeTodoStatusRequest): Promise<TodoItem> {
        const response = await axios.patch(`${API_BASE_URL}/todos/${id}/status`, status);
        return response.data.data;
    }
};

