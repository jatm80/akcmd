package config

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type NetworkRegistry struct {
	Name      string
	ShortPath string
	Host      string
}

var NetworkRegisteries = []NetworkRegistry{
	{Name: "edgenet", ShortPath: "ovrclk/net/master/edgenet/meta.json", Host: "https://raw.githubusercontent.com/"},
	{Name: "testnet", ShortPath: "ovrclk/net/master/testnet/meta.json", Host: "https://raw.githubusercontent.com/"},
	{Name: "mainnet", ShortPath: "ovrclk/net/master/mainnet/meta.json", Host: "https://raw.githubusercontent.com/"},
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

func GetNetworkFromRegistry(name string) NetworkRegistry {
	for i := range NetworkRegisteries {
		registry := NetworkRegisteries[i]
		if registry.Name == name {
			return registry
		}
	}
	return NetworkRegistry{Name: "undefined"}
}

func GetNetworkInfo() {
	network := GetNetworkFromRegistry("testnet")

	networkMetaUrl := network.Host + network.ShortPath
	fileUrl := networkMetaUrl

	networkMetaData := fetch(fileUrl)

	fmt.Printf("%s", networkMetaData)
}
