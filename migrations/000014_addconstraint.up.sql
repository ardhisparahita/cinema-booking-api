ALTER TABLE payments
ADD CONSTRAINT uq_payments_booking_id UNIQUE (booking_id);