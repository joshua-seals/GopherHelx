-- Version: 1.1
-- Description: Create table users
-- Ideally we use UUID where SERIAL exists but 
-- currently it is not implemented.
CREATE TABLE users (
    user_id INTEGER,
    user_name TEXT,
    session TEXT UNIQUE,

    PRIMARY KEY (user_id)
);

-- Version: 1.2
-- Description: Create table applications
CREATE TABLE applications (
	app_id   INTEGER,
	app_name TEXT,
    image    TEXT,
	port     TEXT,

	PRIMARY KEY (app_id)
);

-- Version: 1.3
-- Description: Create table dashboard
CREATE TABLE dashboard (
	dash_id      INTEGER,
	user_id      INTEGER,
	app_id  	 INTEGER,
	user_session TEXT,

	PRIMARY KEY (dash_id),
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
	FOREIGN KEY (app_id) REFERENCES applications (app_id) ON DELETE CASCADE,
	FOREIGN KEY (user_session) REFERENCES users (session) ON DELETE CASCADE
);