CREATE TABLE IF NOT EXISTS messages (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      jid TEXT,
      name TEXT,
      channel_jid TEXT,
      message TEXT,
      type TEXT,
      command TEXT,
      timestamp DATE,
      is_group INTEGER,
      ignored INTEGER,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
      number TEXT PRIMARY KEY,
      jid TEXT,
      name TEXT,
      blocked INTEGER,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
