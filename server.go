package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
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
	file := excelize.NewFile()
	sheet := "Sheet1"

	// Add header
	headers := []string{"Column 1", "Column 2", "Column 3", "Column 4", "Column 5"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		file.SetCellValue(sheet, cell, header)
	}

	// Add data rows
	for i, dataRow := range formData.Rows {
		rowData := []string{dataRow.Col1, dataRow.Col2, dataRow.Col3, dataRow.Col4, dataRow.Col5}
		for j, value := range rowData {
			cell := fmt.Sprintf("%c%d", 'A'+j, i+2)
			file.SetCellValue(sheet, cell, value)
		}
	}

	// Write file to response
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", `attachment; filename="form_data.xlsx"`)
	if err := file.Write(w); err != nil {
		http.Error(w, "Failed to generate Excel file", http.StatusInternalServerError)
		return
	}
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
	Col0       string `json:"col0"`
	Col1       string `json:"col1"`
	Col2       string `json:"col2"`
	Col3       string `json:"col3"`
	Col4       string `json:"col4"`
	Col5       string `json:"col5"`
	TotalPrice float64
}

type PriceTable struct {
	Min   float64
	Max   float64
	Price float64
}

type HeavyAndVolume struct {
	H float64
	V float64
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

func getPriceFromTable(t []PriceTable, midu float64) float64 {
	for _, row := range t {
		if midu <= row.Max && midu >= row.Min {
			return row.Price * midu
		}
	}

	fmt.Printf("Error: failed to get price from midu=%f\n", midu)
	return 0
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
		//		fmt.Println("Request Body:", string(body))
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

		// construct price table
		var pricetable []PriceTable
		for _, row := range miduData {
			num1, err := strconv.ParseFloat(row.MinMidu, 64) // 64 specifies double precision
			if err != nil {
				fmt.Printf("Error parsing midu row: MinMidu=%s, MaxMidu=%s, Price=%s\n", row.MinMidu, row.MaxMidu, row.Price)
			}

			num2, err := strconv.ParseFloat(row.MaxMidu, 64) // 64 specifies double precision
			if err != nil {
				fmt.Printf("Error parsing midu row: MinMidu=%s, MaxMidu=%s, Price=%s\n", row.MinMidu, row.MaxMidu, row.Price)
			}

			num3, err := strconv.ParseFloat(row.Price, 64) // 64 specifies double precision
			if err != nil {
				fmt.Printf("Error parsing midu row: MinMidu=%s, MaxMidu=%s, Price=%s\n", row.MinMidu, row.MaxMidu, row.Price)
			}

			pricetable = append(pricetable, PriceTable{num1, num2, num3})
			//fmt.Printf("Received row: MinMidu=%s, MaxMidu=%s, Price=%s\n", row.MinMidu, row.MaxMidu, row.Price)
		}

		// Print received data
		// huohao to row
		var array []string
		numberToHV := make(map[string]HeavyAndVolume)
		for _, row := range columnData {
			array = append(array, row.Col2)
			parsedH, err := strconv.ParseFloat(row.Col3, 64)
			if err != nil {
				fmt.Printf("Error: Failed to parse 重量=%s, 货号=%s\n", row.Col3, row.Col2)
			}
			parsedV, err := strconv.ParseFloat(row.Col4, 64)
			if err != nil {
				fmt.Printf("Error: Failed to parse 体积=%s, 货号=%s\n", row.Col4, row.Col2)
			}

			numberToHV[row.Col2] = HeavyAndVolume{parsedH, parsedV}
			//			fmt.Printf("Received row: col1=%s, col2=%s, col3=%s, col4=%s, col5=%s\n", row.Col1, row.Col2, row.Col3, row.Col4, row.Col5)
		}

		var returnedData []ColumnDataReturnedToHtml

		allCombinations := partition(array)
		/*
		 * allCombinations is like:
		 *
		 */
		for _, e := range allCombinations {

			var rowData string
			var totalPrice float64
			var totalH string
			var totalV string
			var totalHV string
			for _, r := range e {
				rowData = rowData + "(" + strings.Join(r, ", ") + ")"
				var cH float64
				var cV float64
				var sH string
				var sV string
				for _, element := range r {
					cH += (numberToHV[element].H)
					cV += (numberToHV[element].V)
					sH += strings.TrimRight(strings.TrimRight(strconv.FormatFloat(numberToHV[element].H, 'f', 3, 64), "0"), ".") + " "
					sV += strings.TrimRight(strings.TrimRight(strconv.FormatFloat(numberToHV[element].V, 'f', 3, 64), "0"), ".") + " "

				}
				totalH = totalH + "(" + sH + ") "
				totalV = totalV + "(" + sV + ") "
				totalHV = totalHV + "(" + strings.TrimRight(strings.TrimRight(strconv.FormatFloat(cH/cV, 'f', 3, 64), "0"), ".") + ") "
				totalPrice += getPriceFromTable(pricetable, (cH / cV))
			}

			returnedData = append(returnedData, ColumnDataReturnedToHtml{Col0: "0", Col1: rowData,
				Col2: totalH, Col3: totalV, Col4: totalHV,
				Col5:       strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.3f", totalPrice), "0"), "."),
				TotalPrice: totalPrice})
		}

		// Sort the slice by the Age field
		sort.Slice(returnedData, func(i, j int) bool {
			return returnedData[i].TotalPrice < returnedData[j].TotalPrice
		})

		indexNumber := 1
		for i := range returnedData {
			returnedData[i].Col0 = strconv.Itoa(indexNumber)
			indexNumber++
		}

		// Respond back with a success message
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(returnedData)

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func exportExcelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		MiduData   []map[string]string `json:"miduData"`
		FormData   []map[string]string `json:"formData"`
		ResultData []map[string]string `json:"resultData"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	//	log.Printf("Received data: %+v", data)

	f := excelize.NewFile()

	// 密度价格表 sheet
	f.SetSheetName("Sheet1", "密度价格表")
	f.SetCellValue("密度价格表", "A1", "密度最小")
	f.SetCellValue("密度价格表", "B1", "密度最大")
	f.SetCellValue("密度价格表", "C1", "价格")
	for i, row := range data.MiduData {
		f.SetCellValue("密度价格表", fmt.Sprintf("A%d", i+2), row["minMidu"])
		f.SetCellValue("密度价格表", fmt.Sprintf("B%d", i+2), row["maxMidu"])
		f.SetCellValue("密度价格表", fmt.Sprintf("C%d", i+2), row["price"])
	}

	// 拼货 sheet
	f.NewSheet("拼货")
	f.SetCellValue("拼货", "A1", "运输方式")
	f.SetCellValue("拼货", "B1", "货号")
	f.SetCellValue("拼货", "C1", "重量")
	f.SetCellValue("拼货", "D1", "体积")
	f.SetCellValue("拼货", "E1", "单票密度")
	for i, row := range data.FormData {
		f.SetCellValue("拼货", fmt.Sprintf("A%d", i+2), row["col1"])
		f.SetCellValue("拼货", fmt.Sprintf("B%d", i+2), row["col2"])
		f.SetCellValue("拼货", fmt.Sprintf("C%d", i+2), row["col3"])
		f.SetCellValue("拼货", fmt.Sprintf("D%d", i+2), row["col4"])
		f.SetCellValue("拼货", fmt.Sprintf("E%d", i+2), row["col5"])
	}

	// 生成的数据 sheet
	f.NewSheet("生成的数据")
	f.SetCellValue("生成的数据", "A1", "序号")
	f.SetCellValue("生成的数据", "B1", "拼货组合")
	f.SetCellValue("生成的数据", "C1", "拼货重量")
	f.SetCellValue("生成的数据", "D1", "拼货体积")
	f.SetCellValue("生成的数据", "E1", "拼货密度")
	f.SetCellValue("生成的数据", "F1", "拼货价格")
	for i, row := range data.ResultData {
		f.SetCellValue("生成的数据", fmt.Sprintf("A%d", i+2), row["col0"])
		f.SetCellValue("生成的数据", fmt.Sprintf("B%d", i+2), row["col1"])
		f.SetCellValue("生成的数据", fmt.Sprintf("C%d", i+2), row["col2"])
		f.SetCellValue("生成的数据", fmt.Sprintf("D%d", i+2), row["col3"])
		f.SetCellValue("生成的数据", fmt.Sprintf("E%d", i+2), row["col4"])
		f.SetCellValue("生成的数据", fmt.Sprintf("F%d", i+2), row["col5"])
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=exported_data.xlsx")

	if err := f.Write(w); err != nil {
		log.Printf("Error writing Excel file: %v", err)
		http.Error(w, "Failed to generate Excel file", http.StatusInternalServerError)
	} else {
		log.Println("Excel file generated successfully")
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
	http.HandleFunc("/export-excel", exportExcelHandler)

	fmt.Println("Starting server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
