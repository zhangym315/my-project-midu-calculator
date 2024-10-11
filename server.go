package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

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

// Structs for decoding the two arrays
type MiduData struct {
	MinMidu string `json:"minMidu"`
	MaxMidu string `json:"maxMidu"`
	Price   string `json:"price"`
}

type ColumnData struct {
	Col1 string `json:"col1"`
	Col2 string `json:"col2"`
	Col3 string `json:"col3"`
	Col4 string `json:"col4"`
	Col5 string `json:"col5"`
}

type ColumnDataReturnedToHtml struct {
	Col0 string `json:"col0"`
	Col1 string `json:"col1"`
	Col2 string `json:"col2"`
	Col3 string `json:"col3"`
	Col4 string `json:"col4"`
	Col5 string `json:"col5"`
}

// Function to generate all permutations of the elements
func generatePermutations(arr []string) [][]string {
	var result [][]string

	// Recursive function to generate permutations
	var permute func(start int)
	permute = func(start int) {
		// If we've reached the end of the array, add the permutation to the result
		if start == len(arr)-1 {
			// Append a copy of arr to result
			result = append(result, append([]string(nil), arr...))
			return
		}

		// Loop through the array, generating permutations
		for i := start; i < len(arr); i++ {
			// Swap the current element with the starting element
			arr[start], arr[i] = arr[i], arr[start]
			// Recursively generate permutations with the next elements
			permute(start + 1)
			// Backtrack: swap back to the original positions
			arr[start], arr[i] = arr[i], arr[start]
		}
	}

	// Start the permutation generation process
	permute(0)
	return result
}

// Function to generate all partitions of the elements
func partition(arr []string) [][][]string {
	// This will store all partitions
	var result [][][]string

	// Helper function for recursively generating partitions
	var generate func([][]string, int)
	generate = func(current [][]string, index int) {
		// If we've processed all elements, add the current partition to the result
		if index == len(arr) {
			// Append a copy of current partition
			partitionCopy := make([][]string, len(current))
			for i := range current {
				partitionCopy[i] = append([]string(nil), current[i]...)
			}
			result = append(result, partitionCopy)
			return
		}

		// Try to add arr[index] to an existing subset in current partition
		for i := range current {
			current[i] = append(current[i], arr[index])
			generate(current, index+1)
			current[i] = current[i][:len(current[i])-1] // backtrack
		}

		// Or, create a new subset with arr[index]
		newSubset := []string{arr[index]}
		generate(append(current, newSubset), index+1)
	}

	// Start the recursive partition generation
	generate([][]string{}, 0)
	return result
}

func submitAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Read the entire body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		// Print the body as a string to the console
		fmt.Println("Request Body:", string(body))
		// JSON data
		jsonData := string(body)

		// Declare variables to hold the decoded data
		var miduData []MiduData
		var columnData []ColumnData

		// Wrap the two arrays in a single slice of interface{} to decode both
		var parsedData []interface{}

		// Parse the JSON data
		err = json.Unmarshal([]byte(jsonData), &parsedData)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		// Parse the first array (Midu data)
		miduArray, err := json.Marshal(parsedData[0])
		if err != nil {
			fmt.Println("Error encoding first array:", err)
			return
		}
		err = json.Unmarshal(miduArray, &miduData)
		if err != nil {
			fmt.Println("Error decoding Midu data:", err)
			return
		}

		// Parse the second array (Column data)
		columnArray, err := json.Marshal(parsedData[1])
		if err != nil {
			fmt.Println("Error encoding second array:", err)
			return
		}
		err = json.Unmarshal(columnArray, &columnData)
		if err != nil {
			fmt.Println("Error decoding Column data:", err)
			return
		}

		// Output the parsed data
		fmt.Println("Midu Data:", miduData)
		fmt.Println("Column Data:", columnData)

		// Print received data
		var array []string
		for _, row := range columnData {
			array = append(array, row.Col2)
			fmt.Printf("Received row: col1=%s, col2=%s, col3=%s, col4=%s, col5=%s\n", row.Col1, row.Col2, row.Col3, row.Col4, row.Col5)
		}

		for _, row := range miduData {
			fmt.Printf("Received row: col1=%s, col2=%s, col3=%s, col4=%s, col5=%s\n", row.MinMidu, row.MaxMidu, row.Price)
		}

		fmt.Println("array:", array)

		var returnedData []ColumnDataReturnedToHtml

		allCombinations := partition(array)
		indexNumber := 1
		for _, e := range allCombinations {
			var rowData string
			for _, r := range e {
				rowData = rowData + "(" + strings.Join(r, ", ") + ")\n"
			}
			returnedData = append(returnedData, ColumnDataReturnedToHtml{Col0: strconv.Itoa(indexNumber), Col1: rowData})
			indexNumber++
		}

		// Respond back with a success message
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(returnedData)

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
