// Copyright 2021 ThingWave AB.  All rights reserved.
package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	MAJOR int = 0
	MINOR int = 1
	REV   int = 3
)

func main() {
	os.Exit(mainApp())
}

func mainApp() int {
	version := flag.Bool("version", false, "Prints the version")
	command := flag.String("cmd", "sr-echo", "Eclipse Arrowhead command to execute")
	targetUri := flag.String("sr", "https://127.0.0.1:8443/serviceregistry", "Service Registry URI")
	ca := flag.String("cafile", "", "Root CA PEM file")
	cert := flag.String("cert", "", "Client certificate PEM file")
	key := flag.String("key", "", "Client certificate key PEM file")
	verbose := flag.Bool("verbose", false, "Makes client more verbose whens set to true")
	flag.Parse()

	if *version == true {
		fmt.Printf("%d.%d.%d\n", MAJOR, MINOR, REV)
		return 0
	}

	_, err := url.ParseRequestURI(*targetUri)
	if err != nil {
		fmt.Printf("Illegal URI: %s\n", *targetUri)
		return -1
	}

	// Checking for missing arguments
	if (*ca != "" && (*cert == "" || *key == "")) || (*cert != "" && (*ca == "" || *key == "")) || (*key != "" && (*cert == "" || *ca == "")) {
		fmt.Println("Error: missing arguments!")
		fmt.Println("Use 'ahctl --help' to see usage.")
		return -1
	}

	if *command == "sr-echo" || *command == "or-echo" || *command == "au-echo" || *command == "dm-echo" {

	} else if *command == "get-all-systems" {
		// This will get all the systems available

	} else if *command == "get-all-services" {
		// Get all the available services

	} else if *command == "get-grouped" {

	} else {
		// Command received is unknown, return -1 to mark an error
		fmt.Printf("Unknown command: %s\n", *command)
		return -1
	}

	// create a new client variable "client"
	var client http.Client

	// Check the files of certificate etc. If not, return -1. Else run in unsecure mode
	if *ca != "" && *cert != "" && *key != "" {
		ok, _ := fileExists(*ca)
		if !ok {
			fmt.Printf("Error: ca file %s does not exist\n", *ca)
			return -1
		}
		ok, _ = fileExists(*cert)
		if !ok {
			fmt.Printf("Error: certificate file %s does not exist\n", *cert)
			return -1
		}
		ok, _ = fileExists(*key)
		if !ok {
			fmt.Printf("Error: key file %s does not exist\n", *key)
			return -1
		}
		// Read the files after making sure they exist
		caCert, err := ioutil.ReadFile(*ca)
		if err != nil {
			log.Fatalf("Error opening cert file %s, Error: %s", *ca, err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		myCert, err := tls.LoadX509KeyPair(*cert, *key)
		if err != nil {
			log.Fatalf("Error creating x509 keypair from client cert file %s and client key file %s", *cert, *key)
		}

		t := &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{myCert},
				RootCAs:            caCertPool,
				InsecureSkipVerify: false,
			},
		}

		client = http.Client{Transport: t, Timeout: 5 * time.Second}

		if *command == "srecho" {
			fmt.Printf("Calling %s\n", *targetUri+"/echo")
			data, err := getData(client, *targetUri+"/echo")
			if err == nil {
				fmt.Println(string(data))
			}
		}
	} else {
		//fmt.Println("Running unsecure mode ...")
		client = http.Client{Timeout: 10 * time.Second}
	}

	if *command == "sr-echo" { // GET to serviceregistry/echo
		fmt.Printf("Calling %s\n", *targetUri+"/echo")
		data, err := getData(client, *targetUri+"/echo")
		if err != nil {
			fmt.Printf("Could not connect to '%s'\n", *targetUri+"/echo")
			return -1
		} else {
			fmt.Println(string(data))
		}
	} else if *command == "get-all-systems" {
		data, err := getData(client, *targetUri+"/mgmt/systems?direction=ASC&sort_field=id")
		if err == nil {
			var response SystemList
			json.Unmarshal(data, &response)

			empJSON, _ := json.MarshalIndent(response, "", "  ")
			fmt.Println(string(empJSON))
		}
	} else if *command == "get-all-services" {
		data, err := getData(client, *targetUri+"/mgmt/services?direction=ASC&sort_field=id")
		if err == nil {
			//fmt.Println(data)

			var response ServiceDefinitionList
			json.Unmarshal(data, &response)

			empJSON, _ := json.MarshalIndent(response, "", "  ")
			fmt.Println(string(empJSON))
		}
	} else if *command == "get-grouped" {
		data, err := getData(client, *targetUri+"/mgmt/grouped")
		if err == nil {
			//fmt.Println(string(data))

			var response GroupedDTO
			json.Unmarshal(data, &response)

			empJSON, _ := json.MarshalIndent(response, "", "  ")
			fmt.Println(string(empJSON))
		}
	} else if *command == "or-echo" || *command == "au-echo" || *command == "dm-echo" { // GET to orchestrator/echo
		//} else if *command == "au-echo" { // GET to authorization/echo
		//} else if *command == "dm-echo" { // GET to datamanager/echo
		var sreq ServiceQueryRequest
		if *command == "or-echo" {
			sreq.ServiceDefinitionRequirement = "proxy" //XXX
		} else if *command == "au-echo" {
			sreq.ServiceDefinitionRequirement = "proxy" //XXX
		} else if *command == "dm-echo" {
			sreq.ServiceDefinitionRequirement = "proxy"
		}
		sreq.InterfaceRequirements = []string{"HTTP-INSECURE-JSON"}
		var minVerReq int
		sreq.MinVersionRequirement = &minVerReq
		*sreq.MinVersionRequirement = 1

		reqJSON, _ := json.MarshalIndent(sreq, "", "  ")
		if *verbose != false {
			fmt.Println(string(reqJSON))
		}

		req, err := http.NewRequest("POST", *targetUri+"/query", bytes.NewBuffer(reqJSON))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Printf("Response statusCode: '%d'\n", resp.StatusCode)
		if resp.StatusCode == 200 {
			//fmt.Println("response Headers:", resp.Header)
			body, _ := ioutil.ReadAll(resp.Body)
			if *verbose != false {
				fmt.Println("response Body:", string(body))
			}

			serviceQueryResponse, err := ReadData2Object[ServiceQueryResponse](body)
			if err != nil {
				panic(err)
			}
			//fmt.Printf("%+v\n", serviceQueryResponse.ServiceQueryData[0])

			target := "http://" + serviceQueryResponse.ServiceQueryData[0].Provider.Address + ":" + strconv.Itoa(serviceQueryResponse.ServiceQueryData[0].Provider.Port) + serviceQueryResponse.ServiceQueryData[0].ServiceUri
			target = strings.Replace(target, "/proxy", "/echo", 1)

			fmt.Printf("Calling %s\n", target)
			data, err := getData(client, target)
			if err == nil {
				fmt.Println(string(data))
			}
		}
	}

	return 0
}

func ReadData2Object[T any](bytes []byte) (*T, error) {
	out := new(T)
	if err := json.Unmarshal(bytes, out); err != nil {
		return nil, err
	}
	return out, nil
}

func getData(client http.Client, uri string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("final error: %v\n", err)
		return nil, err
	}

	return body, nil
}

func fileExists(fileName string) (bool, error) {
	_, err := os.Stat(fileName)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
