package scheduler

import (
	"bytes"
	"fmt"
	"io"
	"micro-site/api/model"
	"micro-site/database"
	"net/http"
)

/*
* Retry failed leads
* @return void
 */
func RetryFailedLeads() {
	var logs []model.LeadExternalLog
	if err := database.DB.Where("status = ?", "Failed").Find(&logs).Error; err != nil {
		fmt.Println("Error fetching failed logs:", err)
		return
	}

	apiUrl := "https://digiweb.aliceblueonline.com/DigiWebsite"
	client := &http.Client{}

	for _, log := range logs {
		// Use the existing request payload
		reqPayload := log.RequestPayload

		// If payload is empty or not valid JSON, we might skip, but for now we try to send it if it exists.
		if reqPayload == "" {
			continue
		}

		resp, err := client.Post(apiUrl, "application/json", bytes.NewBuffer([]byte(reqPayload)))

		log.ApiUrl = apiUrl // Ensure URL is updated just in case

		updatedValues := model.LeadExternalLog{}

		if err != nil {
			updatedValues.Status = "Failed"
			updatedValues.ErrorMessage = err.Error()
			updatedValues.HttpStatusCode = 0
		} else {
			defer resp.Body.Close()
			bodyBytes, _ := io.ReadAll(resp.Body)

			updatedValues.ResponsePayload = string(bodyBytes)
			updatedValues.HttpStatusCode = resp.StatusCode

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				updatedValues.Status = "Success"
				updatedValues.ErrorMessage = "" // Clear error message on success
			} else {
				updatedValues.Status = "Failed"
			}
		}

		// Update the log entry
		database.DB.Model(&log).Updates(updatedValues)
	}
}
