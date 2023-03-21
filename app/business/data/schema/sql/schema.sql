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
CREATE TABLE dashboard (
	users_dash_id INTEGER ,
	users_session TEXT,
	apps_app_id  INTEGER,
	PRIMARY KEY (users_dash_id, users_session, apps_app_id)
);