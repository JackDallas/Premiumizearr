package premiumizeme

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Premiumizeme struct {
	APIKey string
}

func NewPremiumizemeClient(APIKey string) Premiumizeme {
	return Premiumizeme{APIKey: APIKey}
}

func (pm *Premiumizeme) createPremiumizemeURL(urlPath string) (url.URL, error) {
	u, err := url.Parse("https://www.premiumize.me/api/")
	if err != nil {
		return *u, err
	}
	u.Path = path.Join(u.Path, urlPath)
	q := u.Query()
	q.Set("apikey", pm.APIKey)
	u.RawQuery = q.Encode()
	return *u, nil
}

func (pm *Premiumizeme) GetTransfers() ([]Transfer, error) {
	log.Trace("Getting transfers list from premiumize.me")
	url, err := pm.createPremiumizemeURL("/transfer/list")
	if err != nil {
		return nil, err
	}

	var ret []Transfer
	req, _ := http.NewRequest("GET", url.String(), nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ret, err
	}

	defer resp.Body.Close()
	res := ListTransfersResponse{}
	err = json.NewDecoder(resp.Body).Decode(&res)

	if res.Status != "success" {
		return ret, fmt.Errorf("%s", res.Status)
	}

	if err != nil {
		return ret, err
	}

	log.Tracef("Received %d transfers", len(res.Transfers))
	return res.Transfers, nil
}

func (pm *Premiumizeme) ListFolder(folderID string) ([]Item, error) {
	var ret []Item
	url, err := pm.createPremiumizemeURL("/folder/list")
	if err != nil {
		return ret, err
	}

	q := url.Query()
	q.Set("id", folderID)
	url.RawQuery = q.Encode()

	client := &http.Client{}
	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return ret, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return ret, err
	}

	if resp.StatusCode != 200 {
		return ret, fmt.Errorf("error listing folder: %s (%d)", resp.Status, resp.StatusCode)
	}

	defer resp.Body.Close()
	res := ListFoldersResponse{}
	log.Trace("Reading response")
	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		return ret, err
	}

	if res.Status != "success" {
		return ret, fmt.Errorf(res.Message)
	}

	return res.Content, nil
}

func (pm *Premiumizeme) GetFolders() ([]Item, error) {
	log.Trace("Getting folder list from premiumize.me")
	url, err := pm.createPremiumizemeURL("/folder/list")
	if err != nil {
		return nil, err
	}

	var ret []Item
	req, _ := http.NewRequest("GET", url.String(), nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ret, err
	}

	defer resp.Body.Close()
	res := ListFoldersResponse{}
	err = json.NewDecoder(resp.Body).Decode(&res)

	if res.Status != "success" {
		return ret, fmt.Errorf("%s", res.Status)
	}

	if err != nil {
		return ret, err
	}

	log.Tracef("Received %d Folders", len(res.Content))
	return res.Content, nil
}

func (pm *Premiumizeme) CreateTransfer(filePath string, parentID string) error {
	//TODO: handle file size, i.e. incorrect file being saved
	log.Trace("Opening file: ", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		log.Errorf("First try failed, waiting 1 second and trying to open file: %s again", filePath)
		time.Sleep(1 * time.Second)
		file, err = os.Open(filePath)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	url, err := pm.createPremiumizemeURL("/transfer/create")
	if err != nil {
		return err
	}

	client := &http.Client{}
	var request *http.Request

	switch filepath.Ext(file.Name()) {
	case ".nzb":
		request, err = createNZBRequest(file, &url, parentID)
	case ".magnet":
		request, err = createMagnetRequest(file, &url, parentID)
	}

	if err != nil {
		return err
	}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error creating transfer: %s (%d)", resp.Status, resp.StatusCode)
	}

	defer resp.Body.Close()
	res := CreateTransferResponse{}
	log.Trace("Reading response")
	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		return err
	}

	if res.Status != "success" {
		return fmt.Errorf(res.Message)
	}

	log.Tracef("Transfer created: %+v", res)

	return nil
}

func (pm *Premiumizeme) DeleteFolder(folderID string) error {
	url, err := pm.createPremiumizemeURL("/folder/delete")
	if err != nil {
		return err
	}

	q := url.Query()
	q.Set("id", folderID)
	url.RawQuery = q.Encode()

	client := &http.Client{}
	request, err := http.NewRequest("DELETE", url.String(), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error deleting folder: %s (%d)", resp.Status, resp.StatusCode)
	}

	defer resp.Body.Close()
	res := SimpleResponse{}
	log.Trace("Reading response")
	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		return err
	}

	if res.Status != "success" {
		return fmt.Errorf(res.Message)
	}

	log.Tracef("Folder deleted: %+v", res)

	return nil
}

func (pm *Premiumizeme) CreateFolder(folderName string) (string, error) {
	url, err := pm.createPremiumizemeURL("/folder/create")
	if err != nil {
		return "", err
	}

	q := url.Query()
	q.Set("name", folderName)
	url.RawQuery = q.Encode()

	client := &http.Client{}
	request, err := http.NewRequest("POST", url.String(), nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("error creating folder: %s (%d)", resp.Status, resp.StatusCode)
	}

	defer resp.Body.Close()
	res := CreateFolderResponse{}
	log.Trace("Reading response")
	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		return "", err
	}

	if res.Status != "success" {
		return "", fmt.Errorf(res.Message)
	}

	log.Tracef("Folder created: %+v", res)

	return res.ID, nil
}

func (pm *Premiumizeme) DeleteTransfer(id string) error {
	url, err := pm.createPremiumizemeURL("/transfer/delete")
	if err != nil {
		return err
	}

	client := &http.Client{}
	request, err := createDeleteRequest(id, &url)
	if err != nil {
		return err
	}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error deleting transfer: %s (%d)", resp.Status, resp.StatusCode)
	}

	defer resp.Body.Close()
	res := SimpleResponse{}
	log.Trace("Reading response")
	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		return err
	}

	if res.Status != "success" {
		return fmt.Errorf("failed to delete transfer: %s, message: %+v", id, res.Message)
	}

	log.Tracef("Transfer Deleted: %+v", res)

	return nil
}

func createNZBRequest(file *os.File, url *url.URL, parentID string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("src", filepath.Base(file.Name()))

	if err != nil {
		return nil, err
	}

	io.Copy(part, file)
	writer.Close()

	part, err = writer.CreateFormField("folder_id")

	if err != nil {
		return nil, err
	}

	_, err = part.Write([]byte(parentID))

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url.String(), body)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	if err != nil {
		return nil, err
	}

	return request, nil
}

func createMagnetRequest(file *os.File, url *url.URL, parentID string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormField("src")

	if err != nil {
		return nil, err
	}

	io.Copy(part, file)
	writer.Close()

	part, err = writer.CreateFormField("folder_id")

	if err != nil {
		return nil, err
	}

	_, err = part.Write([]byte(parentID))

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url.String(), body)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	if err != nil {
		return nil, err
	}

	return request, nil
}

func createDeleteRequest(id string, URL *url.URL) (*http.Request, error) {
	// Build Values to send to endpoint
	data := url.Values{}
	data.Set("id", id)

	// Create and encode request
	request, err := http.NewRequest("POST", URL.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	// Setup headers
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	return request, nil
}

type SRCType = int

const (
	SRC_FILE = iota
	SRC_FOLDER
)

func (pm *Premiumizeme) GenerateZippedFileLink(fileID string) (string, error) {
	dlLink, err := pm.generateZip(fileID, SRC_FILE)
	if err != nil {
		return "", err
	}
	return dlLink, nil
}

func (pm *Premiumizeme) GenerateZippedFolderLink(fileID string) (string, error) {
	dlLink, err := pm.generateZip(fileID, SRC_FOLDER)
	if err != nil {
		return "", err
	}
	return dlLink, nil
}

func (pm *Premiumizeme) generateZip(ID string, srcType SRCType) (string, error) {
	// Build URL with apikey
	URL, err := pm.createPremiumizemeURL("/zip/generate")
	if err != nil {
		return "", err
	}

	// Build Values to send to endpoint
	data := url.Values{}

	if srcType == SRC_FILE {
		data.Set("files[]", ID)
	} else if srcType == SRC_FOLDER {
		data.Set("folders[]", ID)
	} else {
		return "", fmt.Errorf("unknown source type: %d", srcType)
	}

	// Create and encode request
	request, err := http.NewRequest("POST", URL.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	// Setup headers
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	//Fire request
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("error getting zip link for: %s response: %s (%d)", ID, resp.Status, resp.StatusCode)
	}

	// Decode response
	defer resp.Body.Close()
	var res GenerateZipResponse
	log.Trace("Reading response")
	err = json.NewDecoder(resp.Body).Decode(&res)

	log.Tracef("Zip Response: %+v", res)
	if err != nil {
		return "", err
	}

	if res.Status != "success" {
		return "", fmt.Errorf("error getting zip link for: %s, Status: %s", ID, res.Status)
	}

	log.Debugf("Zip link created: %+v", res.Location)

	return res.Location, nil
}
