CREATE SCHEMA IF NOT EXISTS league_sim;
USE league_sim;

CREATE TABLE IF NOT EXISTS league
(
    id        INT AUTO_INCREMENT PRIMARY KEY,
    name      VARCHAR(255) NOT NULL,
    leagueId  CHAR(36)     NOT NULL UNIQUE,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS active_league
(
    id               INT AUTO_INCREMENT PRIMARY KEY,
    leagueId         CHAR(36) NOT NULL,
    upcomingFixtures JSON,
    teams            JSON,
    playedFixtures   JSON,
    currentWeek      INT,
    standings        JSON,
    onActiveLeague   BOOLEAN   DEFAULT FALSE,
    createdAt        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (leagueId) REFERENCES league (leagueId)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS match_results
(
    id             INT AUTO_INCREMENT PRIMARY KEY,
    leagueId       CHAR(36)    NOT NULL,
    homeTeam  VARCHAR(36) NOT NULL,
    homeGoals  INT         NOT NULL,
    awayTeam   VARCHAR(36) NOT NULL,
    awayGoals  INT         NOT NULL,
    winnerName VARCHAR(36) NOT NULL,
    matchWeek      INT         NOT NULL,
    createdAt      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (leagueId) REFERENCES league (leagueId)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);





