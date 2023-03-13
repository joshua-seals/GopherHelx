-- Version: 1.1
-- Description: Create table users
CREATE TABLE users (
    user_id UUID,
    name TEXT,
    session TEXT UNIQUE,
    password_hash TEXT,
    date_created TIMESTAMP,
    date_updated TIMESTAMP,

    PRIMARY KEY (user_id)
);

-- Version: 1.2
-- Description: Create table applications
CREATE TABLE applications (
	app_id   UUID,
	name     TEXT,
    image    TEXT,
	port     INT,

	PRIMARY KEY (app_id)
);

-- Version: 1.3
-- Description: Create table dashboard
CREATE TABLE dashboard (
	dash_id      UUID,
	user_id      UUID,
	app_id   UUID,

	PRIMARY KEY (dash_id),
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
	FOREIGN KEY (app_id) REFERENCES application (app_id) ON DELETE CASCADE
);