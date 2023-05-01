package xsetup

import daprSdkClient "github.com/dapr/go-sdk/client"

type DaprClient struct {
	daprClient daprSdkClient.Client
}

func NewDaprClient() (DaprClient, error) {
	daprClient, err := daprSdkClient.NewClient()
	if err != nil {
		return DaprClient{}, err
	}

	return DaprClient{daprClient: daprClient}, nil
}
