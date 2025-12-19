-- 001_gallery.sql

CREATE TABLE IF NOT EXISTS gallery_categories (
  id         VARCHAR(255),
  title      VARCHAR(255) ,
  slug       VARCHAR(255),
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS gallery_photos (
  id          VARCHAR(255) ,
  category_id VARCHAR(255),
  alt         VARCHAR(255) ,
  image_path  VARCHAR(1024),
  sort_order  INT,
  created_at  TIMESTAMP NOT NULL DEFAULT now()
);

-- Демо категории (как на скрине)
INSERT INTO gallery_categories (id, title, slug) VALUES
('gal-bani',     'Бани',     'bani'),
('gal-doma',     'Дома',     'doma'),
('gal-kottedzhi','Коттеджи', 'kottedzhi');

-- Демо фото (пути под твой файловый хендлер /gallery/picture/:name)
-- Файлы положи в: ./uploads/gallery/
INSERT INTO gallery_photos (id, category_id, alt, image_path, sort_order) VALUES
('bani-1', 'gal-bani', 'Баня 1', 'build-house-bg.jpg', 10),
('bani-2', 'gal-bani', 'Баня 2', 'build-house-bg.jpg', 20),
('bani-3', 'gal-bani', 'Баня 3', 'build-house-bg.jpg', 30),

('doma-1', 'gal-doma', 'Дом 1', 'build-house-bg.jpg', 10),
('doma-2', 'gal-doma', 'Дом 2', 'build-house-bg.jpg', 20),

('kot-1', 'gal-kottedzhi', 'Коттедж 1', 'build-house-bg.jpg', 10),
('kot-2', 'gal-kottedzhi', 'Коттедж 2', 'build-house-bg.jpg', 20);
