-- Create database tables for healthcare system

-- Physicians table
CREATE TABLE physicians (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    license_number VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Patients table
CREATE TABLE patients (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Drugs table
CREATE TABLE drugs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    generic_name VARCHAR(255),
    strength VARCHAR(100),
    dosage_form VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Patient-Physician relationships (many-to-many)
CREATE TABLE patient_physicians (
    id SERIAL PRIMARY KEY,
    patient_id INTEGER NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    physician_id INTEGER NOT NULL REFERENCES physicians(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(patient_id, physician_id)
);

-- Prescriptions table
CREATE TABLE prescriptions (
    id SERIAL PRIMARY KEY,
    patient_id INTEGER NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    physician_id INTEGER NOT NULL REFERENCES physicians(id) ON DELETE CASCADE,
    drug_id INTEGER NOT NULL REFERENCES drugs(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    sig TEXT NOT NULL, -- Instructions for use
    prescribed_date DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_prescriptions_patient_id ON prescriptions(patient_id);
CREATE INDEX idx_prescriptions_physician_id ON prescriptions(physician_id);
CREATE INDEX idx_prescriptions_drug_id ON prescriptions(drug_id);
CREATE INDEX idx_prescriptions_prescribed_date ON prescriptions(prescribed_date);
CREATE INDEX idx_patient_physicians_patient_id ON patient_physicians(patient_id);
CREATE INDEX idx_patient_physicians_physician_id ON patient_physicians(physician_id);

-- Composite index for top drugs analytics query
CREATE INDEX idx_prescriptions_date_drug_id ON prescriptions(prescribed_date, drug_id);