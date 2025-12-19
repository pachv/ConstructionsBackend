CREATE TABLE IF NOT EXISTS admin_email_settings (
    id         VARCHAR(255),
    email      VARCHAR(255),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

INSERT INTO admin_email_settings (id, email)
SELECT 'singleton', ''
WHERE NOT EXISTS (
    SELECT 1
    FROM admin_email_settings
    WHERE id = 'singleton'
);
