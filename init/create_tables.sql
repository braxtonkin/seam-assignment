CREATE TABLE IF NOT EXISTS blogs (
    id serial PRIMARY KEY,
    title varchar(255) NOT NULL,
    author varchar(255) NOT NULL,
    content text NOT NULL
);

CREATE TABLE IF NOT EXISTS comments (
    id      serial PRIMARY KEY,
    blog_id INTEGER      NOT NULL,
    author  varchar(255) NOT NULL,
    content text         NOT NULL,
    FOREIGN KEY (blog_id) REFERENCES blogs (id) ON DELETE CASCADE
);