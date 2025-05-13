#!/bin/bash

# 数据库配置
DB_USER="root"
DB_PASSWORD="Zz1122"
DB_HOST="localhost"
DB_PORT="3306"
DB_NAME="todo_app_db"

# 模板目录
TEMPLATE_DIR=""

# 检查是否安装了 xorm
if ! command -v xorm &> /dev/null
then
    echo "xorm 未安装，请先安装 xorm"
    exit
fi

OUTPUT_DIR=""
xorm reverse mysql "$DB_USER:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME" $TEMPLATE_DIR $OUTPUT_DIR

echo "模型代码已生成"
