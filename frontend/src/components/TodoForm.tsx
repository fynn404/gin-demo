import React from 'react';
import { Form, Input, Select, DatePicker, Button, Space } from 'antd';
import { TodoItem } from '../types/todo';
import dayjs from 'dayjs';

interface TodoFormProps {
    initialValues?: TodoItem | null;
    onSubmit: (values: any) => void;
    onCancel: () => void;
}

const TodoForm: React.FC<TodoFormProps> = ({ initialValues, onSubmit, onCancel }) => {
    const [form] = Form.useForm();

    const handleSubmit = async () => {
        try {
            const values = await form.validateFields();
            // 转换日期格式
            if (values.dueDate) {
                values.dueDate = values.dueDate.toISOString();
            }
            onSubmit(values);
        } catch (error) {
            console.error('Validation failed:', error);
        }
    };

    // 设置初始值
    React.useEffect(() => {
        if (initialValues) {
            form.setFieldsValue({
                ...initialValues,
                dueDate: initialValues.dueDate ? dayjs(initialValues.dueDate) : undefined,
            });
        } else {
            form.resetFields();
        }
    }, [initialValues, form]);

    return (
        <Form
            form={form}
            layout="vertical"
            initialValues={{ priority: 'medium' }}
        >
            <Form.Item
                name="title"
                label="标题"
                rules={[{ required: true, message: '请输入标题' }]}
            >
                <Input placeholder="请输入待办事项标题" />
            </Form.Item>

            <Form.Item
                name="description"
                label="描述"
            >
                <Input.TextArea rows={4} placeholder="请输入待办事项描述" />
            </Form.Item>

            <Form.Item
                name="priority"
                label="优先级"
            >
                <Select>
                    <Select.Option value="low">低</Select.Option>
                    <Select.Option value="medium">中</Select.Option>
                    <Select.Option value="high">高</Select.Option>
                </Select>
            </Form.Item>

            <Form.Item
                name="dueDate"
                label="截止日期"
            >
                <DatePicker style={{ width: '100%' }} />
            </Form.Item>

            <Form.Item>
                <Space>
                    <Button type="primary" onClick={handleSubmit}>
                        {initialValues ? '更新' : '创建'}
                    </Button>
                    <Button onClick={onCancel}>取消</Button>
                </Space>
            </Form.Item>
        </Form>
    );
};

export default TodoForm; 