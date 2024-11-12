BEGIN;
    CREATE TABLE tbl_user_matches(
        id SERIAL PRIMARY KEY,
        user_id_1 INTEGER NOT NULL REFERENCES tbl_users(id),
        user_id_2 INTEGER NOT NULL REFERENCES tbl_users(id),
        created_at timestamp NOT NULL DEFAULT (now()),
        updated_at timestamp DEFAULT (now())
    );


COMMIT;