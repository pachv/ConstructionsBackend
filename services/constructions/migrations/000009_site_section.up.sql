-- =========================
-- SITE SECTIONS (landing)
-- =========================

CREATE TABLE IF NOT EXISTS site_sections (
    id               VARCHAR(255) PRIMARY KEY,
    title            VARCHAR(255) NOT NULL,
    label            VARCHAR(255) NOT NULL,
    slug             VARCHAR(255) NOT NULL UNIQUE,
    image_url        VARCHAR(1024) NOT NULL,
    advanteges_text  TEXT DEFAULT '',
    has_gallery      BOOLEAN DEFAULT FALSE,
    has_catalog      BOOLEAN DEFAULT FALSE,
    created_at       TIMESTAMP DEFAULT now()
);

-- преимущества массивом
CREATE TABLE IF NOT EXISTS site_section_advanteges (
    id         VARCHAR(255) PRIMARY KEY,
    section_id VARCHAR(255) NOT NULL REFERENCES site_sections(id) ON DELETE CASCADE,
    text       VARCHAR(1024) NOT NULL DEFAULT '',
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_site_section_advanteges_section_sort
ON site_section_advanteges(section_id, sort_order);

-- Галерея для конкретного section
CREATE TABLE IF NOT EXISTS site_section_gallery (
    id         VARCHAR(255) PRIMARY KEY,
    section_id VARCHAR(255) NOT NULL REFERENCES site_sections(id) ON DELETE CASCADE,
    name       VARCHAR(255) NOT NULL DEFAULT '',
    url        VARCHAR(1024) NOT NULL,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_site_section_gallery_section_sort
ON site_section_gallery(section_id, sort_order);

-- Категории внутри секции (ХРАНИМ title/slug ТУТ)
CREATE TABLE IF NOT EXISTS site_section_catalog_categories (
    id          VARCHAR(255) PRIMARY KEY,
    section_id  VARCHAR(255) NOT NULL REFERENCES site_sections(id) ON DELETE CASCADE,
    category_id VARCHAR(255) NOT NULL,      -- например: "cat-instrument"
    title       VARCHAR(255) NOT NULL,      -- например: "Инструмент"
    slug        VARCHAR(255) NOT NULL,      -- например: "instrument"
    sort_order  INT DEFAULT 0,
    created_at  TIMESTAMP DEFAULT now(),
    UNIQUE(section_id, category_id)
);

CREATE INDEX IF NOT EXISTS idx_ssc_categories_section_sort
ON site_section_catalog_categories(section_id, sort_order);

-- =========================
-- CATALOG ITEMS FOR SITE SECTION
-- =========================

CREATE TABLE IF NOT EXISTS site_section_catalog_items (
    id          VARCHAR(255) PRIMARY KEY,
    section_id  VARCHAR(255) NOT NULL REFERENCES site_sections(id) ON DELETE CASCADE,
    category_id VARCHAR(255) NOT NULL, -- ссылается на site_section_catalog_categories.category_id (логически)
    title       VARCHAR(255) NOT NULL,
    price_rub   INT DEFAULT 0,
    image_url   VARCHAR(1024) NOT NULL DEFAULT '',
    sort_order  INT DEFAULT 0,
    created_at  TIMESTAMP DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_ssc_items_section_sort
ON site_section_catalog_items(section_id, sort_order);

CREATE INDEX IF NOT EXISTS idx_ssc_items_section_category
ON site_section_catalog_items(section_id, category_id);

CREATE TABLE IF NOT EXISTS site_section_catalog_item_specs (
    id         VARCHAR(255) PRIMARY KEY,
    item_id    VARCHAR(255) NOT NULL REFERENCES site_section_catalog_items(id) ON DELETE CASCADE,
    key        VARCHAR(255) NOT NULL,
    value      VARCHAR(255) NOT NULL,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_ssc_item_specs_item_sort
ON site_section_catalog_item_specs(item_id, sort_order);

CREATE TABLE IF NOT EXISTS site_section_catalog_item_badges (
    id         VARCHAR(255) PRIMARY KEY,
    item_id    VARCHAR(255) NOT NULL REFERENCES site_section_catalog_items(id) ON DELETE CASCADE,
    badge      VARCHAR(255) NOT NULL,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_ssc_item_badges_item_sort
ON site_section_catalog_item_badges(item_id, sort_order);

-- =========================
-- SEED
-- =========================

INSERT INTO site_sections (id, title, label, slug, image_url, advanteges_text, has_gallery, has_catalog) VALUES
('sec-landing-metal', 'Металлоконструкции', 'Металлоконструкции', 'metall', 'build-1.jpg', '', TRUE, TRUE),
('sec-landing-bsu',   'БСУ',               'БСУ',               'bsu',    'build-2.jpg', '', TRUE, TRUE),
('sec-landing-bps',   'БПС',               'БПС',               'bps',    'build-3.jpg', '', FALSE, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO site_section_advanteges (id, section_id, text, sort_order) VALUES
('adv-metal-1', 'sec-landing-metal', '', 1),
('adv-metal-2', 'sec-landing-metal', '', 2)
ON CONFLICT (id) DO NOTHING;

INSERT INTO site_section_gallery (id, section_id, name, url, sort_order) VALUES
('gal-metal-1', 'sec-landing-metal', 'Цех снаружи',   'baths-1.jpg', 1),
('gal-metal-2', 'sec-landing-metal', 'Каркас внутри', 'baths-2.jpg', 2),
('gal-metal-3', 'sec-landing-metal', 'Ангар',         'baths-3.jpg', 3),
('gal-bsu-1',   'sec-landing-bsu',   'БСУ 1',         'baths-4.jpg', 1)
ON CONFLICT (id) DO NOTHING;

-- ВАЖНО: теперь нужны title/slug
INSERT INTO site_section_catalog_categories (id, section_id, category_id, title, slug, sort_order) VALUES
('ssc-cat-1', 'sec-landing-metal', 'cat-instrument', 'Инструмент', 'instrument', 1),
('ssc-cat-2', 'sec-landing-metal', 'cat-power',      'Силовое',    'power',      2)
ON CONFLICT (id) DO NOTHING;

INSERT INTO site_section_catalog_items
(id, section_id, category_id, title, price_rub, image_url, sort_order)
VALUES
('prd-1', 'sec-landing-metal', 'cat-instrument', 'Блок верхний доборный', 2484, 'floor.jpg', 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO site_section_catalog_item_badges (id, item_id, badge, sort_order) VALUES
('bad-1', 'prd-1', 'В30', 1),
('bad-2', 'prd-1', 'F300', 2),
('bad-3', 'prd-1', 'W8', 3)
ON CONFLICT (id) DO NOTHING;

INSERT INTO site_section_catalog_item_specs (id, item_id, key, value, sort_order) VALUES
('spec-1', 'prd-1', 'Марка бетона', 'B30', 1),
('spec-2', 'prd-1', 'Морозостойкость', 'F300', 2),
('spec-3', 'prd-1', 'Водонепроницаемость', 'W8', 3),
('spec-4', 'prd-1', 'Применение', 'доборные элементы', 4)
ON CONFLICT (id) DO NOTHING;
