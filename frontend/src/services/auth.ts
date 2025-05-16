import axios from 'axios';
import { LoginRequest, LoginResponse, RegisterRequest, RegisterResponse } from '../types/user';

const API_BASE_URL = '/api/v1';

export const login = async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await axios.post<{ data: LoginResponse }>(`${API_BASE_URL}/auth/login`, data);
    const { access_token } = response.data.data;
    // 保存token到localStorage
    localStorage.setItem('token', access_token);
    // 设置axios默认header
    axios.defaults.headers.common['Authorization'] = `Bearer ${access_token}`;
    return response.data.data;
};

export const register = async (data: RegisterRequest): Promise<RegisterResponse> => {
    const response = await axios.post<{ data: RegisterResponse }>(`${API_BASE_URL}/auth/register`, data);
    return response.data.data;
};

export const logout = async (): Promise<void> => {
    await axios.post(`${API_BASE_URL}/auth/logout`);
    // 清除localStorage中的token
    localStorage.removeItem('token');
    // 清除axios默认header
    delete axios.defaults.headers.common['Authorization'];
};

// 初始化axios认证header
const token = localStorage.getItem('token');
if (token) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
} 