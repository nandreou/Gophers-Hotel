CREATE DATABaSE mydb;

CREATE TABLE users (
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    access_level INT NOT NULL
);

CREATE INDEX idx_users_email ON users(email); 
CREATE INDEX idx_users_access_level ON users(access_level);

CREATE TABLE rooms (
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    room_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE restrictions (
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    restriction_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE reservations (
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    room_id INT,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    start_date DATE,
    end_date DATE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_rooms
        FOREIGN KEY (room_id)
            REFERENCES rooms(id)
);

CREATE TABLE room_restrictions (
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    room_id INT,
    reservation_id INT,
    restriction_id INT,
    start_date DATE,
    end_date DATE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_rooms FOREIGN KEY (room_id) REFERENCES rooms(id),
    CONSTRAINT fk_reservation FOREIGN KEY (reservation_id) REFERENCES reservations(id),
    CONSTRAINT fk_restriction FOREIGN KEY (restriction_id) REFERENCES restrictions(id)
);

