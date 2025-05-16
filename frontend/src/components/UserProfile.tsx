import React, { useEffect, useState } from 'react';
import { Form, Input, Button, Card, message, Avatar, Upload } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined, UploadOutlined } from '@ant-design/icons';
import { getUserProfile, updateUserProfile } from '../services/user';
import { User, UpdateProfileRequest } from '../types/user';
import type { RcFile, UploadProps } from 'antd/es/upload';
import type { UploadFile } from 'antd/es/upload/interface';

const UserProfile: React.FC = () => {
    const [form] = Form.useForm();
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(false);
    const [avatarUrl, setAvatarUrl] = useState<string>();

    useEffect(() => {
        fetchUserProfile();
    }, []);

    const fetchUserProfile = async () => {
        try {
            const userData = await getUserProfile();
            setUser(userData);
            form.setFieldsValue({
                nickname: userData.nickname,
                email: userData.email,
            });
            setAvatarUrl(userData.avatar_url);
        } catch (error) {
            message.error('获取用户信息失败');
        }
    };

    const onFinish = async (values: UpdateProfileRequest) => {
        try {
            setLoading(true);
            const response = await updateUserProfile({
                ...values,
                avatar_url: avatarUrl,
            });
            message.success('个人信息更新成功');
            if (response.has_change_password) {
                message.success('密码修改成功');
            }
            await fetchUserProfile();
        } catch (error) {
            message.error('更新失败，请稍后重试');
        } finally {
            setLoading(false);
        }
    };

    const beforeUpload = (file: RcFile) => {
        const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png';
        if (!isJpgOrPng) {
            message.error('只能上传 JPG/PNG 格式的图片！');
        }
        const isLt2M = file.size / 1024 / 1024 < 2;
        if (!isLt2M) {
            message.error('图片大小不能超过 2MB！');
        }
        return isJpgOrPng && isLt2M;
    };

    const handleChange: UploadProps['onChange'] = (info: any) => {
        if (info.file.status === 'done') {
            // 这里假设后端返回的数据结构中包含图片URL
            setAvatarUrl(info.file.response.url);
            message.success('头像上传成功');
        } else if (info.file.status === 'error') {
            message.error('头像上传失败');
        }
    };

    return (
        <div style={{ maxWidth: 800, margin: '40px auto', padding: '0 20px' }}>
            <Card title="个人信息" bordered={false}>
                <div style={{ textAlign: 'center', marginBottom: 24 }}>
                    <Avatar
                        size={100}
                        src={avatarUrl}
                        icon={<UserOutlined />}
                    />
                    <div style={{ marginTop: 16 }}>
                        <Upload
                            name="avatar"
                            action="/api/v1/upload"
                            beforeUpload={beforeUpload}
                            onChange={handleChange}
                            showUploadList={false}
                        >
                            <Button icon={<UploadOutlined />}>更换头像</Button>
                        </Upload>
                    </div>
                </div>

                <Form
                    form={form}
                    name="userProfile"
                    onFinish={onFinish}
                    layout="vertical"
                >
                    <Form.Item label="用户名">
                        <Input value={user?.username} disabled />
                    </Form.Item>

                    <Form.Item
                        name="nickname"
                        label="昵称"
                    >
                        <Input prefix={<UserOutlined />} placeholder="设置昵称" />
                    </Form.Item>

                    <Form.Item
                        name="email"
                        label="邮箱"
                        rules={[
                            { type: 'email', message: '请输入有效的邮箱地址' }
                        ]}
                    >
                        <Input prefix={<MailOutlined />} placeholder="设置邮箱" />
                    </Form.Item>

                    <Form.Item
                        name="password"
                        label="新密码"
                        rules={[
                            { min: 6, message: '密码至少6个字符' }
                        ]}
                    >
                        <Input.Password prefix={<LockOutlined />} placeholder="不修改请留空" />
                    </Form.Item>

                    <Form.Item
                        name="confirm"
                        label="确认新密码"
                        dependencies={['password']}
                        rules={[
                            ({ getFieldValue }) => ({
                                validator(_, value) {
                                    if (!value || !getFieldValue('password') || getFieldValue('password') === value) {
                                        return Promise.resolve();
                                    }
                                    return Promise.reject(new Error('两次输入的密码不一致'));
                                },
                            }),
                        ]}
                    >
                        <Input.Password prefix={<LockOutlined />} placeholder="确认新密码" />
                    </Form.Item>

                    <Form.Item>
                        <Button type="primary" htmlType="submit" loading={loading} block>
                            保存修改
                        </Button>
                    </Form.Item>
                </Form>
            </Card>
        </div>
    );
};

export default UserProfile; 