-- setup.sql

-- 1. User Table
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 2. Manga Table (With Indexing for fast search)
CREATE TABLE IF NOT EXISTS manga (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT,
    genres TEXT,
    status TEXT,
    total_chapters INTEGER,
    description TEXT
);

CREATE INDEX IF NOT EXISTS idx_manga_title ON manga(title); -- For <500ms search requirement

-- 3. Library/Progress Table (The "Net-centric" engine)
CREATE TABLE IF NOT EXISTS user_progress (
    user_id TEXT,
    manga_id TEXT,
    current_chapter INTEGER DEFAULT 0 CHECK (current_chapter >= 0),
    status TEXT DEFAULT 'Reading',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- CRITICAL for TCP Sync later
    PRIMARY KEY (user_id, manga_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (manga_id) REFERENCES manga(id) ON DELETE CASCADE
);