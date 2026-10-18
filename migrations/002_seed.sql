INSERT INTO halls (id, name) VALUES (1, 'Hall 1');

-- 5 rows (A-E) x 8 seats = 40 seats
INSERT INTO seats (hall_id, row_label, number)
WITH RECURSIVE n(i) AS (SELECT 1 UNION ALL SELECT i + 1 FROM n WHERE i < 8)
SELECT 1, r.l, n.i
FROM (SELECT 'A' AS l UNION ALL SELECT 'B' UNION ALL SELECT 'C'
      UNION ALL SELECT 'D' UNION ALL SELECT 'E') r
CROSS JOIN n;

INSERT INTO movies (id, title, duration_min) VALUES (1, 'Interstellar', 169);

INSERT INTO showtimes (movie_id, hall_id, starts_at, price_cents)
VALUES (1, 1, CAST(strftime('%s', 'now', '+1 day') AS INTEGER), 1200);