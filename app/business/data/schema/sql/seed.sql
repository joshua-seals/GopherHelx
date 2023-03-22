INSERT INTO users (user_id, user_name, session) VALUES
	(1, 'Bob','JufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2'),
	(2, 'Sally', 'UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW')
	ON CONFLICT DO NOTHING;

INSERT INTO applications (app_name, image, port) VALUES
	('jupyter-lab','jupyter/datascience-notebook', 8888),
	('pgadmin', 'pgadmin:latest', 8088),
	('postgresql', 'postgres', 5432),
	('r-studio','rstudio/rstudio-workbench',8787)
	ON CONFLICT DO NOTHING;

INSERT INTO dashboard (users_dash_id, apps_app_id, users_session) VALUES
	(1, 2, 'JufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2' ),
	(1, 3, 'JufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2'),
    (2, 1, 'UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW'), 
    (2, 2, 'UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW'),
    (2, 3, 'UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW')
	ON CONFLICT DO NOTHING;

-- On conflict do nothing means, if already there don't worry about it.