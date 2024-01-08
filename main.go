package main

import (
	"encoding/json"
	"fmt"
	"io"
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

func getRootDirectory() (map[string]File, error) {
	// Retrieves the root directory of dhokla.net (note the lack of an ID in the request)
	jsonData, err := http.Get("https://dhokla.net/api/d/")
	if err != nil {
		return nil, fmt.Errorf("error occurred during request: %v", err)
	}

	// The root directory is special, since it is used to categorize media based on type.
	// movies, tv shows, anime, etc.
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

func getDirectory(id string) (map[string]File, error) {
	/* Makes a GET request to retrieve specific directories using id. */
	jsonData, err := http.Get("https://dhokla.net/api/d/" + id)
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


// keep-alive connection for parallel downloads
func downloadFile(f File) error {
	fmt.Printf("Started downloading %v, ID: %v, size: %v\n", f.Name, f.Id, f.Size)
	URL := "https.dhokla.net/f/" + f.Id

	

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
	// keep-alive connection for parallel downloads
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Connection", "keep-alive")

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


type SearchResponse struct {
	Message string `json:"message"`
	Data    []File `json:"data"`
}

// search for files/directories by name
func search(name string) ([]File, error) {
	jsonData, err := http.Get("https://dhokla.net/api/s/" + name)
	if err != nil {
		return nil, fmt.Errorf("error occurred during request: %v", err)
	}

	defer jsonData.Body.Close()

	body, err := io.ReadAll(jsonData.Body)
	if err != nil {
		return nil, fmt.Errorf("error occurred while reading response body: %v", err)
	}

	// Unmarshal the JSON data into a SearchResponse struct
	var response SearchResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error occurred while unmarshalling JSON data: %v", err)
	}

	return response.Data, nil
}

func main() {
	// [TEST] Print out the files in the root directory
	rootFileDirectories, err := getRootDirectory()
	if err != nil {
		fmt.Printf("error occurred while getting root directory: %v", err)
	}
	for _, file := range rootFileDirectories {
		fmt.Printf(file.Id + " " + file.Name + "\n")
	}

	// [TEST] Print out the files in the directory with ID d0404f70-93ff-4059-984f-9b988941f955
	dir, err := getDirectory("d0404f70-93ff-4059-984f-9b988941f955")
	if err != nil {
		fmt.Printf("error occurred while getting directory: %v", err)
	}
	for _, file := range dir {
		fmt.Printf(file.Id + " " + file.Name + "\n")
	}

	// [TEST] search for files/directories by name
	var searchQuery string
	fmt.Printf("Enter search query: ")
	_, err = fmt.Scanln(&searchQuery)
	if err != nil {
		fmt.Printf("error occurred while scanning input: %v", err)
	}

	searchResults, err := search("'" + searchQuery + "'")
	if err != nil {
		fmt.Printf("error occurred while searching: %v", err)
	}
	for _, file := range searchResults {
		fmt.Printf(file.Id + " " + file.Name + "\n")
	}
}
