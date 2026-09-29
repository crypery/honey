-- MySQL dump 10.13
CREATE DATABASE wp_db;
USE wp_db;
INSERT INTO wp_users (ID, user_login, user_pass, user_email) VALUES
(1, 'admin', '$P$Bx9Yq2vVQmKpLr8sT1wZu4nA7cD3eF6gH0jI5kO9aB', 'admin@example.com');
INSERT INTO wp_options (option_name, option_value) VALUES
('siteurl', 'http://example.com'),
('admin_email', 'admin@example.com');
