-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS todo_app_db DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE todo_app_db;

CREATE TABLE `users_tab` (
                             `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户唯一标识符，自增',
                             `username` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名，必须唯一',
                             `password_hash` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码哈希',
                             `role` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户角色',
                             `status` enum('enabled','disabled') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'enabled' COMMENT '用户状态',
                             `nickname` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户昵称',
                             `email` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户电子邮件',
                             `avatar_url` varchar(1024) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户头像URL',
                             `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                             `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                             PRIMARY KEY (`id`),
                             UNIQUE KEY `idx_users_username` (`username`),
                             CONSTRAINT `users_tab_chk_1` CHECK ((`role` in (_utf8mb4'user',_utf8mb4'admin')))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE `todo_items_tab` (
                                  `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY, -- 自增主键
                                  `user_id` BIGINT UNSIGNED, -- 外键，引用 users 表的 id
                                  `title` VARCHAR(255) NOT NULL, -- 任务标题，不能为空
                                  `description` TEXT, -- 任务描述
                                  `priority` ENUM('low', 'medium', 'high') NOT NULL DEFAULT 'medium', -- 优先级
                                  `due_date` DATETIME, -- 截止日期
                                  `completed` BOOLEAN DEFAULT FALSE, -- 是否完成
                                  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间，默认当前时间
                                  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 更新时间，自动更新
                                  KEY `idx_user_id` (`user_id`), -- user_id索引
                                  FOREIGN KEY (`user_id`) REFERENCES `users_tab`(`id`) ON DELETE CASCADE -- 外键约束
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;