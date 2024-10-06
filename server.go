package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/tealeg/xlsx"
)

// Template for the HTML form with dynamically addable rows
var tmpl = template.Must(template.New("form").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Dynamic Row Form</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f9;
        }
        .container {
            background-color: #fff;
            padding: 20px;
            margin-top: 20px;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
            width: 80%;
            margin: auto;
        }
        h1 {
            text-align: center;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
        }
        th, td {
            padding: 10px;
            border: 1px solid #ccc;
            text-align: center;
        }
        input[type="text"] {
            width: 100%;
            padding: 8px;
            border: 1px solid #ccc;
            border-radius: 4px;
        }
        .btn {
            margin-top: 10px;
            padding: 10px;
            background-color: #28a745;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
        }
        .btn:hover {
            background-color: #218838;
        }
        .delete-btn {
            background-color: #dc3545;
            color: white;
            border: none;
            padding: 5px 10px;
            cursor: pointer;
            border-radius: 4px;
        }
        .delete-btn:hover {
            background-color: #c82333;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Dynamic Row Form</h1>
        <form method="POST" action="/">
            <table id="dataTable">
                <thead>
                    <tr>
                        <th>Column 1</th>
                        <th>Column 2</th>
                        <th>Column 3</th>
                        <th>Column 4</th>
                        <th>Column 5</th>
                        <th>Column 6</th>
                        <th>Column 7</th>
                        <th>Action</th>
                    </tr>
                </thead>
                <tbody id="tableBody">
                    {{if .Submitted}}
                        {{range .Rows}}
                        <tr>
                            <td><input type="text" name="col1[]" value="{{.Col1}}"></td>
                            <td><input type="text" name="col2[]" value="{{.Col2}}"></td>
                            <td><input type="text" name="col3[]" value="{{.Col3}}"></td>
                            <td><input type="text" name="col4[]" value="{{.Col4}}"></td>
                            <td><input type="text" name="col5[]" value="{{.Col5}}"></td>
                            <td><input type="text" name="col6[]" value="{{.Col6}}"></td>
                            <td><input type="text" name="col7[]" value="{{.Col7}}"></td>
                            <td><button type="button" class="delete-btn" onclick="deleteRow(this)">Delete</button></td>
                        </tr>
                        {{end}}
                    {{else}}
                        <tr>
                            <td><input type="text" name="col1[]"></td>
                            <td><input type="text" name="col2[]"></td>
                            <td><input type="text" name="col3[]"></td>
                            <td><input type="text" name="col4[]"></td>
                            <td><input type="text" name="col5[]"></td>
                            <td><input type="text" name="col6[]"></td>
                            <td><input type="text" name="col7[]"></td>
                            <td><button type="button" class="delete-btn" onclick="deleteRow(this)">Delete</button></td>
                        </tr>
                    {{end}}
                </tbody>
            </table>
            <button type="button" class="btn" onclick="addRow()">Add Row</button><br><br>
            <input type="submit" value="Submit" class="btn">
        </form>

        {{if .Submitted}}
        <h2>Submitted Data</h2>
        <table>
            <thead>
                <tr>
                    <th>Column 1</th>
                    <th>Column 2</th>
                    <th>Column 3</th>
                    <th>Column 4</th>
                    <th>Column 5</th>
                    <th>Column 6</th>
                    <th>Column 7</th>
                </tr>
            </thead>
            <tbody>
                {{range .Rows}}
                <tr>
                    <td>{{.Col1}}</td>
                    <td>{{.Col2}}</td>
                    <td>{{.Col3}}</td>
                    <td>{{.Col4}}</td>
                    <td>{{.Col5}}</td>
                    <td>{{.Col6}}</td>
                    <td>{{.Col7}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
        <br>
        <a href="/download-xlsx" class="btn">Download as XLSX</a>
        {{end}}
    </div>

    <script>
        // Function to add a new row
        function addRow() {
            var tableBody = document.getElementById("tableBody");
            var row = document.createElement("tr");

            row.innerHTML = 
                '<td><input type="text" name="col1[]"></td>' +
                '<td><input type="text" name="col2[]"></td>' +
                '<td><input type="text" name="col3[]"></td>' +
                '<td><input type="text" name="col4[]"></td>' +
                '<td><input type="text" name="col5[]"></td>' +
                '<td><input type="text" name="col6[]"></td>' +
                '<td><input type="text" name="col7[]"></td>' +
                '<td><button type="button" class="delete-btn" onclick="deleteRow(this)">Delete</button></td>';

            tableBody.appendChild(row);
        }

        // Function to delete a row
        function deleteRow(button) {
            var row = button.parentElement.parentElement;
            row.remove();
        }
    </script>
</body>
</html>
`))

// Struct to hold form data for each row
type RowData struct {
	Col1 string
	Col2 string
	Col3 string
	Col4 string
	Col5 string
	Col6 string
	Col7 string
}

// Struct to hold all the form data
type FormData struct {
	Submitted bool
	Rows      []RowData
}

// Variable to store the form data globally for XLSX export
var formData FormData

// Handler to render the form and handle form submission
func formHandler(w http.ResponseWriter, r *http.Request) {
	formData = FormData{} // Reset the global formData

	if r.Method == http.MethodPost {
		// Parse the form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusInternalServerError)
			return
		}

		// Extract row data
		col1 := r.Form["col1[]"]
		col2 := r.Form["col2[]"]
		col3 := r.Form["col3[]"]
		col4 := r.Form["col4[]"]
		col5 := r.Form["col5[]"]
		col6 := r.Form["col6[]"]
		col7 := r.Form["col7[]"]

		// Populate rows
		for i := range col1 {
			row := RowData{
				Col1: col1[i],
				Col2: col2[i],
				Col3: col3[i],
				Col4: col4[i],
				Col5: col5[i],
				Col6: col6[i],
				Col7: col7[i],
			}
			formData.Rows = append(formData.Rows, row)
		}

		formData.Submitted = true
	}

	// Render the form and any submitted data
	tmpl.Execute(w, formData)
}

// Handler to export the form data as XLSX
func exportXLSXHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new Excel file
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("FormData")
	if err != nil {
		http.Error(w, "Failed to create sheet", http.StatusInternalServerError)
		return
	}

	// Add the header row
	headerRow := sheet.AddRow()
	headerRow.WriteSlice(&[]string{"Column 1", "Column 2", "Column 3", "Column 4", "Column 5", "Column 6", "Column 7"}, -1)

	// Add data rows
	for _, row := range formData.Rows {
		dataRow := sheet.AddRow()
		dataRow.WriteSlice(&[]string{row.Col1, row.Col2, row.Col3, row.Col4, row.Col5, row.Col6, row.Col7}, -1)
	}

	// Set the content type and trigger download
	w.Header().Set("Content-Disposition", "attachment;filename=data.xlsx")
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// Write the Excel file to the response writer
	err = file.Write(w)
	if err != nil {
		http.Error(w, "Failed to write XLSX file", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", formHandler)                    // Serve the form and handle submissions
	http.HandleFunc("/download-xlsx", exportXLSXHandler) // Serve the XLSX file

	fmt.Println("Starting server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
