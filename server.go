package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/tealeg/xlsx"
)

type MiduPriceData struct {
	StartMidu string `json:"startMidu"`
	EndMidu   string `json:"endMidu"`
	MiduPrice string `json:"miduPrice"`
}

// Struct to hold all form data
type MiduPricesData struct {
	MiDuPricesRows []MiduPriceData
}

// Struct to hold form data for each row
type RowData struct {
	Col1 string `json:"col1"`
	Col2 string `json:"col2"`
	Col3 string `json:"col3"`
	Col4 string `json:"col4"`
	Col5 string `json:"col5"`
}

// Struct to hold all form data
type FormData struct {
	Rows []RowData
}

// Global variable to store the form data
var formDataMidu = MiduPricesData{MiDuPricesRows: make([]MiduPriceData, 1)}
var formData = FormData{Rows: make([]RowData, 1)}

type CombinedFormData struct {
	MiduPricesData
	FormData
}

// Handler to serve the form
func formHandler(w http.ResponseWriter, r *http.Request) {
	// Load the template file
	tmplPath := filepath.Join("templates", "form.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	combinedData := CombinedFormData{
		MiduPricesData: formDataMidu,
		FormData:       formData,
	}

	// Execute the template with form data
	err = tmpl.Execute(w, combinedData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Auto-save handler to update form data
func autoSaveHandler(w http.ResponseWriter, r *http.Request) {
	var update struct {
		RowIndex int    `json:"rowIndex"`
		ColName  string `json:"colName"`
		Value    string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Failed to parse request", http.StatusBadRequest)
		return
	}

	// Ensure the row exists, or create it
	for len(formData.Rows) <= update.RowIndex {
		formData.Rows = append(formData.Rows, RowData{})
	}

	// Update the correct column based on the column name
	switch update.ColName {
	case "col1":
		formData.Rows[update.RowIndex].Col1 = update.Value
	case "col2":
		formData.Rows[update.RowIndex].Col2 = update.Value
	case "col3":
		formData.Rows[update.RowIndex].Col3 = update.Value
	case "col4":
		formData.Rows[update.RowIndex].Col4 = update.Value
	case "col5":
		formData.Rows[update.RowIndex].Col5 = update.Value
	}
}

// Handler to add a new row to the form
func addRowMiduHandler(w http.ResponseWriter, r *http.Request) {
	formDataMidu.MiDuPricesRows = append(formDataMidu.MiDuPricesRows, MiduPriceData{})
	w.WriteHeader(http.StatusOK)
}

// Handler to add a new row to the form
func addRowHandler(w http.ResponseWriter, r *http.Request) {
	formData.Rows = append(formData.Rows, RowData{})
	w.WriteHeader(http.StatusOK)
}

// Handler to submit the form
func submitHandler(w http.ResponseWriter, r *http.Request) {
	// Return the submitted form data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(formData.Rows)
}

// Export XLSX file handler
func deleteRowHandler(w http.ResponseWriter, r *http.Request) {

}

func exportXLSXHandler(w http.ResponseWriter, r *http.Request) {
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	row := sheet.AddRow()

	// Add header
	row.WriteSlice(&[]string{"Column 1", "Column 2", "Column 3", "Column 4", "Column 5"}, -1)

	// Add data rows
	for _, dataRow := range formData.Rows {
		row = sheet.AddRow()
		row.WriteSlice(&[]string{dataRow.Col1, dataRow.Col2, dataRow.Col3, dataRow.Col4, dataRow.Col5}, -1)
	}

	// Write file to response
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", `attachment; filename="form_data.xlsx"`)
	file.Write(w)
}

func submitAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse the form data
		r.ParseForm()
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {
	// Serve static files (e.g., CSS and JS files)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Handle form and actions
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/autosave", autoSaveHandler)       // Auto-save endpoint
	http.HandleFunc("/add-row", addRowHandler)          // Add new row endpoint
	http.HandleFunc("/add-row-midu", addRowMiduHandler) // Add new midu row endpoint
	http.HandleFunc("/submit", submitHandler)           // Submit form endpoint
	http.HandleFunc("/download-xlsx", exportXLSXHandler)
	http.HandleFunc("/delete-row", deleteRowHandler)
	http.HandleFunc("/submitall", submitAllHandler)

	fmt.Println("Starting server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
