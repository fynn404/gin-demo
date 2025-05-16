export interface User {
    user_id: string;
    username: string;
    nickname?: string;
    email?: string;
    avatar_url?: string;
    created_at: string;
    updated_at: string;
}

export interface LoginRequest {
    username: string;
    password: string;
}

export interface LoginResponse {
    user_id: string;
    username: string;
    access_token: string;
}

export interface RegisterRequest {
    username: string;
    password: string;
    role: string;
    nickname?: string;
    email?: string;
}

export interface RegisterResponse {
    user_id: string;
    username: string;
    nickname?: string;
    email?: string;
}

export interface UpdateProfileRequest {
    nickname?: string;
    email?: string;
    password?: string;
    avatar_url?: string;
}

export interface UpdateProfileResponse {
    nickname: string;
    email: string;
    avatar_url: string;
    has_change_password: boolean;
} 