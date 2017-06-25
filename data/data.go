// Package "data" contains utility functions for working with data objects.
// ToDo: 
// Настроить запись поля "контент" в БД
package data

import (
    "database/sql"
    _ "github.com/lib/pq"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

type DataInput struct {
    Id int `json:"id"`
    Name string `json:"name"`
    Project int `json:"project"`
    Parent int `json:"parent"`
    Coordinates []byte `json:"coordinates"`
    Content []byte `json:"content"`
}
type DataOutput struct {
    Id int `json:"id"`
    Name string `json:"name"`
    Project int `json:"project"`
    Parent int `json:"parent"`
    Coordinates map[string]int `json:"coordinates"`
    Content []FieldGroup `json:"content"`
}
type FieldGroup struct {
    Name string `json:"name"`
    Order int `json:"order"`
    Fields []Field `json:"fields"`
}
type Field struct {
    Type string `json:"type"`
    Value string `json:"value"`
    Order int `json:"order"`
}

var db *sql.DB
var err error

func init() {
    db, err = sql.Open("postgres", "user=urivsky password=123581321 dbname=mindassistant sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }

    if err = db.Ping(); err != nil {
        log.Fatal(err)
    }
}

// Function "List" show list of data by id
func List(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        http.Error(w, http.StatusText(405), 405)
        return
    }

    rows, err := db.Query("SELECT * FROM data")
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }
    defer rows.Close()

    dataList := make([]*DataOutput, 0)
    for rows.Next() {
        data := new(DataInput)
        
        err := rows.Scan(&data.Id, &data.Name, &data.Project, &data.Parent, &data.Coordinates, &data.Content)
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            return
        }
        
        var coordinates map[string]int
        var content []FieldGroup

        json.Unmarshal([]byte(data.Coordinates), &coordinates)
        json.Unmarshal([]byte(data.Content), &content)

        output := &DataOutput{
            Id: data.Id,
            Name: data.Name,
            Project: data.Project,
            Parent: data.Parent,
            Coordinates: coordinates,
            Content: content,
        }

        dataList = append(dataList, output)
    }
    if err = rows.Err(); err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    result, err := json.Marshal(&dataList)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(result)
}

// Function "ListByProject" show list of data by project id
func ListByProject(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        http.Error(w, http.StatusText(405), 405)
        return
    }

    project := r.FormValue("project")
    if project == "" {
        http.Error(w, http.StatusText(400), 400)
        return
    }

    rows, err := db.Query("SELECT * FROM data WHERE project = $1", project)
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }
    defer rows.Close()

    dataList := make([]*DataOutput, 0)
    for rows.Next() {
        data := new(DataInput)
        
        err := rows.Scan(&data.Id, &data.Name, &data.Project, &data.Parent, &data.Coordinates, &data.Content)
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            return
        }
        
        var coordinates map[string]int
        var content []FieldGroup

        json.Unmarshal([]byte(data.Coordinates), &coordinates)
        json.Unmarshal([]byte(data.Content), &content)

        output := &DataOutput{
            Id: data.Id,
            Name: data.Name,
            Project: data.Project,
            Parent: data.Parent,
            Coordinates: coordinates,
            Content: content,
        }

        dataList = append(dataList, output)
    }
    if err = rows.Err(); err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    result, err := json.Marshal(&dataList)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(result)
}

// Function "Item" show data by id
func Item(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        http.Error(w, http.StatusText(405), 405)
        return
    }

    id := r.FormValue("id")
    if id == "" {
        http.Error(w, http.StatusText(400), 400)
        return
    }

    row := db.QueryRow("SELECT * FROM data WHERE id = $1", id)
    data := new(DataInput)

    err := row.Scan(&data.Id, &data.Name, &data.Project, &data.Parent, &data.Coordinates, &data.Content)
    if err == sql.ErrNoRows {
        http.NotFound(w, r)
        return
    } else if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }
    
    var coordinates map[string]int
    var content []FieldGroup

    json.Unmarshal([]byte(data.Coordinates), &coordinates)
    json.Unmarshal([]byte(data.Content), &content)

    output := &DataOutput{
        Id: data.Id,
        Name: data.Name,
        Project: data.Project,
        Parent: data.Parent,
        Coordinates: coordinates,
        Content: content,
    }    

    result, err := json.Marshal(output)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(result)
}

// Function "Create" creates a new data object by json
func Create(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, http.StatusText(405), 405)
        return
    }

    decoder := json.NewDecoder(r.Body)
    data := new(DataOutput)
    err := decoder.Decode(&data)
    if err != nil {
        panic(err)
    }
    defer r.Body.Close()

    if data.Name == "" || data.Project <= 0 {
        http.Error(w, http.StatusText(400), 400)
        return
    }

    if data.Coordinates == nil {
        data.Coordinates = map[string]int{"x":0,"y":0}
    }

    stmt, err := db.Prepare("INSERT INTO data(name, project) VALUES($1, $2);")
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    result, err := stmt.Exec(data.Name, data.Project)
    if err != nil {
        log.Fatal(err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    if rowsAffected > 0 {
        fmt.Fprintf(w, "%t\n", true)
    }
}

// Function "Update" updates a new data by json
func Update(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, http.StatusText(405), 405)
        return
    }

    decoder := json.NewDecoder(r.Body)
    data := new(DataOutput)
    err := decoder.Decode(&data)
    if err != nil {
        panic(err)
    }
    defer r.Body.Close()    

    if data.Id == 0 || data.Project <= 0 {
        http.Error(w, http.StatusText(400), 400)
        return
    }

    if data.Coordinates == nil {
        data.Coordinates = map[string]int{"x":0,"y":0}
    }
    coordinates, _ := json.Marshal(data.Coordinates)

    var content []byte
    if data.Content != nil && len(data.Content) > 0 {
        content, _ = json.Marshal(data.Content)
    } else {
        content = nil
        // content = pq.Array(content) 
    }

    stmt, err := db.Prepare("update data set name = $2, project = $3, parent = $4, coordinates = $5, content = $6 where id=$1")
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    result, err := stmt.Exec(data.Id, data.Name, data.Project, data.Parent, coordinates, content)
    if err != nil {
        log.Fatal(err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    if rowsAffected > 0 {
        fmt.Fprintf(w, "%t\n", true)
    }
}

// Function "Delete" delete a data by id
func Delete(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, http.StatusText(405), 405)
        return
    }

    id := r.FormValue("id")
    if id == "" {
        http.Error(w, http.StatusText(400), 400)
        return
    }

    stmt, err := db.Prepare("delete from data where id=$1")
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    result, err := stmt.Exec(id)
    if err != nil {
        log.Fatal(err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    if rowsAffected > 0 {
        fmt.Fprintf(w, "%t\n", true)
    }
}