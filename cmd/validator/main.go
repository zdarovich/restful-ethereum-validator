package main

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"restful-ethereum-validator/beaconclient"
	"restful-ethereum-validator/config"
	"restful-ethereum-validator/server"
	"restful-ethereum-validator/service"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
		return
	}

	ethClient, err := ethclient.Dial(cfg.RPCDialURL)
	if err != nil {
		logger.Fatalf("Failed to connect to Ethereum node: %v", err)
	}
	defer ethClient.Close()
	beaconClient := beaconclient.NewBeaconClient(cfg.RPCDialURL)
	ethService := service.NewEthereumService(beaconClient, ethClient)

	apiHandler := server.NewAPIHandler(logger, ethService)

	r := gin.Default()

	r.GET("/blockreward/:slot", apiHandler.GetBlockReward)
	r.GET("/syncduties/:slot", apiHandler.GetSyncDuties)

	if err := r.Run(":8080"); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
