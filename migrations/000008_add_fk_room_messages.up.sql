ALTER TABLE room_messages
ADD CONSTRAINT fk_room_messages_room
FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE;
