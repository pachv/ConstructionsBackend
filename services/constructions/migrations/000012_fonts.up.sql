CREATE TABLE IF NOT EXISTS admin_fonts (
  id         TEXT PRIMARY KEY,
  name       TEXT NOT NULL,
  file_path  TEXT NOT NULL,
  selected   BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Гарантия: selected=true может быть только у одной строки
CREATE UNIQUE INDEX IF NOT EXISTS ux_admin_fonts_selected_one
ON admin_fonts ((selected))
WHERE selected = TRUE;

-- Ровно 1 стартовый шрифт (файл должен лежать в build: /app/build/static/media/...)
INSERT INTO admin_fonts (id, name, file_path, selected, created_at, updated_at)
VALUES (
  'default',
  'Стандартный',
  'static/media/AppFont.31d6cfe0d16ae931b73c.woff2',
  TRUE,
  now(),
  now()
)
ON CONFLICT (id) DO NOTHING;
