package helper

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type FieldErrorResponse struct {
	Status  bool              `json:"status"`
	Message string            `json:"message"`
	Error   map[string]string `json:"errors"`
}

func FormatValidationError(err error) map[string]string {
	errors := map[string]string{}

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			//field := strings.ToLower(fe.Field()) // field name like 'email'
			field := toSnakeCase(fe.Field()) // field name like 'email'
			switch fe.Tag() {
			case "required":
				errors[field] = field + " is required"
			case "email":
				errors[field] = "Invalid email address"
			case "uni_users_email":
				errors[field] = "Email already exists"
			case "gte":
				errors[field] = field + " must be at least " + fe.Param()
			case "lte":
				errors[field] = field + " must be at most " + fe.Param()
			case "passwordmatch":
				errors[field] = "Password and confirm password do not match"
			default:
				errors[field] = "invalid value for " + field
			}
		}
	}

	return errors
}

func DataBaseValidationError1(err error) map[string]string {
	if err == nil || !strings.Contains(err.Error(), "Duplicate entry") {
		return nil
	}

	//define mapping of DB constraint name to user-friendly filed message

	constantMaps := map[string]struct {
		Field   string
		Message string
	}{
		"uni_users_email":         {"email", "Email already registered"},
		"uni_users_employee_id":   {"employee_id", "Employee ID already taken"},
		"uni_users_mobile_number": {"mobile_number", "Mobile number already used"},
		"uni_designations_name":   {"name", "Designation name already exists"},
		"uni_departments_name":    {"name", "Department ID already exists"},
	}

	for constraint, v := range constantMaps {
		if strings.Contains(err.Error(), constraint) {
			//Derive filed name from constraint (Optional logic)
			// filed := strings.TrimPrefix(constraint, "uni_")
			// filed = strings.Split(filed, "_")[1]
			return map[string]string{
				v.Field: v.Message,
			}
		}
	}

	// Default fallback
	return map[string]string{
		"error": err.Error(),
	}
}

func toSnakeCase(str string) string {
	snake := regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}

func DataBaseValidationError(err error) *FieldErrorResponse {

	if err == nil || !strings.Contains(err.Error(), "Duplicate entry") {
		return nil
	}

	// Define mapping of DB constraint names to user-friendly field messages
	constraintMap := map[string]string{
		"uni_users_email":         "Email already registered",
		"uni_users_username":      "Username already taken",
		"uni_users_mobile_number": "Mobile number already used",
		"uni_companies_name":      "Company name already exists",
		"uni_users_employees_id":  "Employee ID already exists",
		"uni_designations_name":   "Designations name is taken",
		"uni_departments_name":    "Department name is taken",
		"uni_levae_types_name":    "The leave is already taken",
	}

	for constraint, message := range constraintMap {
		if strings.Contains(err.Error(), constraint) {
			// Derive field name from constraint (optional logic)
			field := strings.TrimPrefix(constraint, "uni_")
			field1 := strings.Split(field, "_")[1] // e.g., "email" from "uni_users_email"
			//	var field3 string
			if len(strings.Split(constraint, "_")) > 3 {
				field2 := strings.Split(field, "_")[2]
				field1 = field1 + "_" + field2
			}

			return &FieldErrorResponse{
				Status:  false,
				Message: "Failed",
				Error: map[string]string{
					field1: message,
				},
			}
		}
	}
	return &FieldErrorResponse{
		Message: err.Error(),
		Error:   nil,
		Status:  false,
	}
}
