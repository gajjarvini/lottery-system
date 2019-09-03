package main

import (
	"encoding/json"
	"io/ioutil"
	log "logger"
)

const (
	configFile = "config/conf.json"
)

/*GetConfigFileValues is used for parse the values from config file*/
func GetConfigFileValues() error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Error.Println("Unable to Fetch Configuration File due to Error: ", err)
		return err
	}

	//	var dataObj configuration
	err = json.Unmarshal(data, &configvalues)
	if err != nil {
		log.Error.Println("Unable to Fetch Config File Data. Please Check file data format. Error: ", err)
		return err
	}
	return nil
}
