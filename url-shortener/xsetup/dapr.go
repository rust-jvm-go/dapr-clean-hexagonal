package xsetup

import daprSdkClient "github.com/dapr/go-sdk/client"

type Dapr struct {
	daprClient daprSdkClient.Client
}

func NewDaprClient() (*Dapr, error) {
	daprClient, err := daprSdkClient.NewClient()
	if err != nil {
		return &Dapr{}, err
	}

	return &Dapr{daprClient: daprClient}, nil
}
