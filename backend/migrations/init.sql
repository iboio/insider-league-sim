CREATE TABLE league
(
    leagueId  CHAR(36) PRIMARY KEY UNIQUE NOT NULL,
    name      VARCHAR(255)                NOT NULL,
    teamCount INT                         NOT NULL,
    isActive  BOOLEAN  DEFAULT FALSE,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS teams
(
    id           INT AUTO_INCREMENT PRIMARY KEY,
    leagueId     CHAR(36) NOT NULL,
    teamName     CHAR(32) NOT NULL,
    attackPower  FLOAT    NOT NULL,
    defensePower FLOAT    NOT NULL,
    morale       FLOAT    NOT NULL,
    stamina      FLOAT    NOT NULL,
    createdAt    DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (leagueId) REFERENCES league (leagueId)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE TABLE standings
(
    id        INT AUTO_INCREMENT PRIMARY KEY,
    leagueId  CHAR(36) NOT NULL,
    teamName  CHAR(32) NOT NULL,
    goals     INT      NOT NULL,
    against   INT      NOT NULL,
    diff      INT GENERATED ALWAYS AS (goals - against) STORED,
    played    INT      NOT NULL,
    wins      INT      NOT NULL,
    losses    INT      NOT NULL,
    draws     INT      NOT NULL,
    points    INT      NOT NULL,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (leagueId) REFERENCES league (leagueId)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS matches
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    leagueId   CHAR(36)    NOT NULL,
    home       VARCHAR(32) NOT NULL,
    homeGoals  INT         NOT NULL,
    away       VARCHAR(32) NOT NULL,
    awayGoals  INT         NOT NULL,
    winnerName VARCHAR(32) NOT NULL,
    matchWeek  INT         NOT NULL,
    isPlayed   BOOLEAN   DEFAULT FALSE,
    createdAt  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (leagueId) REFERENCES league (leagueId)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);