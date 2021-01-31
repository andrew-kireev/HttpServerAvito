CREATE TABLE hotels
(
    id serial not null primary key,
    description text,
    cost int
);

CREATE TABLE bookings
(
    booking_id INT GENERATED ALWAYS AS IDENTITY,
    hotel_id INT,
    begin_data date,
    end_date date,
    PRIMARY KEY (booking_id),
    CONSTRAINT fk_hotel
        FOREIGN KEY (hotel_id)
            REFERENCES hotels(id)
);