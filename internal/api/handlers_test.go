package api

import (
	"database/sql"
	"testing"
	"time"

	"health-api/internal/models"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=health_api_test sslmode=disable")
	if err != nil {
		t.Skip("Skipping test: unable to connect to test database")
	}

	if err = db.Ping(); err != nil {
		t.Skip("Skipping test: unable to ping test database")
	}

	// Clean up and set up test data
	_, err = db.Exec(`
		DELETE FROM prescriptions;
		DELETE FROM patient_physicians;
		DELETE FROM drugs;
		DELETE FROM patients;
		DELETE FROM physicians;
		
		INSERT INTO physicians (id, first_name, last_name, license_number, email) VALUES
		(1, 'John', 'Smith', 'MD123456', 'j.smith@hospital.com'),
		(2, 'Sarah', 'Johnson', 'MD234567', 's.johnson@clinic.com');
		
		INSERT INTO patients (id, first_name, last_name, date_of_birth, email) VALUES
		(1, 'Alice', 'Anderson', '1985-03-15', 'alice.anderson@email.com'),
		(2, 'Bob', 'Baker', '1990-07-22', 'bob.baker@email.com');
		
		INSERT INTO drugs (id, name, generic_name, strength, dosage_form) VALUES
		(1, 'Lipitor', 'Atorvastatin', '20mg', 'Tablet'),
		(2, 'Metformin', 'Metformin', '500mg', 'Tablet'),
		(3, 'Lisinopril', 'Lisinopril', '10mg', 'Tablet');
		
		INSERT INTO patient_physicians (patient_id, physician_id) VALUES
		(1, 1), (2, 1), (1, 2);
		
		-- Reset sequences
		ALTER SEQUENCE physicians_id_seq RESTART WITH 3;
		ALTER SEQUENCE patients_id_seq RESTART WITH 3;
		ALTER SEQUENCE drugs_id_seq RESTART WITH 4;
	`)
	if err != nil {
		t.Fatalf("Failed to set up test data: %v", err)
	}

	return db
}

func TestGetTopDrugsQuery(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handlers := NewHandlers(db)

	// Insert test prescriptions with known quantities and dates
	testCases := []struct {
		name           string
		prescriptions  []struct {
			patientID, physicianID, drugID, quantity int
			date                                     string
		}
		fromDate       string
		toDate         string
		limit          int
		expectedDrugs  []models.TopDrug
		expectedCount  int
	}{
		{
			name: "top drugs in date range",
			prescriptions: []struct {
				patientID, physicianID, drugID, quantity int
				date                                     string
			}{
				{1, 1, 1, 30, "2024-10-01"},  // Lipitor: 30
				{2, 1, 1, 60, "2024-10-02"},  // Lipitor: +60 = 90 total
				{1, 1, 2, 90, "2024-10-03"},  // Metformin: 90
				{2, 1, 3, 30, "2024-10-04"},  // Lisinopril: 30
				{1, 2, 2, 30, "2024-10-05"},  // Metformin: +30 = 120 total
			},
			fromDate: "2024-10-01",
			toDate:   "2024-10-05",
			limit:    3,
			expectedDrugs: []models.TopDrug{
				{DrugID: 2, DrugName: "Metformin", GenericName: "Metformin", TotalQuantity: 120, PrescriptionCount: 2},
				{DrugID: 1, DrugName: "Lipitor", GenericName: "Atorvastatin", TotalQuantity: 90, PrescriptionCount: 2},
				{DrugID: 3, DrugName: "Lisinopril", GenericName: "Lisinopril", TotalQuantity: 30, PrescriptionCount: 1},
			},
			expectedCount: 3,
		},
		{
			name: "limit results",
			prescriptions: []struct {
				patientID, physicianID, drugID, quantity int
				date                                     string
			}{
				{1, 1, 1, 50, "2024-10-01"},
				{1, 1, 2, 40, "2024-10-01"},
				{1, 1, 3, 30, "2024-10-01"},
			},
			fromDate: "2024-10-01",
			toDate:   "2024-10-01",
			limit:    2,
			expectedDrugs: []models.TopDrug{
				{DrugID: 1, DrugName: "Lipitor", GenericName: "Atorvastatin", TotalQuantity: 50, PrescriptionCount: 1},
				{DrugID: 2, DrugName: "Metformin", GenericName: "Metformin", TotalQuantity: 40, PrescriptionCount: 1},
			},
			expectedCount: 2,
		},
		{
			name: "date range filtering",
			prescriptions: []struct {
				patientID, physicianID, drugID, quantity int
				date                                     string
			}{
				{1, 1, 1, 100, "2024-09-30"}, // Outside range
				{1, 1, 2, 50, "2024-10-01"},  // Inside range
				{1, 1, 3, 25, "2024-10-06"},  // Outside range
			},
			fromDate: "2024-10-01",
			toDate:   "2024-10-05",
			limit:    10,
			expectedDrugs: []models.TopDrug{
				{DrugID: 2, DrugName: "Metformin", GenericName: "Metformin", TotalQuantity: 50, PrescriptionCount: 1},
			},
			expectedCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean prescriptions table
			_, err := db.Exec("DELETE FROM prescriptions")
			if err != nil {
				t.Fatalf("Failed to clean prescriptions: %v", err)
			}

			// Insert test prescriptions
			for _, p := range tc.prescriptions {
				_, err := db.Exec(`
					INSERT INTO prescriptions (patient_id, physician_id, drug_id, quantity, sig, prescribed_date)
					VALUES ($1, $2, $3, $4, 'Test instructions', $5)`,
					p.patientID, p.physicianID, p.drugID, p.quantity, p.date)
				if err != nil {
					t.Fatalf("Failed to insert test prescription: %v", err)
				}
			}

			// Execute the top drugs query (admin view - no patient restriction)
			fromDate, _ := time.Parse("2006-01-02", tc.fromDate)
			toDate, _ := time.Parse("2006-01-02", tc.toDate)

			query := `
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

			rows, err := db.Query(query, fromDate, toDate, tc.limit)
			if err != nil {
				t.Fatalf("Query failed: %v", err)
			}
			defer rows.Close()

			var results []models.TopDrug
			for rows.Next() {
				var drug models.TopDrug
				err := rows.Scan(&drug.DrugID, &drug.DrugName, &drug.GenericName, &drug.TotalQuantity, &drug.PrescriptionCount)
				if err != nil {
					t.Fatalf("Failed to scan result: %v", err)
				}
				results = append(results, drug)
			}

			// Verify results
			if len(results) != tc.expectedCount {
				t.Errorf("Expected %d results, got %d", tc.expectedCount, len(results))
			}

			for i, expected := range tc.expectedDrugs {
				if i >= len(results) {
					t.Errorf("Missing expected result at index %d: %+v", i, expected)
					continue
				}

				result := results[i]
				if result.DrugID != expected.DrugID {
					t.Errorf("Expected drug ID %d at index %d, got %d", expected.DrugID, i, result.DrugID)
				}
				if result.DrugName != expected.DrugName {
					t.Errorf("Expected drug name %s at index %d, got %s", expected.DrugName, i, result.DrugName)
				}
				if result.TotalQuantity != expected.TotalQuantity {
					t.Errorf("Expected total quantity %d at index %d, got %d", expected.TotalQuantity, i, result.TotalQuantity)
				}
				if result.PrescriptionCount != expected.PrescriptionCount {
					t.Errorf("Expected prescription count %d at index %d, got %d", expected.PrescriptionCount, i, result.PrescriptionCount)
				}
			}
		})
	}
}