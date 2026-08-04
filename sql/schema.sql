CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL, 
    password VARCHAR(200) NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS boards (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title       VARCHAR(200) NOT NULL,
    content    TEXT NOT NULL,
    image_path VARCHAR(200) NOT NULL DEFAULT '',
    is_private BOOLEAN NOT NULL DEFAULT false ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_board_user_id ON boards(user_id);
CREATE INDEX IF NOT EXISTS idx_board_creat_at ON boards(created_at DESC);

CREATE TABLE IF NOT EXISTS board_shares (
    board_id INTEGER NOT NULL REFERENCES boards (id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    PRIMARY KEY (board_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_board_shares_user_id ON board_shares(user_id);

CREATE TABLE IF NOT EXISTS comments (
    id  SERIAL PRIMARY KEY,
    board_id INTEGER NOT NULL REFERENCES boards (id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS index_comments_board_id ON comments(board_id);