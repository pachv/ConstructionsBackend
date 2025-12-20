-- 001_gallery.sql

CREATE TABLE IF NOT EXISTS gallery_categories (
  id         VARCHAR(255),
  title      VARCHAR(255),
  slug       VARCHAR(255),
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS gallery_photos (
  id          VARCHAR(255),
  category_id VARCHAR(255),
  alt         VARCHAR(255),
  image_path  VARCHAR(1024),
  sort_order  INT,
  created_at  TIMESTAMP NOT NULL DEFAULT now()
);

-- Категории
INSERT INTO gallery_categories (id, title, slug) VALUES
('gal-baths',     'Бани',      'baths'),
('gal-houses',    'Дома',      'houses'),
('gal-cottages',  'Коттеджи',  'cottages');

-- Фото (файлы лежат в ./uploads/gallery/)
INSERT INTO gallery_photos (id, category_id, alt, image_path, sort_order) VALUES
-- ===== БАНИ =====
('baths-1', 'gal-baths', 'Баня 1', 'baths-1.jpg', 10),
('baths-2', 'gal-baths', 'Баня 2', 'baths-2.jpg', 20),
('baths-3', 'gal-baths', 'Баня 3', 'baths-3.jpg', 30),
('baths-4', 'gal-baths', 'Баня 4', 'baths-4.jpg', 40),
('baths-5', 'gal-baths', 'Баня 5', 'baths-5.jpg', 50),
('baths-6', 'gal-baths', 'Баня 6', 'baths-6.jpg', 60),
('baths-7', 'gal-baths', 'Баня 7', 'baths-7.jpg', 70),
('baths-8', 'gal-baths', 'Баня 8', 'baths-8.jpg', 80),
('baths-9', 'gal-baths', 'Баня 9', 'baths-9.jpg', 90),

-- ===== КОТТЕДЖИ =====
('cottages-1', 'gal-cottages', 'Коттедж 1', 'cottages-1.jpg', 10),
('cottages-2', 'gal-cottages', 'Коттедж 2', 'cottages-2.jpg', 20),
('cottages-4', 'gal-cottages', 'Коттедж 4', 'cottages-4.jpg', 30),
('cottages-5', 'gal-cottages', 'Коттедж 5', 'cottages-5.jpg', 40),
('cottages-6', 'gal-cottages', 'Коттедж 6', 'cottages-6.jpg', 50),
('cottages-7', 'gal-cottages', 'Коттедж 7', 'cottages-7.jpg', 60),
('cottages-8', 'gal-cottages', 'Коттедж 8', 'cottages-8.jpg', 70),

-- ===== ДОМА =====
('houses-1', 'gal-houses', 'Дом 1', 'houses-1.jpg', 10),
('houses-2', 'gal-houses', 'Дом 2', 'houses-2.jpg', 20),
('houses-3', 'gal-houses', 'Дом 3', 'houses-3.jpg', 30),
('houses-4', 'gal-houses', 'Дом 4', 'houses-4.jpg', 40),
('houses-5', 'gal-houses', 'Дом 5', 'houses-5.jpg', 50),
('houses-6', 'gal-houses', 'Дом 6', 'houses-6.jpg', 60),
('houses-7', 'gal-houses', 'Дом 7', 'houses-7.jpg', 70),
('houses-8', 'gal-houses', 'Дом 8', 'houses-8.jpg', 80);
