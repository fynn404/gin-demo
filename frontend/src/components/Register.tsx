import React from 'react';
import { Form, Input, Button, Card, message } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined, UserAddOutlined } from '@ant-design/icons';
import { register } from '../services/auth';
import { RegisterRequest } from '../types/user';
import { useNavigate } from 'react-router-dom';

const Register: React.FC = () => {
    const navigate = useNavigate();

    const onFinish = async (values: RegisterRequest) => {
        try {
            // 默认注册为普通用户
            const registerData = { ...values, role: 'user' };
            await register(registerData);
            message.success('注册成功，请登录');
            navigate('/login');
        } catch (error) {
            message.error('注册失败，请稍后重试');
        }
    };

    return (
        <div style={{ maxWidth: 400, margin: '100px auto', padding: '0 20px' }}>
            <Card title="注册" bordered={false}>
                <Form
                    name="register"
                    onFinish={onFinish}
                    autoComplete="off"
                >
                    <Form.Item
                        name="username"
                        rules={[
                            { required: true, message: '请输入用户名' },
                            { min: 3, message: '用户名至少3个字符' }
                        ]}
                    >
                        <Input prefix={<UserOutlined />} placeholder="用户名" />
                    </Form.Item>

                    <Form.Item
                        name="password"
                        rules={[
                            { required: true, message: '请输入密码' },
                            { min: 6, message: '密码至少6个字符' }
                        ]}
                    >
                        <Input.Password prefix={<LockOutlined />} placeholder="密码" />
                    </Form.Item>

                    <Form.Item
                        name="confirm"
                        dependencies={['password']}
                        rules={[
                            { required: true, message: '请确认密码' },
                            ({ getFieldValue }) => ({
                                validator(_, value) {
                                    if (!value || getFieldValue('password') === value) {
                                        return Promise.resolve();
                                    }
                                    return Promise.reject(new Error('两次输入的密码不一致'));
                                },
                            }),
                        ]}
                    >
                        <Input.Password prefix={<LockOutlined />} placeholder="确认密码" />
                    </Form.Item>

                    <Form.Item
                        name="nickname"
                        rules={[{ required: false }]}
                    >
                        <Input prefix={<UserAddOutlined />} placeholder="昵称（选填）" />
                    </Form.Item>

                    <Form.Item
                        name="email"
                        rules={[
                            { type: 'email', message: '请输入有效的邮箱地址' },
                            { required: false }
                        ]}
                    >
                        <Input prefix={<MailOutlined />} placeholder="邮箱（选填）" />
                    </Form.Item>

                    <Form.Item>
                        <Button type="primary" htmlType="submit" block>
                            注册
                        </Button>
                    </Form.Item>

                    <Form.Item>
                        <Button type="link" block onClick={() => navigate('/login')}>
                            已有账号？立即登录
                        </Button>
                    </Form.Item>
                </Form>
            </Card>
        </div>
    );
};

export default Register; 