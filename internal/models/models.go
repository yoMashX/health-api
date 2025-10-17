package models

import "time"

type Physician struct {
	ID            int       `json:"id" db:"id"`
	FirstName     string    `json:"first_name" db:"first_name"`
	LastName      string    `json:"last_name" db:"last_name"`
	LicenseNumber string    `json:"license_number" db:"license_number"`
	Email         string    `json:"email" db:"email"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type Patient struct {
	ID          int       `json:"id" db:"id"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	DateOfBirth time.Time `json:"date_of_birth" db:"date_of_birth"`
	Email       string    `json:"email" db:"email"`
	Phone       *string   `json:"phone" db:"phone"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Drug struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	GenericName *string   `json:"generic_name" db:"generic_name"`
	Strength    *string   `json:"strength" db:"strength"`
	DosageForm  *string   `json:"dosage_form" db:"dosage_form"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Prescription struct {
	ID             int       `json:"id" db:"id"`
	PatientID      int       `json:"patient_id" db:"patient_id"`
	PhysicianID    int       `json:"physician_id" db:"physician_id"`
	DrugID         int       `json:"drug_id" db:"drug_id"`
	Quantity       int       `json:"quantity" db:"quantity"`
	Sig            string    `json:"sig" db:"sig"`
	PrescribedDate time.Time `json:"prescribed_date" db:"prescribed_date"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type CreatePrescriptionRequest struct {
	PatientID   int    `json:"patient_id" validate:"required"`
	PhysicianID int    `json:"physician_id" validate:"required"`
	DrugID      int    `json:"drug_id" validate:"required"`
	Quantity    int    `json:"quantity" validate:"required,min=1"`
	Sig         string `json:"sig" validate:"required"`
}

type TopDrug struct {
	DrugID         int    `json:"drug_id" db:"drug_id"`
	DrugName       string `json:"drug_name" db:"drug_name"`
	GenericName    string `json:"generic_name" db:"generic_name"`
	TotalQuantity  int    `json:"total_quantity" db:"total_quantity"`
	PrescriptionCount int `json:"prescription_count" db:"prescription_count"`
}

type Role string

const (
	RolePhysician Role = "physician"
	RolePatient   Role = "patient"
	RoleAdmin     Role = "admin"
)