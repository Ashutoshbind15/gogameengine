-- carefully execute this once to create the tables in the database

DROP TABLE IF EXISTS GAMEPLAYERS;
DROP TABLE IF EXISTS USERS;

CREATE TABLE USERS (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL
);

CREATE TABLE GAMEPLAYERS (
    id SERIAL PRIMARY KEY,
    playername VARCHAR(255),
    avatar VARCHAR(255)
)