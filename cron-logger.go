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
	Name    string `json:"team_name"`
	Members uint32 `json:"members"`
}

func main() {

	config := cmdline.GetSettings()
	if *config.Host == "" || *config.Port == "" || *config.File == "" {
		cmdline.PrintUsage()
		return
	}

	handle, ferr := os.OpenFile(*config.File, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if ferr != nil {
		log.Println("Failed to open log file -> ", *config.File)
	}

	var message = time.Now().UTC().String() + "\n"

	defer func(msg *string) {
		_, err := handle.Write([]byte(message))
		if err != nil {
			log.Println("failed to write log file")
		}
		_ = handle.Close() // ignore cloe error
	}(&message)

	resp, err := resty.R().Get("http://" + config.GetAddress() + "/api/v1/teams/occupants")
	if err != nil || !resp.IsSuccess() {
		log.Println("failed to access rest api ...")
		message += fmt.Sprintf("rest call result: %v, err: %v\n", resp, err)
		return
	}

	result := make([]TeamItem, 0)
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		log.Println("failed to unmarshall json ...")
		message += "failed to unmarshall REST response -> " + err.Error() + "\n"
		return
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	for _, a := range result {
		message += fmt.Sprintf("team name: %16s,  number of occupants: %d\n", a.Name, a.Members)
	}
	message += "\n"

}
