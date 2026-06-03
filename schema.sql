  CREATE TABLE IF NOT EXISTS events (
      uuid TEXT PRIMARY KEY,
      summary TEXT NOT NULL,
      location TEXT,
      description TEXT NOT NULL,
      timezone TEXT NOT NULL,
      all_day BOOLEAN NOT NULL,
      start TEXT NOT NULL,
      end TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

