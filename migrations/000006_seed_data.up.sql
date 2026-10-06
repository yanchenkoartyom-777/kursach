-- Seed data for development and testing

-- Insert sample admin user (password: 123456)
-- Password hash generated with bcrypt cost 10
INSERT INTO users (email, password_hash, first_name, last_name, role)
VALUES ('formed86@gmail.com', '$2a$10$3XZdryAisgPUMjRx5Ny6bOwW5MwhnTv2/4bsjjw/0ZLZeOHi9GX.a', 'Admin', 'User', 'admin')
ON CONFLICT (email) DO NOTHING;

-- Insert sample client users (password: 123456)
INSERT INTO users (email, password_hash, first_name, last_name, role)
VALUES 
    ('baloban.oleg2016@yandex.ru', '$2a$10$3XZdryAisgPUMjRx5Ny6bOwW5MwhnTv2/4bsjjw/0ZLZeOHi9GX.a', 'Олег', 'Балобан', 'client'),
    ('ivan.petrov@example.com', '$2a$10$3XZdryAisgPUMjRx5Ny6bOwW5MwhnTv2/4bsjjw/0ZLZeOHi9GX.a', 'Иван', 'Петров', 'client'),
    ('maria.sidorova@example.com', '$2a$10$3XZdryAisgPUMjRx5Ny6bOwW5MwhnTv2/4bsjjw/0ZLZeOHi9GX.a', 'Мария', 'Сидорова', 'client'),
    ('alex.kozlov@example.com', '$2a$10$3XZdryAisgPUMjRx5Ny6bOwW5MwhnTv2/4bsjjw/0ZLZeOHi9GX.a', 'Александр', 'Козлов', 'client')
ON CONFLICT (email) DO NOTHING;

-- Insert sample cars with images
INSERT INTO cars (make, model, year, price, color, engine_type, transmission, mileage, description, image_url, is_available)
VALUES 
    ('Toyota', 'Camry', 2023, 28500.00, 'Silver', 'Hybrid', 'Automatic', 5000, 'Надежный седан с гибридным двигателем', 'https://avatars.mds.yandex.net/get-verba/997355/2a0000018e0fd1de72d4ecb8fef86d78a72e/cattouchret', true),
    ('Honda', 'Accord', 2023, 27800.00, 'Black', 'Gasoline', 'Automatic', 3000, 'Комфортный семейный автомобиль', 'https://avatars.mds.yandex.net/get-autoru-vos/16813311/1880abd06e03e3c8450eb13c206135ec/1200x900n', true),
    ('BMW', '3 Series', 2022, 45000.00, 'White', 'Gasoline', 'Automatic', 15000, 'Премиальный спортивный седан', 'https://avatars.mds.yandex.net/get-autoru-vos/5227974/5d2156d05934119cd6ab4201be65bc75/1200x900n', true),
    ('Mercedes-Benz', 'C-Class', 2022, 48500.00, 'Blue', 'Diesel', 'Automatic', 12000, 'Элегантный и мощный автомобиль', 'https://avatars.mds.yandex.net/get-autoru-vos/2171774/02292a6419349ba0559f9fe131ac8a2e/1200x900n', true),
    ('Audi', 'A4', 2023, 42000.00, 'Gray', 'Gasoline', 'Automatic', 8000, 'Современный дизайн и технологии', 'https://autocapital.by/wp-content/uploads/2024/07/image_2024-07-12_12-43-27.jpg', true),
    ('Tesla', 'Model 3', 2023, 52000.00, 'Red', 'Electric', 'Automatic', 2000, 'Электромобиль будущего', 'https://digitalassets.tesla.com/tesla-contents/image/upload/f_auto,q_auto/Model-3-Main-Hero-Desktop-LHD.jpg', true),
    ('Ford', 'Mustang', 2022, 38000.00, 'Yellow', 'Gasoline', 'Manual', 10000, 'Легендарный американский маслкар', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcS1g7vL6QemI0_x5yCOtd--vqWR2BnXSRKr1Q&s', true),
    ('Volkswagen', 'Passat', 2023, 32000.00, 'Silver', 'Diesel', 'Automatic', 6000, 'Практичный и экономичный седан', 'https://avcdn.av.by/cargeneration/0001/2630/8575.jpg', true),
    ('Mazda', 'CX-5', 2023, 35000.00, 'Blue', 'Gasoline', 'Automatic', 4000, 'Стильный кроссовер для семьи', 'https://avcdn.av.by/cargeneration/0000/0025/0423.jpeg', true),
    ('Hyundai', 'Sonata', 2023, 26500.00, 'White', 'Gasoline', 'Automatic', 7000, 'Доступный и надежный седан', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRSZtW57WJcTdzlNgPJyOakMB8PZXsCUqVM0g&s', true)
ON CONFLICT DO NOTHING;


-- Insert sample orders with historical data (last 6 months)
INSERT INTO orders (user_id, total_amount, status, created_at)
VALUES 
    -- 6 months ago
    (2, 28500.00, 'completed', NOW() - INTERVAL '180 days'),
    (3, 45000.00, 'completed', NOW() - INTERVAL '175 days'),
    (4, 52000.00, 'completed', NOW() - INTERVAL '170 days'),
    -- 5 months ago
    (2, 27800.00, 'completed', NOW() - INTERVAL '150 days'),
    (3, 48500.00, 'completed', NOW() - INTERVAL '145 days'),
    (5, 42000.00, 'completed', NOW() - INTERVAL '140 days'),
    -- 4 months ago
    (2, 35000.00, 'completed', NOW() - INTERVAL '120 days'),
    (4, 38000.00, 'completed', NOW() - INTERVAL '115 days'),
    (3, 32000.00, 'completed', NOW() - INTERVAL '110 days'),
    -- 3 months ago
    (5, 28500.00, 'completed', NOW() - INTERVAL '90 days'),
    (2, 45000.00, 'completed', NOW() - INTERVAL '85 days'),
    (4, 52000.00, 'completed', NOW() - INTERVAL '80 days'),
    -- 2 months ago
    (3, 27800.00, 'completed', NOW() - INTERVAL '60 days'),
    (5, 48500.00, 'completed', NOW() - INTERVAL '55 days'),
    (2, 42000.00, 'completed', NOW() - INTERVAL '50 days'),
    -- 1 month ago
    (4, 35000.00, 'completed', NOW() - INTERVAL '30 days'),
    (3, 38000.00, 'completed', NOW() - INTERVAL '25 days'),
    (5, 32000.00, 'processing', NOW() - INTERVAL '20 days'),
    -- Current month
    (2, 28500.00, 'processing', NOW() - INTERVAL '15 days'),
    (4, 45000.00, 'completed', NOW() - INTERVAL '10 days'),
    (3, 52000.00, 'completed', NOW() - INTERVAL '5 days'),
    (5, 35000.00, 'pending', NOW() - INTERVAL '2 days')
ON CONFLICT DO NOTHING;

-- Insert sample order items
INSERT INTO order_items (order_id, car_id, price, created_at)
VALUES 
    -- 6 months ago
    (1, 1, 28500.00, NOW() - INTERVAL '180 days'),
    (2, 3, 45000.00, NOW() - INTERVAL '175 days'),
    (3, 6, 52000.00, NOW() - INTERVAL '170 days'),
    -- 5 months ago
    (4, 2, 27800.00, NOW() - INTERVAL '150 days'),
    (5, 4, 48500.00, NOW() - INTERVAL '145 days'),
    (6, 5, 42000.00, NOW() - INTERVAL '140 days'),
    -- 4 months ago
    (7, 9, 35000.00, NOW() - INTERVAL '120 days'),
    (8, 7, 38000.00, NOW() - INTERVAL '115 days'),
    (9, 8, 32000.00, NOW() - INTERVAL '110 days'),
    -- 3 months ago
    (10, 1, 28500.00, NOW() - INTERVAL '90 days'),
    (11, 3, 45000.00, NOW() - INTERVAL '85 days'),
    (12, 6, 52000.00, NOW() - INTERVAL '80 days'),
    -- 2 months ago
    (13, 2, 27800.00, NOW() - INTERVAL '60 days'),
    (14, 4, 48500.00, NOW() - INTERVAL '55 days'),
    (15, 5, 42000.00, NOW() - INTERVAL '50 days'),
    -- 1 month ago
    (16, 9, 35000.00, NOW() - INTERVAL '30 days'),
    (17, 7, 38000.00, NOW() - INTERVAL '25 days'),
    (18, 8, 32000.00, NOW() - INTERVAL '20 days'),
    -- Current month
    (19, 1, 28500.00, NOW() - INTERVAL '15 days'),
    (20, 3, 45000.00, NOW() - INTERVAL '10 days'),
    (21, 6, 52000.00, NOW() - INTERVAL '5 days'),
    (22, 9, 35000.00, NOW() - INTERVAL '2 days')
ON CONFLICT DO NOTHING;
