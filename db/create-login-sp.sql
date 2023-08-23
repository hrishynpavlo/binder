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

GRANT SELECT, UPDATE, INSERT, REFERENCES ON ALL TABLES IN SCHEMA public TO binder_app;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO binder_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO binder_app; 