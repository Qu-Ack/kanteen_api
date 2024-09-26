-- +goose Up
create table otp(id UUID PRIMARY KEY, mobile TEXT NOT NULL, OTP TEXT NOT NULL);


-- +goose Down
drop table otp;