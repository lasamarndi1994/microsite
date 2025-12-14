package service

import (
	"encoding/base64"
	"os"
	"strings"
)

/*
* Upload Base64 image
* @param imgData string
* @param fileName string
* @param path string
* @return bool
 */
func UploadBase64Image(imgData string, fileName string, path string) bool {
	os.MkdirAll("./uploads/"+path, os.ModePerm)
	// Remove Base64 header part: "data:image/png;base64,"
	if idx := strings.Index(imgData, ";base64,"); idx != -1 {
		imgData = imgData[idx+8:]
	}
	// Decode base64 string
	data, err := base64.StdEncoding.DecodeString(imgData)
	if err != nil {
		return false
	}
	// Create file name
	//fileName := fmt.Sprintf(file_name+"_%d.png", os.Getpid())

	// Write file
	err = os.WriteFile("./uploads/"+path+"/"+fileName, data, 0644)
	if err != nil {
		return false
	} else {
		return true
	}
}

/*
* Delete Image
* @param fileName string
* @param path string
* @return error
 */
func DeleteImage(fileName string, path string) error {
	err := os.Remove("./uploads/" + path + "/" + fileName)
	if err != nil {
		return err
	}
	return nil
}
