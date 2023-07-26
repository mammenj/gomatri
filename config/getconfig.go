package config

import (
	"encoding/json"
	"log"
	"os"
)

// Configuration struct for application config
type Configuration struct {
	Serverhost string `json:"serverhost"`
	Serverport string `json:"serverport"`
	Dbuser     string `json:"dbuser"`
	Dbport     string `json:"dbport"`
	Dbpassword string `json:"dbpassword"`
	Db         string `json:"db"`
	Dbhost     string `json:"dbhost"`
	Dbengine   string `json:"dbengine"`
	Dbssl      string `json:"dbssl"`
}

// GetConfiguration gets al the application configs
func GetConfiguration(configfile string) (Configuration, error) {
	myconfig := Configuration{}

	file, err := os.Open(configfile)
	if err != nil {
		return myconfig, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&myconfig)
	if err != nil {
		log.Printf("GetConfiguration Errror Return Config, Error %v, %v", myconfig, err)
		return myconfig, err
	}
	return myconfig, nil
}
