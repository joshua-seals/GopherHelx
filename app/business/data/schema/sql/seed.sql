INSERT INTO users (user_id, name, session) VALUES
	(1, 'Bob','$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a'),
	(2, 'Sally', '$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW')
	ON CONFLICT DO NOTHING;

INSERT INTO dashboard (dash_id, user_id, app) VALUES
	(1, 1, 2 ),
	(1, 1, 3),
    (2, 2, 1), 
    (2, 2, 2),
    (2, 2, 3)
	ON CONFLICT DO NOTHING;

INSERT INTO applications (app_id, name, image, port) VALUES
	(1, 'jupyter-lab','dockerhub.io/jupyter-lab:latest', 8080),
	(2, 'pgadmin', 'dockerhub.io/pgadmin:latest', 8088),
	(3, 'postgresql', 'dockerhub.io/postgresql:latest', 5432)
	ON CONFLICT DO NOTHING;

-- On conflict do nothing means, if already there don't worry about it.