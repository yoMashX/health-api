-- Seed data for healthcare system

-- Insert physicians
INSERT INTO physicians (first_name, last_name, license_number, email) VALUES
('John', 'Smith', 'MD123456', 'j.smith@hospital.com'),
('Sarah', 'Johnson', 'MD234567', 's.johnson@clinic.com'),
('Michael', 'Brown', 'MD345678', 'm.brown@medical.com'),
('Emily', 'Davis', 'MD456789', 'e.davis@health.com'),
('Robert', 'Wilson', 'MD567890', 'r.wilson@care.com');

-- Insert patients
INSERT INTO patients (first_name, last_name, date_of_birth, email, phone) VALUES
('Alice', 'Anderson', '1985-03-15', 'alice.anderson@email.com', '555-0101'),
('Bob', 'Baker', '1990-07-22', 'bob.baker@email.com', '555-0102'),
('Carol', 'Clark', '1978-11-08', 'carol.clark@email.com', '555-0103'),
('David', 'Evans', '1992-01-30', 'david.evans@email.com', '555-0104'),
('Emma', 'Foster', '1988-05-12', 'emma.foster@email.com', '555-0105'),
('Frank', 'Garcia', '1975-09-25', 'frank.garcia@email.com', '555-0106'),
('Grace', 'Harris', '1993-12-03', 'grace.harris@email.com', '555-0107'),
('Henry', 'Irving', '1982-06-18', 'henry.irving@email.com', '555-0108');

-- Insert drugs
INSERT INTO drugs (name, generic_name, strength, dosage_form) VALUES
('Lipitor', 'Atorvastatin', '20mg', 'Tablet'),
('Metformin', 'Metformin', '500mg', 'Tablet'),
('Lisinopril', 'Lisinopril', '10mg', 'Tablet'),
('Amlodipine', 'Amlodipine', '5mg', 'Tablet'),
('Omeprazole', 'Omeprazole', '20mg', 'Capsule'),
('Simvastatin', 'Simvastatin', '40mg', 'Tablet'),
('Hydrochlorothiazide', 'Hydrochlorothiazide', '25mg', 'Tablet'),
('Atenolol', 'Atenolol', '50mg', 'Tablet'),
('Furosemide', 'Furosemide', '40mg', 'Tablet'),
('Warfarin', 'Warfarin', '5mg', 'Tablet');

-- Insert patient-physician relationships
INSERT INTO patient_physicians (patient_id, physician_id) VALUES
(1, 1), (1, 2), -- Alice sees Dr. Smith and Dr. Johnson
(2, 1), (2, 3), -- Bob sees Dr. Smith and Dr. Brown
(3, 2), (3, 4), -- Carol sees Dr. Johnson and Dr. Davis
(4, 3), (4, 5), -- David sees Dr. Brown and Dr. Wilson
(5, 1), (5, 4), -- Emma sees Dr. Smith and Dr. Davis
(6, 2), (6, 5), -- Frank sees Dr. Johnson and Dr. Wilson
(7, 3), (7, 1), -- Grace sees Dr. Brown and Dr. Smith
(8, 4), (8, 5); -- Henry sees Dr. Davis and Dr. Wilson

-- Insert prescriptions with dates spread over the last 6 months
INSERT INTO prescriptions (patient_id, physician_id, drug_id, quantity, sig, prescribed_date) VALUES
-- Recent prescriptions (last 30 days)
(1, 1, 1, 30, 'Take one tablet daily with food', '2024-09-15'),
(2, 1, 2, 60, 'Take one tablet twice daily with meals', '2024-09-20'),
(3, 2, 3, 30, 'Take one tablet daily in the morning', '2024-09-25'),
(4, 3, 4, 30, 'Take one tablet daily', '2024-09-28'),
(5, 1, 5, 30, 'Take one capsule daily before breakfast', '2024-10-01'),
(6, 2, 6, 30, 'Take one tablet daily at bedtime', '2024-10-05'),
(7, 3, 7, 30, 'Take one tablet daily in the morning', '2024-10-08'),
(8, 4, 8, 30, 'Take one tablet twice daily', '2024-10-10'),
(1, 2, 9, 30, 'Take one tablet daily with water', '2024-10-12'),
(2, 3, 10, 30, 'Take as directed by physician', '2024-10-14'),

-- Older prescriptions (2-6 months ago)
(3, 4, 1, 30, 'Take one tablet daily with food', '2024-08-15'),
(4, 5, 2, 60, 'Take one tablet twice daily with meals', '2024-08-20'),
(5, 4, 3, 30, 'Take one tablet daily in the morning', '2024-07-25'),
(6, 5, 4, 30, 'Take one tablet daily', '2024-07-28'),
(7, 1, 5, 30, 'Take one capsule daily before breakfast', '2024-07-01'),
(8, 5, 6, 30, 'Take one tablet daily at bedtime', '2024-06-15'),
(1, 1, 2, 90, 'Take one tablet twice daily with meals', '2024-06-20'),
(2, 1, 1, 30, 'Take one tablet daily with food', '2024-06-25'),
(3, 2, 7, 30, 'Take one tablet daily in the morning', '2024-05-15'),
(4, 3, 8, 60, 'Take one tablet twice daily', '2024-05-20');