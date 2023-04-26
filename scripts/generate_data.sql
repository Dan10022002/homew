/*
INSERT INTO users(name, surname, age) VALUES ('Dan', 'Yasenko', 21),
                                             ('Nikita', 'Mishin', 35),
                                             ('Julia', 'Lebedeva', 13),
                                             ('Alexander', 'Pushkin', 28);

INSERT INTO tickets(user_id, cost, start_time) VALUES (2, 1500, '2020-02-03'),
                                                      (2, 2500, '2020-02-05'),
                                                      (3, 1000, '2020-09-25'),
                                                      (1, 777, '2020-04-04'),
                                                      (1, 888, '2020-05-05'),
                                                      (1, 999, '2020-06-06'),
                                                      (1, 555, '2020-07-07');
*/

INSERT INTO users(name, surname, age)
SELECT upper(substr(md5(random()::text), 0, 15)), upper(substr(md5(random()::text), 0, 25)), random() * 100
FROM generate_series(1, 50);

INSERT INTO tickets(user_id, cost, place)
SELECT 1 + random() * (SELECT id FROM users ORDER BY id DESC OFFSET 1 LIMIT 1), 1000 + random() * 5000, 1 + random() * 149
FROM generate_series(1, 70);