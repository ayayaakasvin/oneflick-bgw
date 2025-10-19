SELECT 
    e.event_uuid,
    e.title,
    e.description,
    e.starting_time,
    e.ending_time,
    e.status,
    e.capacity,
    e.image_url,
    c.name AS category_name,
    u.username AS organizer_username,
    COALESCE(SUM(t.sold), 0) AS total_tickets_sold,
    ROW_NUMBER() OVER (ORDER BY COALESCE(SUM(t.sold), 0) DESC) AS rank
FROM events e
LEFT JOIN tickets t ON e.event_uuid = t.event_uuid
LEFT JOIN category c ON e.category_id = c.category_id
LEFT JOIN users u ON e.organizer_id = u.user_id
WHERE e.status = 'active'
GROUP BY e.event_uuid, e.title, e.description, e.starting_time, e.ending_time, 
         e.status, e.capacity, e.image_url, c.name, u.username
HAVING COALESCE(SUM(t.sold), 0) > 0
ORDER BY total_tickets_sold DESC
LIMIT 10;