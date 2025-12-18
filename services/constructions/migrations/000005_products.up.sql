-- =========================
-- CATEGORIES
-- =========================
CREATE TABLE catalog_categories (
    id VARCHAR(255),
    title VARCHAR(255),
    slug VARCHAR(255),
    image_path VARCHAR(1024),
    created_at TIMESTAMP
);

INSERT INTO catalog_categories (id, title, slug, image_path, created_at) VALUES
('cat-instrument','Инструмент','instrument','/categories/picture/tools.jpg',now()),
('cat-floor','Напольные покрытия','floor','/categories/picture/floor.jpg',now()),
('cat-tile','Плитка','tiles','/categories/picture/tile.jpg',now()),
('cat-mixes','Сухие смеси и грунтовки','mixes','/categories/picture/mixes.jpg',now()),
('cat-power','Электроинструменты','power-tools','/categories/picture/power-tools.jpg',now());

-- =========================
-- SECTIONS
-- =========================
CREATE TABLE catalog_sections (
    id VARCHAR(255),
    title VARCHAR(255),
    slug VARCHAR(255),
    image_path VARCHAR(1024),
    created_at TIMESTAMP
);

CREATE TABLE catalog_category_sections (
    category_id VARCHAR(255),
    section_id  VARCHAR(255),
    created_at  TIMESTAMP
);

INSERT INTO catalog_sections (id, title, slug, image_path, created_at) VALUES
('sec-measuring',   'Измерительный инструмент',          'measuring',    '/images/sections/measuring.png',    now()),
('sec-tile',        'Инструмент для укладки плитки',     'tile',         '/images/sections/tile.png',         now()),
('sec-locksmith',   'Слесарные инструменты',             'locksmith',    '/images/sections/locksmith.png',    now()),
('sec-laminate',    'Ламинат',                           'laminate',     '/images/sections/laminate.png',     now()),
('sec-linoleum',    'Линолеум',                          'linoleum',     '/images/sections/linoleum.png',     now()),
('sec-parquet',     'Паркетная доска',                   'parquet',      '/images/sections/parquet.png',      now()),
('sec-floor-tile',  'Напольная плитка',                  'floor-tile',   '/images/sections/floor-tile.png',   now()),
('sec-wall-tile',   'Настенная плитка',                  'wall-tile',    '/images/sections/wall-tile.png',    now()),
('sec-mosaic',      'Плитка-мозаика',                    'mosaic',       '/images/sections/mosaic.png',       now()),
('sec-cement',      'Цемент',                            'cement',       '/images/sections/cement.png',       now()),
('sec-putty',       'Шпатлевка',                         'putty',        '/images/sections/putty.png',        now()),
('sec-plaster',     'Штукатурка',                        'plaster',      '/images/sections/plaster.png',      now()),
('sec-drill',       'Дрель',                             'drill',        '/images/sections/drill.png',        now()),
('sec-heat-gun',    'Фен технический',                   'heat-gun',     '/images/sections/heat-gun.png',     now()),
('sec-screwdriver', 'Шуруповерт',                        'screwdriver',  '/images/sections/screwdriver.png',  now());

-- связи category -> sections
INSERT INTO catalog_category_sections (category_id, section_id, created_at) VALUES
-- Инструмент
('cat-instrument', 'sec-measuring', now()),
('cat-instrument', 'sec-tile', now()),
('cat-instrument', 'sec-locksmith', now()),
-- Напольные покрытия
('cat-floor', 'sec-laminate', now()),
('cat-floor', 'sec-linoleum', now()),
('cat-floor', 'sec-parquet', now()),
-- Плитка
('cat-tile', 'sec-floor-tile', now()),
('cat-tile', 'sec-wall-tile', now()),
('cat-tile', 'sec-mosaic', now()),
-- Сухие смеси
('cat-mixes', 'sec-cement', now()),
('cat-mixes', 'sec-putty', now()),
('cat-mixes', 'sec-plaster', now()),
-- Электроинструменты
('cat-power', 'sec-drill', now()),
('cat-power', 'sec-heat-gun', now()),
('cat-power', 'sec-screwdriver', now());

-- =========================
-- PRODUCTS + BADGES
-- =========================
CREATE TABLE catalog_products (
    id VARCHAR(255),
    title VARCHAR(255),
    slug VARCHAR(255),

    category_slug VARCHAR(255),
    section_slug VARCHAR(255),

    brand VARCHAR(255),
    type VARCHAR(255),

    price INT,
    old_price INT,

    in_stock BOOLEAN,

    sale_percent INT,
    image_path VARCHAR(1024),

    created_at TIMESTAMP
);

CREATE TABLE product_badges (
    id VARCHAR(255),
    code VARCHAR(255),
    title VARCHAR(255),
    created_at TIMESTAMP
);

CREATE TABLE product_badge_links (
    product_id VARCHAR(255),
    badge_id   VARCHAR(255),
    created_at TIMESTAMP
);

INSERT INTO product_badges (id, code, title, created_at) VALUES
('badge-hit', 'hit', 'Хит', now()),
('badge-new', 'new', 'Новинка', now()),
('badge-sale', 'sale', 'Акция', now()),
('badge-recommended', 'recommended', 'Советуем', now());

INSERT INTO catalog_products (
  id, title, slug,
  category_slug, section_slug,
  brand, type,
  price, old_price,
  in_stock,
  sale_percent,
  image_path,
  created_at
) VALUES
-- instrument / measuring
('p-1','Дальномер лазерный BOSCH Professional GLM 40','bosch-glm-40','instrument','measuring','BOSCH','дальномер',4500,4900,true,-8,'/img/products/bosch-glm40.png',now()),
('p-2','Дальномер ELITECH ЛД 40Н лазерный','elitech-ld-40n','instrument','measuring','ELITECH','дальномер',2600,2600,true,0,'/img/products/elitech-ld40n.png',now()),
('p-3','Лазерный нивелир BOSCH Universal Level 2 Set','bosch-universal-level-2','instrument','measuring','BOSCH','нивелир',4500,4500,true,0,'/img/products/bosch-level2.png',now()),

-- instrument / tile
('p-4','Зубчатая гладилка LUX-TOOLS Classic сталь 130×270 мм','lux-trowel-130x270','instrument','tile','LUX-TOOLS','гладилка',229,229,true,0,'/img/products/trowel.png',now()),
('p-5','Захват LUX вакуумный двойной 120 мм','lux-vacuum-double-120','instrument','tile','LUX-TOOLS','захват',1130,1450,true,-22,'/img/products/vacuum-double.png',now()),
('p-6','Вакуумный захват LUX-TOOLS Classic тройной','lux-vacuum-triple','instrument','tile','LUX-TOOLS','захват',1450,1450,true,0,'/img/products/vacuum-triple.png',now()),

-- instrument / locksmith
('p-7','Молоток кровельщика LUX-TOOLS 600 гр','lux-roof-hammer-600','instrument','locksmith','LUX-TOOLS','молоток',1500,1500,true,0,'/img/products/hammer.png',now()),
('p-8','Молоток слесарный дерево 300 гр','hammer-300-wood','instrument','locksmith','PR-VOLGA','молоток',90,120,true,-25,'/img/products/hammer-wood.png',now()),
('p-9','Набор отверток LUX-TOOLS Classic 6 предметов','lux-screwdrivers-6','instrument','locksmith','LUX-TOOLS','отвертка крестовая',710,710,true,0,'/img/products/screwdrivers.png',now()),

-- floor / laminate
('p-floor-1','Ламинат «Дуб Монгольский» 33 класс','laminate-dub-mongolskiy-33','floor','laminate','ARTENS','ламинат',480,480,true,0,'/img/products/laminate-1.png',now()),
('p-floor-2','Ламинат «Дуб Ривер» 31 класс','laminate-dub-river-31','floor','laminate','ARTENS','ламинат',228,280,true,-18,'/img/products/laminate-2.png',now()),
('p-floor-3','Ламинат «Дуб Хадсон» 31 класс','laminate-dub-hudson-31','floor','laminate','ARTENS','ламинат',228,228,true,0,'/img/products/laminate-3.png',now()),

-- floor / linoleum
('p-floor-4','Линолеум «Дуб Венский» 32 класс 3 м','linoleum-dub-venskiy-32','floor','linoleum','TARKETT','линолеум',508,508,true,0,'/img/products/linoleum-1.png',now()),
('p-floor-5','Линолеум «Дуб Маренго» 21 класс 3 м','linoleum-dub-marengo-21','floor','linoleum','TARKETT','линолеум',286,560,true,-48,'/img/products/linoleum-2.png',now()),
('p-floor-6','Линолеум «Дуб Прованс» 32 класс 3.5 м','linoleum-dub-provans-32','floor','linoleum','TARKETT','линолеум',614,614,true,0,'/img/products/linoleum-3.png',now()),

-- floor / parquet
('p-floor-7','Паркетная доска «Дуб Натуральный»','parquet-dub-natural','floor','parquet','BARLINEK','паркетная доска',3200,3200,true,0,'/img/products/parquet-1.png',now()),
('p-floor-8','Паркетная доска «Дуб Селект»','parquet-dub-select','floor','parquet','BARLINEK','паркетная доска',3650,3650,false,0,'/img/products/parquet-2.png',now()),

-- tiles / floor-tile
('p-tile-1','Плитка напольная керамогранит 600×600','floor-tile-600x600','tiles','floor-tile','KERAMA MARAZZI','керамогранит',1190,1190,true,0,'/img/products/tile-floor.png',now()),
('p-tile-floor-1','Керамогранит «Милан Брера» 30x30 см','keramogranit-milan-brera-30x30','tiles','floor-tile','KERAMA MARAZZI','керамогранит',1090,1340,true,-18,'/img/products/tile-milan-brera.png',now()),
('p-tile-floor-2','Керамогранит «Сланец» 30x30','keramogranit-slanec-30x30','tiles','floor-tile','KERAMA MARAZZI','керамогранит',1200,1200,true,0,'/img/products/tile-slanec.png',now()),
('p-tile-floor-3','Керамогранит EZ01 40x40','keramogranit-ez01-40x40','tiles','floor-tile','ESTIMA','керамогранит',1620,1620,true,0,'/img/products/tile-ez01.png',now()),

-- mixes / cement
('p-mix-1','Цемент М500 Д0 50 кг','cement-m500-50kg','mixes','cement','EUROCEMENT','цемент',520,520,true,0,'/img/products/cement.png',now()),
('p-mix-cement-1','Портландцемент Holcim M500 ЦЕМ II/А-И 42.5 25 кг','holcim-m500-25kg','mixes','cement','HOLCIM','цемент',183,220,true,-16,'/img/products/cement-holcim-25.png',now()),
('p-mix-cement-2','Портландцемент Евроцемент М500 ЦЕМ II/А-Ш 42.5 Н 50 кг','eurocement-m500-50kg','mixes','cement','ЕВРОЦЕМЕНТ','цемент',413,413,true,0,'/img/products/cement-eurocement-50.png',now()),
('p-mix-cement-3','Цемент Axton 5 кг','axton-cement-5kg','mixes','cement','AXTON','цемент',96,96,true,0,'/img/products/cement-axton-5.png',now()),
('p-mix-cement-4','Цемент монтажный Ceresit CX5 водоостанавливающий, 2 кг','ceresit-cx5-2kg','mixes','cement','CERESIT','цемент',330,370,true,-10,'/img/products/cement-ceresit-cx5.png',now()),
('p-mix-cement-5','Цемент Севряковцемент ПЦ-500 Д20, 50 кг','severyakovcement-pc500-50kg','mixes','cement','СЕВРЯКОВЦЕМЕНТ','цемент',305,305,true,0,'/img/products/cement-sevryakov-50.png',now()),

-- mixes / putty
('p-mix-putty-1','Шпаклевка полимерная финишная Axton 5 кг','axton-putty-finish-5kg','mixes','putty','AXTON','шпаклевка',160,160,true,0,'/img/products/putty-axton-5.png',now()),
('p-mix-putty-2','Шпаклёвка полимерная финишная Weber Vetonit LR Plus 20 кг','vetonit-lr-plus-20kg','mixes','putty','WEBER','шпаклевка',590,590,true,0,'/img/products/putty-vetonit-lr-plus.png',now()),
('p-mix-putty-3','Шпаклевка цементная Axton базовая, 25 кг','axton-putty-cement-base-25kg','mixes','putty','AXTON','шпаклевка',286,286,true,0,'/img/products/putty-axton-cement.png',now()),
('p-mix-putty-4','Шпаклёвка цементная финишная Weber Vetonit Facade white 20 кг','vetonit-facade-white-20kg','mixes','putty','WEBER','шпаклевка',560,670,true,-16,'/img/products/putty-vetonit-facade.png',now()),

-- mixes / plaster
('p-mix-plaster-1','Штукатурка гипсовая Plitonit GT 30 кг','plitonit-gt-30kg','mixes','plaster','PLITONIT','штукатурка',420,420,true,0,'/img/products/plaster-plitonit.png',now()),

-- power-tools / screwdriver
('p-power-1','Шуруповерт аккумуляторный BOSCH GSR 120-LI','bosch-gsr-120','power-tools','screwdriver','BOSCH','шуруповерт',6900,6900,true,0,'/img/products/screwdriver-bosch.png',now()),
('pt-6','Дрель-шуруповерт ELITECH ДА 10.8СЛК2 (2 АКБ, без ЗУ)','elitech-da-10-8','power-tools','screwdriver','ELITECH','шуруповерт аккумуляторный',2990,2990,true,0,'/img/products/power/elitech-10-8.png',now()),
('pt-7','Дрель-шуруповерт аккумуляторная 12 В (1x АКБ)','screwdriver-12v-1akb','power-tools','screwdriver','NoName','шуруповерт аккумуляторный',946,946,true,0,'/img/products/power/screwdriver-12v.png',now()),
('pt-8','Дрель-шуруповерт аккумуляторная Makita DDF453SYX5 (2x АКБ)','makita-ddf453syx5','power-tools','screwdriver','Makita','шуруповерт аккумуляторный',8990,9990,true,-10,'/img/products/power/makita-ddf453.png',now()),

-- power-tools / drill
('pt-1','Аккумуляторная дрель-шуруповерт BOSCH EasyDrill 12-2 (без АКБ)','bosch-easydrill-12-2','power-tools','drill','BOSCH','дрель-шуруповерт',2290,2290,true,0,'/img/products/power/bosch-easydrill.png',now()),
('pt-2','Дрель аккумуляторная Einhell TC-CD 18/35 Li 18 В','einhell-tc-cd-18-35','power-tools','drill','Einhell','дрель аккумуляторная',5990,7990,true,-25,'/img/products/power/einhell-18-35.png',now()),
('pt-3','Дрель ударная Bosch Professional GSB 13 RE (БЗП)','bosch-gsb-13-re','power-tools','drill','BOSCH','дрель ударная',4100,4100,true,0,'/img/products/power/bosch-gsb13.png',now()),

-- power-tools / heat-gun
('pt-4','Фен технический BOSCH UniversalHeat 600','bosch-universalheat-600','power-tools','heat-gun','BOSCH','фен технический',3990,3990,true,0,'/img/products/power/bosch-heatgun.png',now()),
('pt-5','Фен технический Makita HG5030K','makita-hg5030k','power-tools','heat-gun','Makita','фен технический',5290,5290,true,0,'/img/products/power/makita-heatgun.png',now());
-- ВАЖНО: тут точка с запятой!

-- =========================
-- BADGE LINKS (product_badge_links)
-- =========================

-- p-1: new, sale
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-1', id, now() FROM product_badges WHERE code IN ('new','sale');

-- p-3: hit, new
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-3', id, now() FROM product_badges WHERE code IN ('hit','new');

-- p-5: new, sale
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-5', id, now() FROM product_badges WHERE code IN ('new','sale');

-- p-6: hit
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-6', id, now() FROM product_badges WHERE code IN ('hit');

-- p-7: new
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-7', id, now() FROM product_badges WHERE code IN ('new');

-- p-8: hit, sale, recommended
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-8', id, now() FROM product_badges WHERE code IN ('hit','sale','recommended');

-- p-9: new
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-9', id, now() FROM product_badges WHERE code IN ('new');

-- p-floor-1: hit
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-floor-1', id, now() FROM product_badges WHERE code IN ('hit');

-- p-floor-2: sale
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-floor-2', id, now() FROM product_badges WHERE code IN ('sale');

-- p-floor-3: new, recommended
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-floor-3', id, now() FROM product_badges WHERE code IN ('new','recommended');

-- p-floor-5: sale
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-floor-5', id, now() FROM product_badges WHERE code IN ('sale');

-- p-floor-6: hit
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-floor-6', id, now() FROM product_badges WHERE code IN ('hit');

-- p-floor-7: recommended
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-floor-7', id, now() FROM product_badges WHERE code IN ('recommended');

-- p-tile-1: hit
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-tile-1', id, now() FROM product_badges WHERE code IN ('hit');

-- p-power-1: hit, recommended
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-power-1', id, now() FROM product_badges WHERE code IN ('hit','recommended');

-- p-tile-floor-1: hit, sale
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-tile-floor-1', id, now() FROM product_badges WHERE code IN ('hit','sale');

-- p-tile-floor-2: hit
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-tile-floor-2', id, now() FROM product_badges WHERE code IN ('hit');

-- p-tile-floor-3: new
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-tile-floor-3', id, now() FROM product_badges WHERE code IN ('new');

-- p-mix-cement-1: hit, sale
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-mix-cement-1', id, now() FROM product_badges WHERE code IN ('hit','sale');

-- p-mix-cement-2: hit
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-mix-cement-2', id, now() FROM product_badges WHERE code IN ('hit');

-- p-mix-cement-3: new
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-mix-cement-3', id, now() FROM product_badges WHERE code IN ('new');

-- p-mix-cement-4: sale
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-mix-cement-4', id, now() FROM product_badges WHERE code IN ('sale');

-- p-mix-cement-5: recommended
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-mix-cement-5', id, now() FROM product_badges WHERE code IN ('recommended');

-- p-mix-putty-2: hit, new
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-mix-putty-2', id, now() FROM product_badges WHERE code IN ('hit','new');

-- p-mix-putty-4: hit, new, sale
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-mix-putty-4', id, now() FROM product_badges WHERE code IN ('hit','new','sale');

-- p-mix-plaster-1: hit
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'p-mix-plaster-1', id, now() FROM product_badges WHERE code IN ('hit');

-- pt-1: hit, recommended
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'pt-1', id, now() FROM product_badges WHERE code IN ('hit','recommended');

-- pt-2: hit, sale, recommended
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'pt-2', id, now() FROM product_badges WHERE code IN ('hit','sale','recommended');

-- pt-3: hit, recommended
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'pt-3', id, now() FROM product_badges WHERE code IN ('hit','recommended');

-- pt-4: new, recommended
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'pt-4', id, now() FROM product_badges WHERE code IN ('new','recommended');

-- pt-5: hit
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'pt-5', id, now() FROM product_badges WHERE code IN ('hit');

-- pt-6: hit
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'pt-6', id, now() FROM product_badges WHERE code IN ('hit');

-- pt-7: new
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'pt-7', id, now() FROM product_badges WHERE code IN ('new');

-- pt-8: new, sale
INSERT INTO product_badge_links (product_id, badge_id, created_at)
SELECT 'pt-8', id, now() FROM product_badges WHERE code IN ('new','sale');
