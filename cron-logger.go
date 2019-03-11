package main

import (
	"cron-logger/cmdline"
	"encoding/json"
	"fmt"
	"gopkg.in/resty.v1"
	"log"
	"os"
	"sort"
	"time"
)

type TeamItem struct {
	Name      string `json:"team_name"`
	Occupants uint32 `json:"occupants"`
}

func main() {

	// check config
	config := cmdline.GetSettings()
	if *config.Host == "" || *config.Port == "" || *config.File == "" {
		cmdline.PrintUsage()
		return
	}

	// open file to which we will be writing
	handle, fErr := os.OpenFile(*config.File, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if fErr != nil {
		log.Println("Failed to open log file -> ", *config.File)
		os.Exit(1)
	}

	var message = "::start of log entry@" + time.Now().UTC().String() + "\n"

	// upon logger close write out to log file
	defer func(msg *string) {
		_, err := handle.Write([]byte(*msg))
		if err != nil {
			log.Println("failed to write log file")
		}
		_ = handle.Close() // ignore cloe error
	}(&message)

	// get table from rest api
	resp, err := resty.R().Get("http://" + config.GetAddress() + "/api/v1/teams/occupants")
	if err != nil || !resp.IsSuccess() {
		log.Println("failed to access rest api ...")
		message += fmt.Sprintf("::error - rest call failed, err: %v\n", err)
		os.Exit(1)
	}

	// process result from rest api
	result := make([]TeamItem, 0)
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		log.Println("failed to unmarshall json ...")
		message += "::error - failed to unmarshall REST response -> " + err.Error() + "\n"
		os.Exit(1)
	}

	// process result and add it to output
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	message += fmt.Sprintf("::row count[%d]\n", len(result))

	for _, a := range result {
		message += fmt.Sprintf("%s=%d\n", a.Name, a.Occupants)
	}

	message += "\n"

}
