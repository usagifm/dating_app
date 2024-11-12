
BEGIN;

    CREATE TABLE tbl_user_swipes(
        id SERIAL PRIMARY KEY,
        swiper_id INTEGER NOT NULL REFERENCES tbl_users(id),
        swiped_id INTEGER NOT NULL REFERENCES tbl_users(id),
        swipe_type VARCHAR(4) NOT NULL,
        created_at timestamp NOT NULL DEFAULT (now()),
        updated_at timestamp DEFAULT (now())
    );

COMMIT;