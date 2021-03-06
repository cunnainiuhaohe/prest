package controllers

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/nuveo/prest/adapters/postgres"
	"github.com/nuveo/prest/api"
	"github.com/nuveo/prest/statements"
)

// GetTables list all (or filter) tables
func GetTables(w http.ResponseWriter, r *http.Request) {
	requestWhere := postgres.WhereByRequest(r)
	sqlTables := statements.Tables
	if requestWhere != "" {
		sqlTables = fmt.Sprint(
			statements.TablesSelect,
			statements.TablesWhere,
			" AND ",
			requestWhere,
			statements.TablesOrderBy)
	}

	object, err := postgres.Query(sqlTables)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(object)
}

// GetTablesByDatabaseAndSchema list all (or filter) tables based on database and schema
func GetTablesByDatabaseAndSchema(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	database, ok := vars["database"]
	if !ok {
		log.Println("Unable to parse database in URI")
		http.Error(w, "Unable to parse database in URI", http.StatusInternalServerError)
		return
	}
	schema, ok := vars["schema"]
	if !ok {
		log.Println("Unable to parse schema in URI")
		http.Error(w, "Unable to parse schema in URI", http.StatusInternalServerError)
		return
	}
	requestWhere := postgres.WhereByRequest(r)
	sqlSchemaTables := statements.SchemaTables
	if requestWhere != "" {
		sqlSchemaTables = fmt.Sprint(
			statements.SchemaTablesSelect,
			statements.SchemaTablesWhere,
			" AND ",
			requestWhere,
			statements.SchemaTablesOrderBy)
	}
	sqlSchemaTables = fmt.Sprint(sqlSchemaTables, " ", postgres.PaginateIfPossible(r))

	object, err := postgres.Query(sqlSchemaTables, database, schema)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(object)
}

// SelectFromTables perform select in database
func SelectFromTables(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	database, ok := vars["database"]
	if !ok {
		log.Println("Unable to parse database in URI")
		http.Error(w, "Unable to parse database in URI", http.StatusInternalServerError)
		return
	}
	schema, ok := vars["schema"]
	if !ok {
		log.Println("Unable to parse schema in URI")
		http.Error(w, "Unable to parse schema in URI", http.StatusInternalServerError)
		return
	}
	table, ok := vars["table"]
	if !ok {
		log.Println("Unable to parse table in URI")
		http.Error(w, "Unable to parse table in URI", http.StatusInternalServerError)
		return
	}

	query := fmt.Sprintf("%s %s.%s.%s", statements.SelectInTable, database, schema, table)
	requestWhere := postgres.WhereByRequest(r)
	sqlSelect := query
	if requestWhere != "" {
		sqlSelect = fmt.Sprint(
			query,
			" WHERE ",
			requestWhere)
	}
	sqlSelect = fmt.Sprint(sqlSelect, " ", postgres.PaginateIfPossible(r))

	object, err := postgres.Query(sqlSelect)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(object)
}

// InsertInTables perform insert in specific table
func InsertInTables(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	database, ok := vars["database"]
	if !ok {
		log.Println("Unable to parse database in URI")
		http.Error(w, "Unable to parse database in URI", http.StatusInternalServerError)
		return
	}
	schema, ok := vars["schema"]
	if !ok {
		log.Println("Unable to parse schema in URI")
		http.Error(w, "Unable to parse schema in URI", http.StatusInternalServerError)
		return
	}
	table, ok := vars["table"]
	if !ok {
		log.Println("Unable to parse table in URI")
		http.Error(w, "Unable to parse table in URI", http.StatusInternalServerError)
		return
	}
	req := api.Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	object, err := postgres.Insert(database, schema, table, req)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(object)
}

// DeleteFromTable perform delete sql
func DeleteFromTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	database, ok := vars["database"]
	if !ok {
		log.Println("Unable to parse database in URI")
		http.Error(w, "Unable to parse database in URI", http.StatusInternalServerError)
		return
	}
	schema, ok := vars["schema"]
	if !ok {
		log.Println("Unable to parse schema in URI")
		http.Error(w, "Unable to parse schema in URI", http.StatusInternalServerError)
		return
	}
	table, ok := vars["table"]
	if !ok {
		log.Println("Unable to parse table in URI")
		http.Error(w, "Unable to parse table in URI", http.StatusInternalServerError)
		return
	}

	where := postgres.WhereByRequest(r)

	object, err := postgres.Delete(database, schema, table, where)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(object)
}

// UpdateTable perform update table
func UpdateTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	database, ok := vars["database"]
	if !ok {
		log.Println("Unable to parse database in URI")
		http.Error(w, "Unable to parse database in URI", http.StatusInternalServerError)
		return
	}
	schema, ok := vars["schema"]
	if !ok {
		log.Println("Unable to parse schema in URI")
		http.Error(w, "Unable to parse schema in URI", http.StatusInternalServerError)
		return
	}
	table, ok := vars["table"]
	if !ok {
		log.Println("Unable to parse table in URI")
		http.Error(w, "Unable to parse table in URI", http.StatusInternalServerError)
		return
	}

	where := postgres.WhereByRequest(r)
	req := api.Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	object, err := postgres.Update(database, schema, table, where, req)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(object)
}
