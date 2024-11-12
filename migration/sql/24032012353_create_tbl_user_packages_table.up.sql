BEGIN;

    CREATE TABLE tbl_user_packages(
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES tbl_users(id),
        package_id INTEGER NOT NULL REFERENCES tbl_packages(id),
        valid_date timestamp NOT NULL,
        created_at timestamp NOT NULL DEFAULT (now()),
        updated_at timestamp DEFAULT (now()),
        unique(user_id,package_id)
    );

COMMIT;