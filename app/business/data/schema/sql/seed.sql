INSERT INTO users (user_id, user_name, session) VALUES
	(1, 'Bob','JufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2'),
	(2, 'Sally', 'UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW')
	ON CONFLICT DO NOTHING;

INSERT INTO applications (app_id, app_name, image, port) VALUES
	(1, 'jupyter-lab','dockerhub.io/jupyter-lab:latest', 8080),
	(2, 'pgadmin', 'dockerhub.io/pgadmin:latest', 8088),
	(3, 'postgresql', 'dockerhub.io/postgresql:latest', 5432)
	ON CONFLICT DO NOTHING;

INSERT INTO dashboard (dash_id, user_id, app_id, user_session) VALUES
	(1, 1, 2, 'JufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2' ),
	(1, 1, 3, 'JufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2'),
    (2, 2, 1, 'UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW'), 
    (2, 2, 2, 'UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW'),
    (2, 2, 3, 'UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW')
	ON CONFLICT DO NOTHING;

-- On conflict do nothing means, if already there don't worry about it.