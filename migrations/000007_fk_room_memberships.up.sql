ALTER TABLE room_memberships
ADD CONSTRAINT fk_room
FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE;
