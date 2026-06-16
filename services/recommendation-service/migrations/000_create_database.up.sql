CREATE DATABASE recommendation_service OWNER dietician_user;
GRANT ALL PRIVILEGES ON DATABASE recommendation_service TO dietician_user;
\c recommendation_service;
GRANT USAGE ON SCHEMA public TO dietician_user;
GRANT CREATE ON SCHEMA public TO dietician_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO dietician_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO dietician_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO dietician_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO dietician_user;
