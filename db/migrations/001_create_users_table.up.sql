-- Create users table with complete schema
-- 001_create_users_table.up.sql

CREATE TABLE users(
    id SERIAL NOT NULL,
    name varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    password varchar(255) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    phone varchar(20) NOT NULL DEFAULT ''::character varying,
    status varchar(20) DEFAULT 'active'::character varying,
    role varchar(50) DEFAULT 'user'::character varying,
    last_login_at timestamp without time zone,
    PRIMARY KEY(id)
);

-- Create indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_phone ON users(phone);

-- Add check constraints for data integrity
ALTER TABLE users ADD CONSTRAINT chk_users_status
    CHECK (status IN ('active', 'inactive', 'suspended', 'pending'));
