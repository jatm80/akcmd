package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/interact"
)

type artifactInterface struct {
	Id         int
	Name       string
	Url        string
	Created_at time.Time
}

type artifactPayloadInterface struct {
	Total_count int
	Artifacts   []artifactInterface
}

func main() {

	artifactsRes := fetch("https://api.github.com/repos/ovrclk/akcmd/actions/artifacts")
	artifacts := toJSON(artifactsRes)

	win_arm64 := getLatestArtifact(getArtifactByName("akcmd-windows-latest_arm64", artifacts.Artifacts))
	win_amd64 := getLatestArtifact(getArtifactByName("akcmd-windows-latest_amd64", artifacts.Artifacts))
	mac_arm64 := getLatestArtifact(getArtifactByName("akcmd-macos-latest_arm64", artifacts.Artifacts))
	mac_amd64 := getLatestArtifact(getArtifactByName("akcmd-macos-latest_amd64", artifacts.Artifacts))
	linux_amd64 := getLatestArtifact(getArtifactByName("akcmd-ubuntu-latest_amd64", artifacts.Artifacts))
	linux_arm64 := getLatestArtifact(getArtifactByName("akcmd-ubuntu-latest_arm64", artifacts.Artifacts))

	downloadURIs := map[string]string{
		"linux_amd64":   linux_amd64.Url + "/zip",
		"linux_arm64":   linux_arm64.Url + "/zip",
		"mac_amd64":     mac_amd64.Url + "/zip",
		"mac_arm64":     mac_arm64.Url + "/zip",
		"windows_amd64": win_arm64.Url + "/zip",
		"windows_arm64": win_amd64.Url + "/zip",
	}

	downloadKey := ""
	osType := runtime.GOOS
	switch osType {
	case "windows":
		downloadKey = "windows"
	case "darwin":
		downloadKey = "mac"
	case "linux":
		downloadKey = "linux"
	}

	plat := runtime.GOARCH

	switch plat {
	case "arm64":
		downloadKey = downloadKey + "_arm64"
	case "amd64":
		downloadKey = downloadKey + "_amd64"
	}

	downloadUrl := downloadURIs[downloadKey]

	color.Redp(`
 AKASH COMMAND CENTER INSTALLER

          ///////////         
            //////////        
             ///////////      
               //////////     
    /////////// //////////,   
   //////////.    //////////  
 ///////////       ////////// 
 /////////***********/////////
  ///////   **********/////// 
    ///*     ***********///   
     /         **********/` + "\n\n")

	// color.Info.Println("Be\n\n")
	// select default network
	_ = interact.SelectOne(
		"Continue installing Akash Command Center",
		[]string{
			"Yes",
		},
		"0",
	)

	fmt.Println("")

	homeDirectory := getHome()
	installDirectory := homeDirectory + "/.bin"

	_ = interact.SelectOne(
		"Install location ",
		[]string{
			installDirectory + "/akcmd",
			"Custom",
		},
		"0",
	)

	fmt.Println("")

	_ = interact.SelectOne(
		"Install akcmd for platform",
		[]string{
			runtime.GOOS + "_" + runtime.GOARCH,
			"No",
		},
		"0",
	)

	fmt.Println("")

	_ = interact.SelectOne(
		"Download version",
		[]string{
			"Nightly",
		},
		"0",
	)

	fmt.Println("")

	err := DownloadFile("akcmd_latest.zip", downloadUrl)
	if err != nil {
		panic(err)
	}

	fmt.Println("Downloading: " + downloadUrl)

	_, err = Unzip("akcmd_latest.zip", installDirectory)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chmod(installDirectory+"/akcmd", 0700)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(`
	
	Installation complete, ensure that ` + installDirectory + ` is included in your $PATH, your can run

	export PATH=$PATH:` + installDirectory + `

	add this line to the end of your .bash_profile, .bashrc, .zshrc
	`)
}

// find the artifacts with the name in the list
func getArtifactByName(name string, artifacts []artifactInterface) []artifactInterface {
	returnArtifacts := []artifactInterface{}
	for i := range artifacts {
		artifact := artifacts[i]
		if artifact.Name == name {
			returnArtifacts = append(returnArtifacts, artifact)
		}
	}
	return returnArtifacts
}

// given a list of artifacts get the latest one
func getLatestArtifact(artifacts []artifactInterface) artifactInterface {
	returnArtifact := artifacts[0]
	for i := range artifacts {
		artifact := artifacts[i]
		if returnArtifact.Created_at.Before(artifact.Created_at) {
			returnArtifact = artifact
		}
	}
	return returnArtifact
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	color.Cyan.Println("Downloading akcmd binary")

	// get the binary from our builds this token can only read public repos
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", `Basic ZG1pa2V5OmdocF9kMXphdEs1VUtzamdaRWNUMmt5VW9HVk45dU5wa0YzaWdabFo=`)
	req.Header.Add("Accept", `application/vnd.github.v3+json`)
	client := &http.Client{}
	resp, _ := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func fetch(url string) string {
	resp, err := http.Get(url)
	if err == nil {
		defer resp.Body.Close()
		responseData, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			responseString := string(responseData)
			return responseString
		}
	}
	return ""
}

func toJSON(jsonString string) artifactPayloadInterface {
	var result artifactPayloadInterface
	json.Unmarshal([]byte(jsonString), &result)
	return result
}

func getHome() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return dirname
}
