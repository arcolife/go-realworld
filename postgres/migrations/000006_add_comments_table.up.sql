BEGIN;

CREATE TABLE IF NOT EXISTS comments (
    id serial primary key,
    body text not null,
    author_id int not null,
    article_id int not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint fk_author foreign key(author_id) references users(id) on delete cascade,
    constraint fk_article foreign key(article_id) references articles(id)
);

COMMIT;