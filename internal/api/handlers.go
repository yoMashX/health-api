package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"health-api/internal/models"
)

type Handlers struct {
	db *sql.DB
}

func NewHandlers(db *sql.DB) *Handlers {
	return &Handlers{db: db}
}

func (h *Handlers) CreatePrescription(w http.ResponseWriter, r *http.Request) {
	role, ok := GetRoleFromContext(r.Context())
	if !ok {
		http.Error(w, "Role not found in context", http.StatusInternalServerError)
		return
	}

	if role == models.RolePatient {
		http.Error(w, "Patients cannot create prescriptions", http.StatusForbidden)
		return
	}

	var req models.CreatePrescriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := validatePrescriptionRequest(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if role == models.RolePhysician {
		userID, ok := GetUserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}

		physicianID, err := strconv.Atoi(userID)
		if err != nil {
			http.Error(w, "Invalid physician ID", http.StatusBadRequest)
			return
		}

		if req.PhysicianID != physicianID {
			http.Error(w, "Physician can only create prescriptions for themselves", http.StatusForbidden)
			return
		}

		canPrescribe, err := h.canPhysicianPrescribeToPatient(physicianID, req.PatientID)
		if err != nil {
			http.Error(w, "Error checking physician-patient relationship", http.StatusInternalServerError)
			return
		}
		if !canPrescribe {
			http.Error(w, "Physician is not authorized to prescribe to this patient", http.StatusForbidden)
			return
		}
	}

	if !h.entityExists("patients", req.PatientID) {
		http.Error(w, "Patient not found", http.StatusBadRequest)
		return
	}
	if !h.entityExists("physicians", req.PhysicianID) {
		http.Error(w, "Physician not found", http.StatusBadRequest)
		return
	}
	if !h.entityExists("drugs", req.DrugID) {
		http.Error(w, "Drug not found", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO prescriptions (patient_id, physician_id, drug_id, quantity, sig, prescribed_date)
		VALUES ($1, $2, $3, $4, $5, CURRENT_DATE)
		RETURNING id, prescribed_date, created_at, updated_at`

	var prescription models.Prescription
	err := h.db.QueryRow(query, req.PatientID, req.PhysicianID, req.DrugID, req.Quantity, req.Sig).Scan(
		&prescription.ID, &prescription.PrescribedDate, &prescription.CreatedAt, &prescription.UpdatedAt)
	if err != nil {
		http.Error(w, "Failed to create prescription", http.StatusInternalServerError)
		return
	}

	prescription.PatientID = req.PatientID
	prescription.PhysicianID = req.PhysicianID
	prescription.DrugID = req.DrugID
	prescription.Quantity = req.Quantity
	prescription.Sig = req.Sig

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(prescription)
}

func (h *Handlers) GetTopDrugs(w http.ResponseWriter, r *http.Request) {
	role, ok := GetRoleFromContext(r.Context())
	if !ok {
		http.Error(w, "Role not found in context", http.StatusInternalServerError)
		return
	}

	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")
	limitStr := r.URL.Query().Get("limit")

	if limitStr == "" {
		limitStr = "10"
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		http.Error(w, "Invalid limit parameter. Must be between 1 and 100", http.StatusBadRequest)
		return
	}

	var fromDate, toDate time.Time
	if fromStr != "" {
		fromDate, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			http.Error(w, "Invalid from date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	} else {
		fromDate = time.Now().AddDate(0, -6, 0)
	}

	if toStr != "" {
		toDate, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			http.Error(w, "Invalid to date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	} else {
		toDate = time.Now()
	}

	var query string
	var args []interface{}

	if role == models.RolePatient {
		userID, ok := GetUserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}

		patientID, err := strconv.Atoi(userID)
		if err != nil {
			http.Error(w, "Invalid patient ID", http.StatusBadRequest)
			return
		}

		query = `
			SELECT 
				d.id as drug_id,
				d.name as drug_name,
				COALESCE(d.generic_name, '') as generic_name,
				SUM(p.quantity) as total_quantity,
				COUNT(p.id) as prescription_count
			FROM prescriptions p
			JOIN drugs d ON p.drug_id = d.id
			WHERE p.patient_id = $1 AND p.prescribed_date >= $2 AND p.prescribed_date <= $3
			GROUP BY d.id, d.name, d.generic_name
			ORDER BY total_quantity DESC
			LIMIT $4`
		args = []interface{}{patientID, fromDate, toDate, limit}
	} else {
		query = `
			SELECT 
				d.id as drug_id,
				d.name as drug_name,
				COALESCE(d.generic_name, '') as generic_name,
				SUM(p.quantity) as total_quantity,
				COUNT(p.id) as prescription_count
			FROM prescriptions p
			JOIN drugs d ON p.drug_id = d.id
			WHERE p.prescribed_date >= $1 AND p.prescribed_date <= $2
			GROUP BY d.id, d.name, d.generic_name
			ORDER BY total_quantity DESC
			LIMIT $3`
		args = []interface{}{fromDate, toDate, limit}
	}

	rows, err := h.db.Query(query, args...)
	if err != nil {
		http.Error(w, "Failed to query top drugs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var topDrugs []models.TopDrug
	for rows.Next() {
		var drug models.TopDrug
		err := rows.Scan(&drug.DrugID, &drug.DrugName, &drug.GenericName, &drug.TotalQuantity, &drug.PrescriptionCount)
		if err != nil {
			http.Error(w, "Failed to scan drug data", http.StatusInternalServerError)
			return
		}
		topDrugs = append(topDrugs, drug)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating over results", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topDrugs)
}

func validatePrescriptionRequest(req models.CreatePrescriptionRequest) error {
	if req.PatientID <= 0 {
		return fmt.Errorf("patient_id must be positive")
	}
	if req.PhysicianID <= 0 {
		return fmt.Errorf("physician_id must be positive")
	}
	if req.DrugID <= 0 {
		return fmt.Errorf("drug_id must be positive")
	}
	if req.Quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}
	if req.Sig == "" {
		return fmt.Errorf("sig (instructions) cannot be empty")
	}
	return nil
}

func (h *Handlers) entityExists(table string, id int) bool {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", table)
	h.db.QueryRow(query, id).Scan(&exists)
	return exists
}

func (h *Handlers) canPhysicianPrescribeToPatient(physicianID, patientID int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM patient_physicians WHERE physician_id = $1 AND patient_id = $2)"
	err := h.db.QueryRow(query, physicianID, patientID).Scan(&exists)
	return exists, err
}