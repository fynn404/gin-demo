import axios from 'axios';
import { User, UpdateProfileRequest, UpdateProfileResponse } from '../types/user';

const API_BASE_URL = '/api/v1';

export const getUserProfile = async (): Promise<User> => {
    const response = await axios.get<{ data: User }>(`${API_BASE_URL}/users/profile`);
    return response.data.data;
};

export const updateUserProfile = async (data: UpdateProfileRequest): Promise<UpdateProfileResponse> => {
    const response = await axios.put<{ data: UpdateProfileResponse }>(`${API_BASE_URL}/users/profile`, data);
    return response.data.data;
}; 