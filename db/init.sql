/* create database */

/*CREATE DATABASE binder_all;*/

/* then switch to this database */
/* create tables, enums, views */

-- CREATE EXTENSION IF NOT EXISTS decoderbufs;

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
    user_id BIGINT NOT NULL UNIQUE,
    interests interest[] NOT NULL,
    last_update TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_photos(
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    photo_urls TEXT[] NOT NULL,
    last_update TIMESTAMP NOT NULL DEFAULT NOW(),
    primary_photo_index SMALLINT DEFAULT 0,
    CONSTRAINT fk_user_photo FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_filters(
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    min_distance_km SMALLINT NOT NULL CONSTRAINT positive_min_distance CHECK(min_distance_km > 0),
    max_distance_km SMALLINT NOT NULL CONSTRAINT positive_max_distance CHECK(max_distance_km > min_distance_km),
    min_age SMALLINT NOT NULL CONSTRAINT positive_min_age CHECK(min_age > 10),
    max_age SMALLINT NOT NULL CONSTRAINT positive_max_age CHECK(max_age > min_age AND max_age < 120),
    CONSTRAINT fk_user_filters FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS subscription_plans(
    id SERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL UNIQUE,
    description VARCHAR(512) NOT NULL,
    display_name VARCHAR(256) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_active BOOLEAN NOT NULL,
    refresh_period_in_hours SMALLINT,
    matching_limit SMALLINT
);

CREATE TABLE IF NOT EXISTS user_subscriptions(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL, 
    subscription_plan_id INT NOT NULL, 
    since TIMESTAMP NOT NULL DEFAULT NOW(),
    is_active BOOLEAN NOT NULL,
    start_period TIMESTAMP NOT NULL DEFAULT NOW(),
    end_period TIMESTAMP NOT NULL,
    CONSTRAINT fk_user_subscriptions_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_subscriptions_subscription_plan FOREIGN KEY(subscription_plan_id) REFERENCES subscription_plans(id)
);

CREATE TABLE IF NOT EXISTS user_geos(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    country_code VARCHAR(3) NOT NULL,
    state_code VARCHAR(128),
    city VARCHAR(128),
    geolocation POINT NOT NULL,
    last_modified TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user_geo FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_feeds(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    country_code VARCHAR(3) NOT NULL,
    state_code VARCHAR(128),
    city VARCHAR(128),
    feed JSON NOT NULL, 
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    feed_offset SMALLINT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL,
    CONSTRAINT fk_user_feed_users FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_country_code ON user_geos (country_code);
CREATE INDEX idx_state_code ON user_geos(state_code);
CREATE INDEX idx_city ON user_geos (city);

CREATE INDEX idx_feed_activity ON user_feeds(is_active);

CREATE OR REPLACE VIEW users_info AS 
    SELECT u.id, u.email, u.first_name, u.last_name, u.display_name, u.date_of_birth, u.country, u.geolocation, ui.interests, up.photo_urls, up.primary_photo_index, uf.min_distance_km, uf.max_distance_km, uf.min_age, uf.max_age
    FROM public.users u 
    LEFT JOIN public.user_interests ui ON ui.user_id = u.id
    LEFT JOIN public.user_photos up ON up.user_id = u.id
    LEFT JOIN public.user_filters uf ON uf.user_id = u.id;

CREATE OR REPLACE VIEW user_matching AS 
    SELECT u.id, u.geolocation[0] as latitude, u.geolocation[1] as longitude, date_part('years', age(u.date_of_birth)) as age, array_to_string(ui.interests, ',') as interests, uf.max_distance_km, uf.min_age, uf.max_age
    FROM users u
    LEFT JOIN public.user_interests ui ON ui.user_id = u.id
    LEFT JOIN public.user_filters uf ON uf.user_id = u.id
    WHERE ui.interests IS NOT NULL;

/* create sql functions to upsert data */

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
) RETURNS SETOF users_info
AS
$$
DECLARE
    new_user_id BIGINT;
BEGIN
    BEGIN
        INSERT INTO users (email, password_hash, first_name, last_name, display_name, date_of_birth, country, geolocation)
        VALUES (email_param, password_hash_param, first_name_param, last_name_param, display_name_param, date_of_birth_param, country_param, POINT(latitude_param, longitude_param))
        RETURNING id INTO new_user_id;

        INSERT INTO user_subscriptions (user_id, subscription_plan_id, is_active, end_period)
        VALUES (new_user_id, 1, true, NOW() + INTERVAL '1 months');
    EXCEPTION 
        WHEN OTHERS THEN
            RAISE;
    END;

    RETURN QUERY SELECT * FROM users_info WHERE id = new_user_id;
END;
$$ 
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION sp_update_user_interests(
    user_id_param BIGINT,
    interests_param interest[]
) RETURNS SETOF users_info
AS
$$
BEGIN
    INSERT INTO user_interests (user_id, interests)
    VALUES (user_id_param, interests_param)
    ON CONFLICT (user_id)
    DO
        UPDATE SET interests = interests_param, last_update = NOW();

    RETURN QUERY SELECT * FROM users_info WHERE id = user_id_param;
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION sp_update_user_photos(
    user_id_param BIGINT,
    photo_urls_param TEXT[]
) RETURNS SETOF users_info
AS
$$
BEGIN
    INSERT INTO user_photos (user_id, photo_urls)
    VALUES (user_id_param, photo_urls_param)
    ON CONFLICT (user_id)
    DO
        UPDATE SET photo_urls = photo_urls_param, last_update = NOW();

    RETURN QUERY SELECT * FROM users_info WHERE id = user_id_param;
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION sp_update_user_filters(
    user_id_param BIGINT,
    min_distance_km_param SMALLINT,
    max_distance_km_param SMALLINT,
    min_age_param SMALLINT,
    max_age_param SMALLINT
) RETURNS SETOF users_info
AS
$$
BEGIN
    INSERT INTO user_filters(user_id, min_distance_km, max_distance_km, min_age, max_age)
    VALUES (user_id_param, min_distance_km_param, max_distance_km_param, min_age_param, max_age_param)
    ON CONFLICT (user_id)
    DO
        UPDATE SET min_distance_km = min_distance_km_param, max_distance_km = max_distance_km_param, min_age = min_age_param, max_age = max_age_param;

    RETURN QUERY SELECT * FROM users_info WHERE id = user_id_param;
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION sp_update_user_geo(
    user_id_param BIGINT,
    country_code_param VARCHAR(3),
    state_code_param VARCHAR(128),
    city_param VARCHAR(128),
    latitude_param NUMERIC,
    longitude_param NUMERIC
) RETURNS VOID
AS
$$
BEGIN
    INSERT INTO user_geos(user_id, country_code, state_code, city, geolocation)
    VALUES(user_id_param, country_code_param, state_code_param, city_param, POINT(latitude_param, longitude_param))
    ON CONFLICT (user_id)
    DO
        UPDATE SET geolocation = POINT(latitude_param, longitude_param);
END;
$$
LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION sp_find_users_to_match(
    user_id_param BIGINT
) RETURNS TABLE (
    id BIGINT,
    email VARCHAR(256),
    first_name VARCHAR(256),
    last_name VARCHAR(256),
    display_name VARCHAR(256),
    date_of_birth DATE,
    country_code VARCHAR(3),
    state_code VARCHAR(128),
    city VARCHAR(128),
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    interests INTEREST[],
    photo_urls TEXT[],
    primary_photo_index SMALLINT,
    min_distance_km SMALLINT,
    max_distance_km SMALLINT,
    min_age SMALLINT,
    max_age SMALLINT
)
AS
$$
BEGIN
    RETURN QUERY
    WITH user_location AS (SELECT ug.country_code, ug.state_code, ug.city FROM user_geos ug WHERE ug.user_id = user_id_param),
         feed_users AS (SELECT ug.user_id, ug.geolocation as actual_geolocation, ug.country_code, ug.state_code, ug.city FROM user_geos ug WHERE ug.country_code IN (SELECT ul.country_code FROM user_location ul) AND ug.state_code IN (SELECT ul.state_code FROM user_location ul) AND ug.city IN (SELECT ul.city FROM user_location ul))
    SELECT ui.id, ui.email, ui.first_name, ui.last_name, ui.display_name, ui.date_of_birth, fu.country_code, fu.state_code, fu.city, fu.actual_geolocation[0] AS latitude, fu.actual_geolocation[1] AS longitude, ui.interests, ui.photo_urls, ui.primary_photo_index, ui.min_distance_km, ui.max_distance_km, ui.min_age, ui.max_age 
    FROM users_info ui 
    JOIN feed_users fu ON fu.user_id = ui.id;
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION sp_create_feed_snapshot(
    user_id_param BIGINT,
    country_code_param VARCHAR(3),
    state_code_param VARCHAR(128),
    city_param VARCHAR(128),
    feed_param JSON
) RETURNS VOID
AS
$$
BEGIN
    INSERT INTO user_feeds (user_id, country_code, state_code, city, feed, is_active)
    VALUES (user_id_param, country_code_param, state_code_param, city_param, feed_param, TRUE);
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION sp_login_user(
    email_param VARCHAR(512)
) RETURNS TABLE (id BIGINT, password_hash VARCHAR)
AS
$$
BEGIN
    RETURN QUERY 
    SELECT u.id, u.password_hash FROM users u WHERE u.email = email_param;
END;
$$
LANGUAGE plpgsql;

/*add fake data*/

INSERT INTO subscription_plans (name, description, display_name, is_active, refresh_period_in_hours, matching_limit)
VALUES ('BasePlan', 'Base subscription plan for newcomers', 'Base plan', true, 12, 100);

/*
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
    FOR i IN 1..500 LOOP
        email := 'user' || i || '@example.com';
        password_hash := md5(random()::text);
        first_name := (SELECT name FROM (VALUES ('Alice'), ('Emma'), ('Olivia'), ('Sophia'), ('Ava'), ('Isabella'), ('Mia'), ('Amelia'), ('Harper'), ('Evelyn')) AS female_names(name) OFFSET floor(random() * 10) LIMIT 1);
        last_name := (SELECT name FROM (VALUES ('Smith'), ('Johnson'), ('Brown'), ('Jones'), ('Miller'), ('Davis'), ('Garcia'), ('Rodriguez'), ('Martinez'), ('Hernandez')) AS surnames(name) OFFSET floor(random() * 10) LIMIT 1);
        display_name := first_name || last_name;
        date_of_birth := DATE '1980-01-01' + (random() * 8000)::integer;
        country := 'US';
        latitude := random() * 180 - 90;
        longitude := random() * 360 - 180;

        PERFORM sp_create_user(email, password_hash, first_name, last_name, display_name, date_of_birth, country, latitude, longitude);
    END LOOP;
END;
$$;
*/

CREATE USER binder_app WITH PASSWORD 'binder_best_app';

GRANT SELECT, UPDATE, INSERT, REFERENCES ON ALL TABLES IN SCHEMA public TO binder_app;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO binder_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO binder_app; 