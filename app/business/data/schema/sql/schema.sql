-- Version: 1.1
-- Description: Create table users
-- Ideally we use UUID where SERIAL exists but 
-- currently it is not implemented.
CREATE TABLE users (
    user_id SERIAL,
    user_name TEXT,
    session TEXT UNIQUE,

    PRIMARY KEY (user_id)
);

-- Version: 1.2
-- Description: Create table applications
CREATE TABLE applications (
	app_id   SERIAL,
	app_name TEXT,
    image    TEXT,
	port     INTEGER,

	PRIMARY KEY (app_id)
);

-- Version: 1.3
-- Description: Create table dashboard
-- dash_id should be created based on session and user id
CREATE TABLE dashboard (
	dash_id      SERIAL,
	user_id      INTEGER,
	app_id  	 INTEGER,
	user_session TEXT,

	PRIMARY KEY (dash_id),
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
	FOREIGN KEY (app_id) REFERENCES applications (app_id) ON DELETE CASCADE,
	FOREIGN KEY (user_session) REFERENCES users (session) ON DELETE CASCADE
);