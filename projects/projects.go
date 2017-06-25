// Package "projects" contains utility functions for working with projects.
package projects

import (
    "database/sql"
    _ "github.com/lib/pq"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

type Project struct {
    Id int
    Name string
    Pages int
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

// Function "List" show list of projects by id
func List(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        http.Error(w, http.StatusText(405), 405)
        return
    }

    rows, err := db.Query("SELECT * FROM projects")
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }
    defer rows.Close()

    projects := make([]*Project, 0)
    for rows.Next() {
        project := new(Project)
        
        err := rows.Scan(&project.Id, &project.Name, &project.Pages)
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            return
        }

        projects = append(projects, project)
    }
    if err = rows.Err(); err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    result, err := json.Marshal(projects)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(result)
}

// Function "Item" show project by id
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

    row := db.QueryRow("SELECT * FROM projects WHERE id = $1", id)
    project := new(Project)

    err := row.Scan(&project.Id, &project.Name, &project.Pages)
    if err == sql.ErrNoRows {
        http.NotFound(w, r)
        return
    } else if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    result, err := json.Marshal(project)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(result)
}

// Function "Create" creates a new project by json
func Create(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, http.StatusText(405), 405)
        return
    }

    decoder := json.NewDecoder(r.Body)
    project := new(Project)
    err := decoder.Decode(&project)
    if err != nil {
        panic(err)
    }
    defer r.Body.Close()

    if project.Name == "" {
        http.Error(w, http.StatusText(400), 400)
        return
    }

    stmt, err := db.Prepare("INSERT INTO projects(name, pages) VALUES($1, $2);")
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    result, err := stmt.Exec(project.Name, project.Pages)
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

// Function "Update" updates project by json
func Update(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, http.StatusText(405), 405)
        return
    }

    decoder := json.NewDecoder(r.Body)
    project := new(Project)
    err := decoder.Decode(&project)
    if err != nil {
        panic(err)
    }
    defer r.Body.Close()    

    if project.Id == 0 || project.Name == "" || project.Pages == 0 {
        http.Error(w, http.StatusText(400), 400)
        return
    }

    stmt, err := db.Prepare("update projects set name = $1, pages = $2 where id=$3")
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    result, err := stmt.Exec(project.Name, project.Pages, project.Id)
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

// Function "Delete" delete a project by id
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

    deleteData, err := db.Prepare("delete from data where project=$1")
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    deleteProject, err := db.Prepare("delete from projects where id=$1")
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    deleteData.Exec(id)
    projectResult, err := deleteProject.Exec(id)
    if err != nil {
        log.Fatal(err)
    }

    rowsAffected, err := projectResult.RowsAffected()
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    if rowsAffected > 0 {
        fmt.Fprintf(w, "%t\n", true)
    }
}