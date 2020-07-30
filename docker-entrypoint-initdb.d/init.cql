CREATE KEYSPACE public
    WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

USE public;

CREATE TABLE message (
    id UUID PRIMARY KEY,
    email        text,
    title        text,
    content      text,
    magic_number int,
);

CREATE INDEX ON message(email);
CREATE INDEX ON message(magic_number);