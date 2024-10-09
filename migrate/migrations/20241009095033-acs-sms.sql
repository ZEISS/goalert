-- +migrate Up
CREATE SEQUENCE IF NOT EXISTS acs_sms_callbacks_id_seq;

-- Table Definition
CREATE TABLE acs_sms_callbacks (
    "phone_number" text NOT NULL,
    "callback_id" uuid NOT NULL,
    "code" int4 NOT NULL,
    "id" int8 NOT NULL DEFAULT nextval('acs_sms_callbacks_id_seq'::regclass),
    "sent_at" timestamptz NOT NULL DEFAULT now(),
    "alert_id" int8,
    "service_id" uuid,
    CONSTRAINT "acs_sms_callbacks_alert_id_fkey" FOREIGN KEY ("alert_id") REFERENCES alerts("id") ON DELETE CASCADE,
    CONSTRAINT "acs_sms_callbacks_service_id_fkey" FOREIGN KEY ("service_id") REFERENCES services("id") ON DELETE CASCADE
);

CREATE SEQUENCE IF NOT EXISTS acs_sms_errors_id_seq;

CREATE TABLE acs_sms_errors (
    "phone_number" text NOT NULL,
    "error_message" text NOT NULL,
    "outgoing" bool NOT NULL,
    "occurred_at" timestamptz NOT NULL DEFAULT now(),
    "id" int8 NOT NULL DEFAULT nextval('acs_sms_errors_id_seq'::regclass)
);

-- +migrate Down
DROP TABLE acs_sms_callbacks;

DROP TABLE acs_sms_errors;