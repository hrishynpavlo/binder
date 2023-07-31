/*create database*/
CREATE DATABASE binder_all;

/*then switch to this database*/

CREATE TABLE IF NOT EXISTS users (id BIGSERIAL PRIMARY KEY,
    email VARCHAR(512) NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL,
    first_name VARCHAR(256) NOT NULL,
    last_name VARCHAR(256) NOT NULL,
    display_name VARCHAR(256),
    date_of_birth DATE NOT NULL,
    country VARCHAR(256) NOT NULL,
    geolocation POINT NOT NULL
);

CREATE TYPE interest AS ENUM ('Travel', 'Music', 'Books', 'Movies', 'Sport', 'Adventure', 'Pets', 'Animals', 'Food', 'Wine', 'Coffee', 'Drink', 'Walks', 'Hiking', 'Dancing', 'Gym', 'Tattoo' );

CREATE TABLE IF NOT EXISTS user_interests (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    interests interest[] NOT NULL,
    last_update TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION sp_create_user(
    email_param VARCHAR(512),
    password_hash_param VARCHAR,
    first_name_param VARCHAR(256),
    last_name_param VARCHAR(256),
    display_name_param VARCHAR(256),
    date_of_birth_param DATE,
    country_param VARCHAR(256),
    latitude_param NUMERIC,
    longitude_param NUMERIC
) RETURNS BIGINT
AS
$$
DECLARE
    new_user_id BIGINT;
BEGIN
    INSERT INTO users (email, password_hash, first_name, last_name, display_name, date_of_birth, country, geolocation)
    VALUES (email_param, password_hash_param, first_name_param, last_name_param, display_name_param, date_of_birth_param, country_param, POINT(latitude_param, longitude_param))
    RETURNING id INTO new_user_id;

    RETURN new_user_id;
END;
$$ 
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION sp_update_user_interests(
    user_id_param BIGINT,
    interests_param interest[]
) RETURNS BOOLEAN
AS
$$
BEGIN
    INSERT INTO user_interests (user_id, interests)
    VALUES (user_id_param, interests_param);
    RETURN true;
END;
$$
LANGUAGE plpgsql;

/*add fake data*/

DO
$$
DECLARE
    i INTEGER;
    email VARCHAR(512);
    password_hash VARCHAR;
    first_name VARCHAR(256);
    last_name VARCHAR(256);
    display_name VARCHAR(256);
    date_of_birth DATE;
    country VARCHAR(256);
    latitude NUMERIC;
    longitude NUMERIC;
BEGIN
    FOR i IN 1..100 LOOP
        email := 'user' || i || '@example.com';
        password_hash := md5(random()::text);
        first_name := (SELECT name FROM (VALUES ('Alice'), ('Emma'), ('Olivia'), ('Sophia'), ('Ava'), ('Isabella'), ('Mia'), ('Amelia'), ('Harper'), ('Evelyn')) AS female_names(name) OFFSET floor(random() * 10) LIMIT 1);
        last_name := (SELECT name FROM (VALUES ('Smith'), ('Johnson'), ('Brown'), ('Jones'), ('Miller'), ('Davis'), ('Garcia'), ('Rodriguez'), ('Martinez'), ('Hernandez')) AS surnames(name) OFFSET floor(random() * 10) LIMIT 1);
        display_name := 'Display Name ' || i;
        date_of_birth := DATE '1980-01-01' + (random() * 8000)::integer;
        country := 'USA';
        latitude := random() * 180 - 90;
        longitude := random() * 360 - 180;

        PERFORM sp_create_user(email, password_hash, first_name, last_name, display_name, date_of_birth, country, latitude, longitude);
    END LOOP;
END;
$$;

CREATE OR REPLACE VIEW users_info AS 
    SELECT u.id, u.email, u.first_name, u.last_name, u.display_name, u.date_of_birth, u.country, u.geolocation, ui.interests 
    FROM public.users u 
    LEFT JOIN public.user_interests ui ON ui.user_id = u.id;

CREATE USER binder_usr WITH PASSWORD 'binder_best_app';

GRANT ALL ON ALL TABLES IN SCHEMA public TO binder_usr;
GRANT ALL ON SCHEMA public TO binder_usr;