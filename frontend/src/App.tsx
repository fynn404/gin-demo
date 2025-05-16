import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { Layout } from 'antd';
import { UserProvider } from './contexts/UserContext';
import Navbar from './components/Navbar';
import Login from './components/Login';
import Register from './components/Register';
import UserProfile from './components/UserProfile';
import TodoList from './components/TodoList';
import PrivateRoute from './components/PrivateRoute';
import { useUser } from './contexts/UserContext';

const { Content } = Layout;

const AppContent: React.FC = () => {
    const { user } = useUser();

    return (
        <Layout style={{ minHeight: '100vh' }}>
            <Navbar username={user?.username} avatarUrl={user?.avatar_url} />
            <Content>
                <Routes>
                    <Route path="/login" element={<Login />} />
                    <Route path="/register" element={<Register />} />
                    <Route
                        path="/todos"
                        element={
                            <PrivateRoute>
                                <TodoList />
                            </PrivateRoute>
                        }
                    />
                    <Route
                        path="/profile"
                        element={
                            <PrivateRoute>
                                <UserProfile />
                            </PrivateRoute>
                        }
                    />
                    <Route path="/" element={<Navigate to="/todos" replace />} />
                </Routes>
            </Content>
        </Layout>
    );
};

const App: React.FC = () => {
    return (
        <Router>
            <UserProvider>
                <AppContent />
            </UserProvider>
        </Router>
    );
};

export default App; 