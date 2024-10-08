-- carefully execute this once to create the tables in the database

DROP TABLE IF EXISTS GAMEPLAYERS;
DROP TABLE IF EXISTS USERS;
DROP TABLE IF EXISTS DBSESSIONS;

CREATE TABLE USERS (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE GAMEPLAYERS (
    id SERIAL PRIMARY KEY,
    playername VARCHAR(255),
    avatar VARCHAR(255)
);

CREATE TABLE DBSESSIONS (
    id VARCHAR(255) PRIMARY KEY,
    user_id INTEGER REFERENCES USERS(id),
    validTo TIMESTAMP
);

INSERT INTO USERS (id, username, password) VALUES (1, 'test', 'test');