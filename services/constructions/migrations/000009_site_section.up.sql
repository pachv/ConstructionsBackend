-- =========================
-- SITE SECTIONS (landing)
-- =========================

CREATE TABLE IF NOT EXISTS site_sections (
    id               VARCHAR(255),
    title            VARCHAR(255),
    label            VARCHAR(255),
    slug             VARCHAR(255),
    image_url        VARCHAR(1024),
    advanteges_text  TEXT,
    has_gallery      BOOLEAN DEFAULT FALSE,
    has_catalog      BOOLEAN DEFAULT FALSE,
    created_at       TIMESTAMP DEFAULT now()
);

-- преимущества массивом
CREATE TABLE IF NOT EXISTS site_section_advanteges (
    id         VARCHAR(255),
    section_id VARCHAR(255),
    text       VARCHAR(1024),
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT now()
);

-- Галерея для конкретного section
CREATE TABLE IF NOT EXISTS site_section_gallery (
    id         VARCHAR(255),
    section_id VARCHAR(255),
    name       VARCHAR(255),
    url        VARCHAR(1024),
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT now()
);

-- Категории внутри секции
CREATE TABLE IF NOT EXISTS site_section_catalog_categories (
    id          VARCHAR(255),
    section_id  VARCHAR(255),
    category_id VARCHAR(255),
    sort_order  INT DEFAULT 0,
    created_at  TIMESTAMP DEFAULT now()
);

-- =========================
-- CATALOG ITEMS FOR SITE SECTION
-- =========================

CREATE TABLE IF NOT EXISTS site_section_catalog_items (
    id          VARCHAR(255),
    section_id  VARCHAR(255),
    category_id VARCHAR(255),
    title       VARCHAR(255),
    price_rub   INT DEFAULT 0,
    image_url   VARCHAR(1024),
    sort_order  INT DEFAULT 0,
    created_at  TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS site_section_catalog_item_specs (
    id         VARCHAR(255),
    item_id    VARCHAR(255),
    key        VARCHAR(255),
    value      VARCHAR(255),
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS site_section_catalog_item_badges (
    id         VARCHAR(255),
    item_id    VARCHAR(255),
    badge      VARCHAR(255),
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT now()
);

-- =========================
-- SEED
-- =========================

INSERT INTO site_sections (id, title, label, slug, image_url, advanteges_text, has_gallery, has_catalog) VALUES
('sec-landing-metal', 'Металлоконструкции', 'Металлоконструкции', 'metall', '/images/sections/build-1.jpg', '', TRUE, TRUE),
('sec-landing-bsu',   'БСУ',               'БСУ',               'bsu',    '/images/sections/build-2.jpg', '', TRUE, TRUE),
('sec-landing-bps',   'БПС',               'БПС',               'bps',    '/images/sections/build-3.jpg', '', TRUE, TRUE);

INSERT INTO site_section_advanteges (id, section_id, text, sort_order) VALUES
('adv-metal-1', 'sec-landing-metal', '', 1),
('adv-metal-2', 'sec-landing-metal', '', 2);

INSERT INTO site_section_gallery (id, section_id, name, url, sort_order) VALUES
('gal-metal-1', 'sec-landing-metal', 'Цех снаружи',   '/sections/gallery/picture/metall-1.jpg', 1),
('gal-metal-2', 'sec-landing-metal', 'Каркас внутри', '/sections/gallery/picture/metall-2.jpg', 2),
('gal-metal-3', 'sec-landing-metal', 'Ангар',         '/sections/gallery/picture/metall-3.jpg', 3),
('gal-bsu-1',   'sec-landing-bsu',   'БСУ 1',         '/sections/gallery/picture/bsu-1.jpg',    1);

INSERT INTO site_section_catalog_categories (id, section_id, category_id, sort_order) VALUES
('ssc-cat-1', 'sec-landing-metal', 'cat-instrument', 1),
('ssc-cat-2', 'sec-landing-metal', 'cat-power',      2);

INSERT INTO site_section_catalog_items
(id, section_id, category_id, title, price_rub, image_url, sort_order)
VALUES
('prd-1', 'sec-landing-metal', 'cat-instrument', 'Блок верхний доборный', 2484, '/catalog/metall/blocks/top.jpg', 1);

INSERT INTO site_section_catalog_item_badges (id, item_id, badge, sort_order) VALUES
('bad-1', 'prd-1', 'В30', 1),
('bad-2', 'prd-1', 'F300', 2),
('bad-3', 'prd-1', 'W8', 3);

INSERT INTO site_section_catalog_item_specs (id, item_id, key, value, sort_order) VALUES
('spec-1', 'prd-1', 'Марка бетона', 'B30', 1),
('spec-2', 'prd-1', 'Морозостойкость', 'F300', 2),
('spec-3', 'prd-1', 'Водонепроницаемость', 'W8', 3),
('spec-4', 'prd-1', 'Применение', 'доборные элементы', 4);
