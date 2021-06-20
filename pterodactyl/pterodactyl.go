package pterodactyl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gamemann/Pterodactyl-Game-Server-Watch/config"
)

// Attributes struct from /api/client/servers/xxxx/resources.
type Attributes struct {
	State string `json:"current_state"`
}

// Utilization struct from /api/client/servers/xxxx/resources.
type Utilization struct {
	Attributes Attributes `json:"attributes"`
}

// Retrieves all servers and add them to the config.
func AddServers(cfg *config.Config) bool {
	// Build endpoint.
	urlstr := cfg.APIURL + "/" + "api/client"

	// Setup HTTP GET request.
	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("GET", urlstr, nil)

	// Set authorization header.
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	// Accept only JSON.
	req.Header.Set("Accept", "application/json")

	// Perform HTTP request and check for errors.
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)

		return false
	}

	// Close body at the end.
	defer resp.Body.Close()

	// Read body.
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)

		return false
	}

	// Create utilization struct.
	var dataobj interface{}

	// Parse JSON.
	err = json.Unmarshal([]byte(string(body)), &dataobj)

	if err != nil {
		fmt.Println(err)

		return false
	}

	// Loop through each data item (server).
	for _, j := range dataobj.(map[string]interface{})["data"].([]interface{}) {
		item := j.(map[string]interface{})

		// Make sure we have a server object.
		if item["object"] == "server" {
			attr := item["attributes"].(map[string]interface{})

			// Build new server structure.
			var sta config.Server

			sta.Enable = true
			sta.UID = attr["identifier"].(string)
			sta.ScanTime = 5
			sta.MaxFails = 10
			sta.MaxRestarts = 2
			sta.RestartInt = 120

			// Retrieve default IP/port.
			for _, i := range attr["relationships"].(map[string]interface{})["allocations"].(map[string]interface{})["data"].([]interface{}) {
				if i.(map[string]interface{})["object"].(string) != "allocation" {
					continue
				}

				alloc := i.(map[string]interface{})["attributes"].(map[string]interface{})

				if alloc["is_default"].(bool) {
					sta.IP = alloc["ip"].(string)
					sta.Port = int(alloc["port"].(float64))
				}
			}

			// Append to servers slice.
			cfg.Servers = append(cfg.Servers, sta)

			fmt.Println("[API] Adding server " + sta.IP + ":" + strconv.Itoa(sta.Port) + " with UID " + sta.UID)
		}
	}

	// Otherwise, return true meaning the container is online.
	return true
}

// Checks the status of a Pterodactyl server. Returns true if on and false if off.
func CheckStatus(apiURL string, apiToken string, uid string) bool {
	// Build endpoint.
	urlstr := apiURL + "/" + "api/client/servers/" + uid + "/resources"

	// Setup HTTP GET request.
	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("GET", urlstr, nil)

	// Set authorization header.
	req.Header.Set("Authorization", "Bearer "+apiToken)

	// Set data to JSON.
	req.Header.Set("Content-Type", "application/json")

	// Accept JSON.
	req.Header.Set("Accept", "application/json")

	// Perform HTTP request and check for errors.
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)

		return false
	}

	// Close body at the end.
	defer resp.Body.Close()

	// Read body.
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)

		return false
	}

	// Create utilization struct.
	var util Utilization

	// Parse JSON.
	json.Unmarshal([]byte(string(body)), &util)

	// Check if the server's state isn't on. If not, return false.
	if util.Attributes.State != "running" {
		return false
	}

	// Otherwise, return true meaning the container is online.
	return true
}

// Kills the specified server.
func KillServer(apiURL string, apiToken string, uid string) {
	// Build endpoint.
	urlstr := apiURL + "/" + "api/client/servers/" + uid + "/" + "power"

	// Setup form data.
	var formdata = []byte(`{"signal": "kill"}`)

	// Setup HTTP GET request.
	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("POST", urlstr, bytes.NewBuffer(formdata))

	// Set authorization header.
	req.Header.Set("Authorization", "Bearer "+apiToken)

	// Set data to JSON.
	req.Header.Set("Content-Type", "application/json")

	// Accept JSON.
	req.Header.Set("Accept", "application/json")

	// Perform HTTP request and check for errors.
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	// Close body at the end.
	resp.Body.Close()
}

// Starts the specified server.
func StartServer(apiURL string, apiToken string, uid string) {
	// Build endpoint.
	urlstr := apiURL + "/" + "api/client/servers/" + uid + "/" + "power"

	// Setup form data.
	var formdata = []byte(`{"signal": "start"}`)

	// Setup HTTP GET request.
	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("POST", urlstr, bytes.NewBuffer(formdata))

	// Set authorization header.
	req.Header.Set("Authorization", "Bearer "+apiToken)

	// Set data to JSON.
	req.Header.Set("Content-Type", "application/json")

	// Accept JSON.
	req.Header.Set("Accept", "application/json")

	// Perform HTTP request and check for errors.
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	// Close body at the end.
	resp.Body.Close()
}
