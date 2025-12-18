CREATE TABLE certificates (
    id         VARCHAR(255) ,
    title      VARCHAR(255) ,
    file_path  VARCHAR(1024), 
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


INSERT INTO certificates (
    id,
    title,
    file_path,
    created_at
) VALUES (
    'cert-main',
    'Основной сертификат',
    '/certificates/file/main.pdf',
    now()
);