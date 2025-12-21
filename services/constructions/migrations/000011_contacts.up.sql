CREATE TABLE IF NOT EXISTS contacts_email_settings (
    id         VARCHAR(255),
    email      VARCHAR(255),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- singleton email
INSERT INTO contacts_email_settings (id, email)
SELECT 'singleton', 'info@example.com'
WHERE NOT EXISTS (
    SELECT 1 FROM contacts_email_settings WHERE id = 'singleton'
);



CREATE TABLE IF NOT EXISTS contacts_numbers (
    id         VARCHAR(255),
    phone      VARCHAR(255),
    label      VARCHAR(255),
    sort_order INTEGER,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_contacts_numbers_sort
    ON contacts_numbers(sort_order);

-- default phone numbers
INSERT INTO contacts_numbers (id, phone, label, sort_order)
SELECT 'phone-main', '+7 (900) 000-00-00', 'Основной', 0
WHERE NOT EXISTS (
    SELECT 1 FROM contacts_numbers WHERE id = 'phone-main'
);

INSERT INTO contacts_numbers (id, phone, label, sort_order)
SELECT 'phone-sales', '+7 (900) 111-11-11', 'Отдел продаж', 1
WHERE NOT EXISTS (
    SELECT 1 FROM contacts_numbers WHERE id = 'phone-sales'
);



CREATE TABLE IF NOT EXISTS contacts_addresses (
    id         VARCHAR(255),
    title      VARCHAR(255),
    address    VARCHAR(1024),
    lat        DOUBLE PRECISION,
    lon        DOUBLE PRECISION,
    sort_order INTEGER,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_contacts_addresses_sort
    ON contacts_addresses(sort_order);

-- default addresses
INSERT INTO contacts_addresses (id, title, address, lat, lon, sort_order)
SELECT
    'addr-office',
    'Офис',
    'г. Москва, ул. Примерная, д. 1',
    55.7558,
    37.6173,
    0
WHERE NOT EXISTS (
    SELECT 1 FROM contacts_addresses WHERE id = 'addr-office'
);

INSERT INTO contacts_addresses (id, title, address, lat, lon, sort_order)
SELECT
    'addr-warehouse',
    'Склад',
    'Московская область, Промзона',
    55.9000,
    37.5000,
    1
WHERE NOT EXISTS (
    SELECT 1 FROM contacts_addresses WHERE id = 'addr-warehouse'
);
