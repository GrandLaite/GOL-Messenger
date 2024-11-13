-- Таблица пользователей
CREATE TABLE Users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,  
    is_premium BOOLEAN NOT NULL DEFAULT FALSE,  
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица сообщений
CREATE TABLE Messages (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES Users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица лайков
CREATE TABLE Likes (
    id SERIAL PRIMARY KEY,
    message_id INTEGER REFERENCES Messages(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES Users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (message_id, user_id)  -- Гарантирует, что пользователь может поставить только один лайк на сообщение
);

-- Таблица суперлайков
CREATE TABLE SuperLikes (
    id SERIAL PRIMARY KEY,
    message_id INTEGER REFERENCES Messages(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES Users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (message_id, user_id)  -- Гарантирует, что пользователь может поставить только один суперлайк на сообщение
);
