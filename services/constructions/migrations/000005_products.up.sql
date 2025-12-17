CREATE TABLE catalog_categories (
    id VARCHAR(255),
    title VARCHAR(255),
    slug VARCHAR(255),
    image_path VARCHAR(1024),
    created_at TIMESTAMP
);

INSERT INTO catalog_categories (id, title, slug, image_path, created_at) VALUES
(
    'cat-instrument',
    'Инструмент',
    'instrument',
    '/categories/picture/tools.jpg',
    now()
),
(
    'cat-floor',
    'Напольные покрытия',
    'floor',
    '/categories/picture/floor.jpg',
    now()
),
(
    'cat-tile',
    'Плитка',
    'tiles',
    '/categories/picture/tile.jpg',
    now()
),
(
    'cat-mixes',
    'Сухие смеси и грунтовки',
    'mixes',
    '/categories/picture/mixes.jpg',
    now()
),
(
    'cat-power',
    'Электроинструменты',
    'power-tools',
    '/categories/picture/power-tools.jpg',
    now()
);

CREATE TABLE catalog_sections (
    id VARCHAR(255),
    title VARCHAR(255),
    slug VARCHAR(255),
    image_path VARCHAR(1024),
    created_at TIMESTAMP
);

-- связь: одна секция может принадлежать многим категориям
CREATE TABLE catalog_category_sections (
    category_id VARCHAR(255),
    section_id  VARCHAR(255),
    created_at  TIMESTAMP
);

