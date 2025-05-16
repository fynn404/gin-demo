import React, { useEffect, useState } from 'react';
import { Table, Button, Space, Tag, Modal, message, Checkbox } from 'antd';
import { DeleteOutlined, EditOutlined } from '@ant-design/icons';
import { TodoItem } from '../types/todo';
import { todoService } from '../services/todoService';
import TodoForm from './TodoForm';
import dayjs from 'dayjs';

const TodoList: React.FC = () => {
    const [todos, setTodos] = useState<TodoItem[]>([]);
    const [total, setTotal] = useState(0);
    const [loading, setLoading] = useState(false);
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [editingTodo, setEditingTodo] = useState<TodoItem | null>(null);
    const [isModalVisible, setIsModalVisible] = useState(false);

    const fetchTodos = async () => {
        try {
            setLoading(true);
            const response = await todoService.listTodos(currentPage, pageSize);
            setTodos(response.todos);
            setTotal(response.total);
        } catch (error) {
            message.error('获取待办事项列表失败');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchTodos();
    }, [currentPage, pageSize]);

    const handleDelete = async (id: number) => {
        try {
            await todoService.deleteTodo(id);
            message.success('删除成功');
            fetchTodos();
        } catch (error) {
            message.error('删除失败');
        }
    };

    const handleEdit = (todo: TodoItem) => {
        setEditingTodo(todo);
        setIsModalVisible(true);
    };

    const handleStatusChange = async (todo: TodoItem) => {
        try {
            await todoService.changeTodoStatus(todo.id, { completed: !todo.completed });
            message.success('状态更新成功');
            fetchTodos();
        } catch (error) {
            message.error('状态更新失败');
        }
    };

    const handleModalSubmit = async (values: any) => {
        try {
            if (editingTodo) {
                await todoService.updateTodo(editingTodo.id, values);
                message.success('更新成功');
            } else {
                await todoService.createTodo(values);
                message.success('创建成功');
            }
            setIsModalVisible(false);
            setEditingTodo(null);
            fetchTodos();
        } catch (error) {
            message.error(editingTodo ? '更新失败' : '创建失败');
        }
    };

    const columns = [
        {
            title: '状态',
            dataIndex: 'completed',
            key: 'completed',
            render: (_: boolean, record: TodoItem) => (
                <Checkbox
                    checked={record.completed}
                    onChange={() => handleStatusChange(record)}
                />
            ),
        },
        {
            title: '标题',
            dataIndex: 'title',
            key: 'title',
        },
        {
            title: '描述',
            dataIndex: 'description',
            key: 'description',
        },
        {
            title: '优先级',
            dataIndex: 'priority',
            key: 'priority',
            render: (priority: string) => {
                const colorMap = {
                    low: 'green',
                    medium: 'orange',
                    high: 'red',
                };
                return <Tag color={colorMap[priority as keyof typeof colorMap]}>{priority}</Tag>;
            },
        },
        {
            title: '截止日期',
            dataIndex: 'dueDate',
            key: 'dueDate',
            render: (date: string) => date ? dayjs(date).format('YYYY-MM-DD') : '-',
        },
        {
            title: '操作',
            key: 'action',
            render: (_: any, record: TodoItem) => (
                <Space size="middle">
                    <Button
                        type="text"
                        icon={<EditOutlined />}
                        onClick={() => handleEdit(record)}
                    >
                        编辑
                    </Button>
                    <Button
                        type="text"
                        danger
                        icon={<DeleteOutlined />}
                        onClick={() => handleDelete(record.id)}
                    >
                        删除
                    </Button>
                </Space>
            ),
        },
    ];

    return (
        <div style={{ padding: '24px' }}>
            <div style={{ marginBottom: '16px' }}>
                <Button
                    type="primary"
                    onClick={() => {
                        setEditingTodo(null);
                        setIsModalVisible(true);
                    }}
                >
                    新建待办事项
                </Button>
            </div>

            <Table
                columns={columns}
                dataSource={todos}
                rowKey="id"
                loading={loading}
                pagination={{
                    current: currentPage,
                    pageSize: pageSize,
                    total: total,
                    onChange: (page, size) => {
                        setCurrentPage(page);
                        setPageSize(size || 10);
                    },
                }}
            />

            <Modal
                title={editingTodo ? '编辑待办事项' : '新建待办事项'}
                open={isModalVisible}
                onCancel={() => {
                    setIsModalVisible(false);
                    setEditingTodo(null);
                }}
                footer={null}
            >
                <TodoForm
                    initialValues={editingTodo}
                    onSubmit={handleModalSubmit}
                    onCancel={() => {
                        setIsModalVisible(false);
                        setEditingTodo(null);
                    }}
                />
            </Modal>
        </div>
    );
};

export default TodoList; 