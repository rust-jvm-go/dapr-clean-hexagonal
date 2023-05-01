package xsetup

import (
	"context"
	"fmt"
	"strconv"
)

var (
	daprConfigStore      = "configstore"
	bunRouterDebugConfig = "BUNROUTER_DEBUG"
	bunRouterPortConfig  = "BUNROUTER_PORT"
	mongoDBURIConfig     = "MONGODB_URI"
	mongoDatabaseConfig  = "MONGODB_DATABASE"
	mongoDBTimeoutConfig = "MONGODB_TIMEOUT"
)

type InitConfig struct {
	Ctx            context.Context
	BunRouterDebug bool
	BunRouterPort  string
	MongoDBURI     string
	MongoDatabase  string
	MongoDBTimeout int
}

func NewInitConfig(ctx context.Context, dapr DaprClient) (InitConfig, error) {
	// Get config items from config store.
	var configurationItems = []string{
		bunRouterDebugConfig, bunRouterPortConfig,
		mongoDBURIConfig, mongoDatabaseConfig, mongoDBTimeoutConfig,
	}
	configItems, err := dapr.daprClient.GetConfigurationItems(ctx, daprConfigStore, configurationItems)
	if err != nil {
		fmt.Println("Could not get config items from DaprClient, err: ", err.Error())
		return InitConfig{}, err
	}

	var bunRouterDebug bool
	var bunRouterPort, mongoDBURI, mongoDatabase string
	var mongoDBTimeout int

	for _, item := range configurationItems {
		configItem, ok := configItems[item]
		if !ok {
			fmt.Printf("Could not get config item %s, err: %v\n", item, err.Error())
			return InitConfig{}, err
		}
		fmt.Printf("Configuration item \"%s\" = %s\n", item, configItem.Value)

		switch item {
		case bunRouterDebugConfig:
			bunRouterDebug, err = strconv.ParseBool(configItem.Value)
			if err != nil {
				fmt.Printf("Could not parse config item %s, err: %v\n", bunRouterDebugConfig, err.Error())
				return InitConfig{}, err
			}
			break
		case bunRouterPortConfig:
			bunRouterPort = configItem.Value
			break
		case mongoDBURIConfig:
			mongoDBURI = configItem.Value
			break
		case mongoDatabaseConfig:
			mongoDatabase = configItem.Value
			break
		case mongoDBTimeoutConfig:
			mongoDBTimeout, err = strconv.Atoi(configItem.Value)
			if err != nil {
				mongoDBTimeout = 10
			}
			break
		}
	}

	return InitConfig{
		Ctx:            ctx,
		BunRouterDebug: bunRouterDebug,
		BunRouterPort:  bunRouterPort,
		MongoDBURI:     mongoDBURI,
		MongoDatabase:  mongoDatabase,
		MongoDBTimeout: mongoDBTimeout,
	}, nil
}
