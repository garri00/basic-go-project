CREATE TABLE IF NOT EXISTS service.accounts
(
    id        UUID PRIMARY KEY DEFAULT service.gen_random_uuid(),
    login     VARCHAR(255) NOT NULL,
    password  VARCHAR(255) NOT NULL,
    is_active BOOLEAN      NOT NULL
);
ALTER TABLE service.accounts
    OWNER TO "serviceadmin";

