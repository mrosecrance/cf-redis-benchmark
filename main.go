package main

import (
	"fmt"
	cf "github.com/cloudfoundry-community/go-cfclient"
	"time"
	"os"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	client, err := createCfClient()
	if err != nil {
		fmt.Printf("Failed to create client: %s", err.Error())
		return
	}

	for i:=0; i<=5; i++ {
		time.Sleep(1000 * time.Millisecond)
		createDedicatedServiceInstance(client, i)
	}


	router := mux.NewRouter()
	if err := http.ListenAndServe(fmt.Sprintf(":%v", getPort()), router); err != nil {
		fmt.Errorf(err.Error())
	}
}
func createCfClient() (*cf.Client, error) {
	c := &cf.Config{
		ApiAddress: "https://api.sys.pikachu.gcp.london.cf-app.com",
		Username:  os.Getenv("CF_USERNAME"),
		Password:   os.Getenv("CF_PASSWORD"),
		SkipSslValidation: true,
	}
	return cf.NewClient(c)
}

func createDedicatedServiceInstance(c *cf.Client, serviceIndex int) {
	fmt.Println("creating instance")

	serviceName := fmt.Sprintf("cf-redis-benchmark-%d", serviceIndex)

	req := cf.ServiceInstanceRequest{
		Name:            serviceName,
		SpaceGuid:       os.Getenv("SPACE_GUID"),
		ServicePlanGuid: os.Getenv("SERVICE_PLAN_GUID"),
	}
	startTime := time.Now()
	_, err := c.CreateServiceInstance(req)
	provisionDuration := time.Since(startTime)

	commonInfoString := fmt.Sprintf("create service %s request at %s took %s and", serviceName,
		startTime,
		provisionDuration)


	if err != nil {
		fmt.Println(fmt.Sprintf("%s failed with error: %s", commonInfoString, err.Error()))
	} else {
		fmt.Println(fmt.Sprintf("%s succeeded", commonInfoString))
	}



}

func getPort() string {
	if configuredPort := os.Getenv("PORT"); configuredPort == "" {
		return "8080"
	} else {
		return configuredPort
	}
}
