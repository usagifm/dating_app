BEGIN;

    CREATE TABLE tbl_packages(
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        description text NOT NULL,
        price NUMERIC NOT NULL,
        periode INTEGER NOT NULL,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_at timestamp NOT NULL DEFAULT (now()),
        updated_at timestamp DEFAULT (now())
    );

        INSERT INTO tbl_packages ( name, description, price, periode, is_active,created_at, updated_at)
    VALUES
        ("Pake Cinta 7 Hari", 'Dapatkan 7 Hari bebas akses unlimited', 15000, 7, TRUE,  NOW(), NOW());

COMMIT;