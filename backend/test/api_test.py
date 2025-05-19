#!/usr/bin/env python3
import pprint

import requests
import json
from typing import Dict, Optional


class APITester:
    def __init__(self, base_url: str = "http://localhost:9527/api/v1"):
        self.base_url = base_url
        self.tokens = {
            'admin': None,
            'user': None,
        }
        self.test_users = {
            'admin': {'username': 'admin_test1', 'password': 'admin123', 'role': 'admin', 'nickname': 'Admin User',
                      'email': 'admin1@test.com'},
            'user': {'username': 'user_test1', 'password': 'user666', 'role': 'user', 'nickname': 'username123',
                        'email': 'student1@test.com'}
        }
        self.test_update_users = {
            'admin': {'username': 'admin_test2', 'password': 'admin234', 'role': 'admin', 'nickname': 'admin_name666',
                      'email': 'admin_new@test.com'},
            'user': {'username': 'user_test1_new', 'password': 'user666', 'role': 'user',
                        'nickname': 'username666New',
                        'email': 'student_new2@test.com'}
        }
        self.test_todo = {
            'title': 'Test Todo Item2222',
            'description': 'This is a test todo item',
            'priority': "low",
            'due_date': '2024-12-31T23:59:59Z'
        }
        self.created_todo_id = 2

    def make_request(self, method: str, endpoint: str, data: Optional[Dict] = None,
                     token: Optional[str] = None) -> requests.Response:
        url = f"{self.base_url}{endpoint}"
        headers = {'Content-Type': 'application/json'}
        # Print the request details
        print(f"\nRequest: {method} {url}")

        if data:
            print(f"Request: {json.dumps(data, indent=2)}")

        if token:
            headers['Authorization'] = f'Bearer {token}'
        print(f"Headers: {headers}")
        try:
            if method == 'GET':
                response = requests.get(url, headers=headers)
            elif method == 'POST':
                response = requests.post(url, json=data, headers=headers)
            elif method == 'PUT':
                response = requests.put(url, json=data, headers=headers)
            elif method == 'DELETE':
                response = requests.delete(url, headers=headers)
            elif method == 'PATCH':
                response = requests.patch(url, json=data,headers=headers)
            else:
                raise ValueError(f"Unsupported HTTP method: {method}")

            print(f"\n{method} {endpoint}")
            print(f"Status Code: {response.status_code}")
            # 将 JSON 字符串解析为 Python 对象
            data = json.loads(response.text)
            pprint.pprint(data)
            return response
        except requests.exceptions.RequestException as e:
            print(f"Error making request: {e}")
            return None

    def register_user(self, user_type: str) -> bool:
        user_data = self.test_users[user_type]
        response = self.make_request('POST', '/auth/register', user_data)
        return response and response.status_code in [201, 200]

    def login_user(self, user_type: str) -> bool:
        user_data = {
            'username': self.test_users[user_type]['username'],
            'password': self.test_users[user_type]['password']
        }
        response = self.make_request('POST', '/auth/login', user_data)
        if response and response.status_code == 200:
            self.tokens[user_type] = response.json().get('data',{}).get('access_token')
            return True
        return False

    def test_user_profile(self, user_type: str) -> bool:
        response = self.make_request('GET', '/users/profile', token=self.tokens[user_type])
        return response and response.status_code == 200

    def update_user_profile(self, user_type: str) -> bool:
        user_data = {
            'nickname': self.test_update_users[user_type]['nickname'],
            'email': self.test_update_users[user_type]['email'],
            # 'password': self.test_update_users[user_type]['password']
        }
        response = self.make_request('PUT', '/users/profile', data=user_data, token=self.tokens[user_type])
        return response and response.status_code == 200


    def create_todo(self, user_type: str) -> bool:
        """Create a new todo item"""
        response = self.make_request('POST', '/todos', data=self.test_todo, token=self.tokens[user_type])
        if response and response.status_code == 201:
            self.created_todo_id = response.json().get('data', {}).get('id')
            return True
        return False

    def list_todos(self, user_type: str, page: int = 1, size: int = 10) -> bool:
        """List todo items with pagination"""
        response = self.make_request('GET', f'/todos?page={page}&size={size}', token=self.tokens[user_type])
        return response and response.status_code == 200

    def get_todo_detail(self, user_type: str) -> bool:
        """Get details of a specific todo item"""
        if not self.created_todo_id:
            print("No todo item created yet")
            return False
        response = self.make_request('GET', f'/todos/{self.created_todo_id}', token=self.tokens[user_type])
        return response and response.status_code == 200

    def update_todo(self, user_type: str) -> bool:
        """Update a todo item"""
        if not self.created_todo_id:
            print("No todo item created yet")
            return False
        update_data = {
            'title': 'Updated Test Todo',
            'description': 'This is an updated test todo item',
            'priority': "high",
            'due_date': '2025-01-31T23:59:59Z'
        }
        response = self.make_request('PUT', f'/todos/{self.created_todo_id}', data=update_data,
                                     token=self.tokens[user_type])
        return response and response.status_code == 200

    def change_todo_status(self, user_type: str, completed: bool = False) -> bool:
        """Change the completion status of a todo item"""
        if not self.created_todo_id:
            print("No todo item created yet")
            return False
        status_data = {'completed': completed}
        response = self.make_request('PATCH', f'/todos/{self.created_todo_id}/status', data=status_data,
                                     token=self.tokens[user_type])
        return response and response.status_code == 200

    def delete_todo(self, user_type: str) -> bool:
        """Delete a todo item"""
        if not self.created_todo_id:
            print("No todo item created yet")
            return False
        response = self.make_request('DELETE', f'/todos/{self.created_todo_id}', token=self.tokens[user_type])
        return response and response.status_code == 200

    def cleanup(self) -> bool:
        return True

    def run_all_tests(self):
        return

if __name__ == "__main__":
    # # Create API tester instance
    tester = APITester()
    user_type = 'admin'
    tester.login_user(user_type)
    # tester.get_todo_detail(user_type)
    tester.delete_todo(user_type)
    # tester.change_todo_status(user_type,completed=True)
    # tester.get_todo_detail(user_type)
    # tester.get_todo_detail(user_type)
    # tester.update_todo(user_type)
    # tester.get_todo_detail(user_type)
    # tester.create_todo(user_type)
    # tester.list_todos(user_type,page=2,size=2)
    # tester.register_user(user_type)
    # tester.login_user(user_type)
    # tester.test_user_profile(user_type)
    # tester.update_user_profile(user_type)
    # tester.login_user(user_type)

