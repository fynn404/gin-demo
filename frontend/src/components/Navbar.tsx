import React from 'react';
import { Layout, Menu, Avatar, Dropdown } from 'antd';
import { UserOutlined, LogoutOutlined, UnorderedListOutlined, ProfileOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { logout } from '../services/auth';

const { Header } = Layout;

interface NavbarProps {
    username?: string;
    avatarUrl?: string;
}

const Navbar: React.FC<NavbarProps> = ({ username, avatarUrl }) => {
    const navigate = useNavigate();

    const handleLogout = async () => {
        try {
            await logout();
            navigate('/login');
        } catch (error) {
            console.error('Logout failed:', error);
        }
    };

    const userMenu = (
        <Menu>
            <Menu.Item key="profile" onClick={() => navigate('/profile')} icon={<ProfileOutlined />}>
                个人信息
            </Menu.Item>
            <Menu.Item key="logout" onClick={handleLogout} icon={<LogoutOutlined />}>
                退出登录
            </Menu.Item>
        </Menu>
    );

    return (
        <Header style={{ background: '#fff', padding: '0 24px', display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
            <div style={{ display: 'flex', alignItems: 'center' }}>
                <h1 style={{ margin: 0, marginRight: 48 }}>Todo App</h1>
                <Menu mode="horizontal" defaultSelectedKeys={['todos']}>
                    <Menu.Item key="todos" icon={<UnorderedListOutlined />} onClick={() => navigate('/todos')}>
                        待办事项
                    </Menu.Item>
                </Menu>
            </div>

            {username && (
                <Dropdown overlay={userMenu} placement="bottomRight">
                    <div style={{ cursor: 'pointer', display: 'flex', alignItems: 'center' }}>
                        <Avatar src={avatarUrl} icon={<UserOutlined />} style={{ marginRight: 8 }} />
                        <span>{username}</span>
                    </div>
                </Dropdown>
            )}
        </Header>
    );
};

export default Navbar; 