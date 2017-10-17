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

	createDedicatedServiceInstance(client)

	router := mux.NewRouter()
	if err := http.ListenAndServe(fmt.Sprintf(":%v", getPort()), router); err != nil {
		fmt.Errorf(err.Error())
	}
}
func createCfClient() (*cf.Client, error) {
	c := &cf.Config{
		ApiAddress: "https://api.sys.pikachu.gcp.london.cf-app.com",
		Username:   "admin",
		Password:   "replace!!!!!!!!!!!!!!!!!!",
		SkipSslValidation: true,
	}
	return cf.NewClient(c)
}

func createDedicatedServiceInstance(c *cf.Client) {
	fmt.Println("creating instance")

	req := cf.ServiceInstanceRequest{
		Name:            fmt.Sprintf("cf-redis-benchmark-%s", "testtrinasdf"),
		SpaceGuid:       os.Getenv("SPACE_GUID"),
		ServicePlanGuid: os.Getenv("SERVICE_PLAN_GUID"),
	}
	startTime := time.Now()
	serviceInstance, err := c.CreateServiceInstance(req)
	duration := time.Since(startTime)

	fmt.Println("request sent")
	fmt.Println(duration)

	instances, err := c.ListServiceInstances()

	if err != nil {
		fmt.Println("failed to list instances")
	}

	fmt.Println(instances)

	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to create instance with err: %s, in time: %f", err.Error(), duration))
	} else {
		fmt.Println(fmt.Sprintf("Succeeded creating instance: %s with guid: %s, in time: %s",
			serviceInstance.Name, serviceInstance.Guid, duration.String()))
	}

}

func getPort() string {
	if configuredPort := os.Getenv("PORT"); configuredPort == "" {
		return "8080"
	} else {
		return configuredPort
	}
}
