package main

import (
	"fmt"
	cf "github.com/cloudfoundry-community/go-cfclient"
	"time"
	"os"
	"net/http"
	"github.com/gorilla/mux"
	"bytes"
	"encoding/json"
)

func main() {
	client, err := createCfClient()
	if err != nil {
		fmt.Printf("Failed to create client: %s", err.Error())
		return
	}

	for i:=0; i<=5; i++ {
		time.Sleep(1000 * time.Millisecond)
		serviceName := fmt.Sprintf("cf-redis-benchmark-%d", i)
		go createDedicatedServiceInstance(client, serviceName)
	}

	bindToServiceInstance(client)


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

func createDedicatedServiceInstance(c *cf.Client, serviceName string) {
	//test that instance creations don't wait on each other
	time.Sleep(1 * time.Second)

	fmt.Println(fmt.Sprintf("creating instance %s", serviceName))

	//

	req := cf.ServiceInstanceRequest{
		Name:            serviceName,
		SpaceGuid:       os.Getenv("SPACE_GUID"),
		ServicePlanGuid: os.Getenv("SERVICE_PLAN_GUID"),
	}
	startTime := time.Now()
	serviceInstance, err := c.CreateServiceInstance(req)
	provisionDuration := time.Since(startTime)

	commonInfoString := fmt.Sprintf("create service %s request at %s took %s and", serviceName,
		startTime,
		provisionDuration)


	if err != nil {
		fmt.Println(fmt.Sprintf("%s failed with error: %s", commonInfoString, err.Error()))
	} else {
		fmt.Println(fmt.Sprintf("%s succeeded creating instance with GUID: %s", commonInfoString, serviceInstance.Guid))
	}

	// now delete the service instance:
	//requestURL := "/v2/service_instances/1aaeb02d-16c3-4405-bc41-80e83d196dff?accepts_incomplete=true"
	//r := c.NewRequest("DELETE", requestURL)
	//response, err := c.DoRequest(r)

}

type ServiceKeyRequest struct {
	Name            string `json:"name"`
	ServiceInstanceGuid       string `json:"service_instance_guid"`
}

func bindToServiceInstance(c *cf.Client) error {

	serviceKeyRequest := ServiceKeyRequest{Name:"myKey", ServiceInstanceGuid: os.Getenv("INSTANCE_FOR_BINDING_GUID")}
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(serviceKeyRequest)
	if err != nil {
		fmt.Println(fmt.Sprintf("Request buffer encoding failed with error: %s", err.Error()))
		return err
	}

	r := c.NewRequestWithBody("POST", "/v2/service_keys", buf)

	startTime := time.Now()
	_, err = c.DoRequest(r)
	bindDuration := time.Since(startTime)

	fmt.Println(bindDuration)

	return nil
}

func getPort() string {
	if configuredPort := os.Getenv("PORT"); configuredPort == "" {
		return "8080"
	} else {
		return configuredPort
	}
}
