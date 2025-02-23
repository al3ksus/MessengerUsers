CREATE TABLE IF NOT EXISTS users
(
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL,
    is_active BOOLEAN DEFAULT TRUE
);


CREATE TABLE IF NOT EXISTS friend_requests
(
    id BIGSERIAL PRIMARY KEY,
    sender_id BIGINT NOT NULL,
    recipient_id BIGINT NOT NULL,
    accepted BOOLEAN DEFAULT FALSE
);



ALTER TABLE friend_requests DROP CONSTRAINT IF EXISTS fk_sender_id_friend_requests;

ALTER TABLE friend_requests DROP CONSTRAINT IF EXISTS fk_recipient_id_friend_requests;

ALTER TABLE friend_requests DROP CONSTRAINT IF EXISTS unequal_sender_id_recipient_id_friend_requests;



ALTER TABLE friend_requests
    ADD CONSTRAINT fk_sender_id_friend_requests FOREIGN KEY (sender_id) REFERENCES users (id);

ALTER TABLE friend_requests
    ADD CONSTRAINT fk_recipient_id_friend_requests FOREIGN KEY (recipient_id) REFERENCES users (id);

ALTER TABLE friend_requests
    ADD CONSTRAINT unequal_sender_id_recipient_id_friend_requests CHECK (sender_id <> recipient_id);



CREATE UNIQUE INDEX IF NOT EXISTS friend_requests_sender_id_recipient_id_key ON friend_requests
    (LEAST(sender_id, recipient_id), GREATEST(sender_id, recipient_id));

CREATE INDEX IF NOT EXISTS sender_id_friend_requests ON friend_requests (sender_id);

CREATE INDEX IF NOT EXISTS recipient_id_friend_requests ON friend_requests (recipient_id);