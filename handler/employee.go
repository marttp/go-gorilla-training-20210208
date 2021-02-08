package handler

import (
	"../model"
	"../util"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

func GetAllEmployees(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	employees := []model.Employee{}
	db.Find(&employees)
	util.RespondJSON(w, http.StatusOK, employees)
}

func CreateEmployee(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	employee := model.Employee{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&employee); err != nil {
		util.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&employee).Error; err != nil {
		util.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.RespondJSON(w, http.StatusCreated, employee)
}

func GetEmployee(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	employee := getEmployeeOr404(db, name, w, r)
	if employee == nil {
		return
	}
	util.RespondJSON(w, http.StatusOK, employee)
}

func UpdateEmployee(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	employee := getEmployeeOr404(db, name, w, r)
	if employee == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&employee); err != nil {
		util.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&employee).Error; err != nil {
		util.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.RespondJSON(w, http.StatusOK, employee)
}

func DeleteEmployee(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	employee := getEmployeeOr404(db, name, w, r)
	if employee == nil {
		return
	}
	if err := db.Delete(&employee).Error; err != nil {
		util.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.RespondJSON(w, http.StatusNoContent, nil)
}

func DisableEmployee(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	employee := getEmployeeOr404(db, name, w, r)
	if employee == nil {
		return
	}
	employee.Disable()
	if err := db.Save(&employee).Error; err != nil {
		util.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.RespondJSON(w, http.StatusOK, employee)
}

func EnableEmployee(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	employee := getEmployeeOr404(db, name, w, r)
	if employee == nil {
		return
	}
	employee.Enable()
	if err := db.Save(&employee).Error; err != nil {
		util.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.RespondJSON(w, http.StatusOK, employee)
}

// getEmployeeOr404 gets a employee instance if exists, or respond the 404 error otherwise
func getEmployeeOr404(db *gorm.DB, name string, w http.ResponseWriter, r *http.Request) *model.Employee {
	employee := model.Employee{}
	if err := db.First(&employee, model.Employee{Name: name}).Error; err != nil {
		util.RespondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &employee
}