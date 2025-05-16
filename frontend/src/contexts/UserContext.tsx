import React, { createContext, useContext, useState, useEffect } from 'react';
import { User } from '../types/user';
import { getUserProfile } from '../services/user';

interface UserContextType {
    user: User | null;
    setUser: (user: User | null) => void;
    loading: boolean;
    refreshUser: () => Promise<void>;
}

const UserContext = createContext<UserContextType>({
    user: null,
    setUser: () => { },
    loading: true,
    refreshUser: async () => { },
});

export const useUser = () => useContext(UserContext);

export const UserProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);

    const refreshUser = async () => {
        try {
            const token = localStorage.getItem('token');
            if (token) {
                const userData = await getUserProfile();
                setUser(userData);
            }
        } catch (error) {
            console.error('Failed to fetch user profile:', error);
            setUser(null);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        refreshUser();
    }, []);

    return (
        <UserContext.Provider value={{ user, setUser, loading, refreshUser }}>
            {children}
        </UserContext.Provider>
    );
};

export default UserContext; 