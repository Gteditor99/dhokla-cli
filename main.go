package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

/*               dhokla-cli
the all-in-one cli tool for dhokla (dot) net
*/

type File struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Dir    bool   `json:"dir"`
	Size   int64  `json:"size,omitempty"`
	Parent string `json:"parent"`
	Mtime  string `json:"mtime"`
}

type Directory struct {
	Message string `json:"message"`
	Data    struct {
		Id     string `json:"id"`
		Name   string `json:"name"`
		Dir    bool   `json:"dir"`
		Parent string `json:"parent"`
		Mtime  string `json:"mtime"`
		Files  []File `json:"files"`
	} `json:"data"`
}

// func getDirectory(id string) (*Directory, error) {
// 	// Make a request to retrieve the directory with the given ID
// 	// and store the response in jsonData
// 	jsonData, err := http.Get("https://dhokla.net/api/directories/" + id)
// 	if err != nil {
// 		return nil, fmt.Errorf("error occurred during request: %v", err)
// 	}

// 	defer jsonData.Body.Close()

// 	body, err := io.ReadAll(jsonData.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("error occurred while reading response body: %v", err)
// 	}

// 	// Unmarshal the JSON data into a Directory struct
// 	var dir Directory
// 	err = json.Unmarshal(body, &dir)
// 	if err != nil {
// 		return nil, fmt.Errorf("error occurred while unmarshalling JSON data: %v", err)
// 	}
// 	return &dir, nil
// }

func getDirectory(id string) (map[string]File, error) {
	// Make a request to retrieve the directory with the given ID
	// and store the response in jsonData
	jsonData, err := http.Get("https://dhokla.net/api/directories/" + id)
	if err != nil {
		return nil, fmt.Errorf("error occurred during request: %v", err)
	}

	defer jsonData.Body.Close()

	body, err := io.ReadAll(jsonData.Body)
	if err != nil {
		return nil, fmt.Errorf("error occurred while reading response body: %v", err)
	}

	// Unmarshal the JSON data into a Directory struct
	var dir Directory
	err = json.Unmarshal(body, &dir)
	if err != nil {
		return nil, fmt.Errorf("error occurred while unmarshalling JSON data: %v", err)
	}
	// fileMap impl for id referencing
	fileMap := make(map[string]File)

	for _, file := range dir.Data.Files {
		fileMap[file.Id] = file
	}
	return fileMap, nil
}

// curl POST "https://dhokla.net/files/b5c7a659-f63c-426e-9f7b-74650ad8a943" >> output.txt
func downloadFile(f File) error {
	fmt.Printf("Started downloading %v, ID: %v, size: %v\n", f.Name, f.Id, f.Size)
	URL := "https.dhokla.net/files/" + f.Id

	// check if URL is up
	fileData, err := http.Get(URL)
	if err != nil {
		return fmt.Errorf(
			"error occurred while downloading %v, ID: %v, size: %v",
			f.Name,
			f.Id,
			f.Size,
		)
	}

	// create file
	out, err := os.Create(f.Name)
	if err != nil {
		return err
	}
	// stream data to file
	defer fileData.Body.Close()

	_, err = io.Copy(out, fileData.Body)
	return err
}

func main() {
	// Get the directory with the given ID
	fmt.Printf("dhokla-cli: Enter directory ID:")
	fmt.Scan()
	var dirID string
	fmt.Scanln(&dirID)

	fileMap, err := getDirectory(dirID)
	if err != nil {
		log.Fatalf("error occurred while getting directory: %v", err)
	}

	// Print out the directory information
	// fmt.Printf("Directory Name: %s\n", dir.Data.Name)
	// fmt.Printf("Directory ID: %s\n", dir.Data.Id)
	// fmt.Printf("Is Directory: %v\n", dir.Data.Dir)
	//  fmt.Printf("Message: %s\n", dir.Message)
	//  fmt.Printf("Parent Directory ID: %s\n", dir.Data.Parent)
	//  fmt.Printf("Modification Time: %s\n", dir.Data.Mtime)

	for _, file := range fileMap {
		fmt.Printf(
			"Name: %s, ID: %s, Size: %d, Modification Time: %s\n",
			file.Name,
			file.Id,
			file.Size,
			file.Mtime,
		)
	}
	file, ok := fileMap["96a6d6e2-583b-4c08-9598-d2849c354856"]
	if ok {
		fmt.Printf("%v\n", file.Id)
	} else {
		fmt.Println("No file found with the given ID")
	}

	// downloadFile(fileMap["96a6d6e2-583b-4c08-9598-d2849c354856"])
}
