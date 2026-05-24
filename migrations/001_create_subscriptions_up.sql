CREATE TABLE subscriptions (
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR(255) NOT NULL,
    user_id   UUID           NOT NULL,
    price     INT            NOT NULL,
    start_date DATE NOT NULL,
    end_date   DATE
);
