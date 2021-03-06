package main

import (
	"time"

	supplyservice "github.com/diadata-org/diadata/internal/pkg/supplyService"
	"github.com/diadata-org/diadata/pkg/dia/helpers/ethhelper"
	models "github.com/diadata-org/diadata/pkg/model"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

func main() {

	ds, err := models.NewDataStore()
	if err != nil {
		log.Fatal("datastore error: ", err)
	}
	conn, err := ethclient.Dial("http://159.69.120.42:8545/")
	if err != nil {
		log.Fatal(err)
	}
	// Fetch token contract addresses from json file
	adds, err := ethhelper.GetAddressesFromFile()
	if err != nil {
		log.Fatal(err)
	}
	// Continuously update supplies once every 24h
	for {
		timeInit := time.Now()
		for _, address := range adds {
			supp, err := supplyservice.GetTotalSupplyfromMainNet(address, conn)
			if err != nil || len(supp.Symbol) < 2 || supp.Supply < 2 {
				continue
			}
			// Hotfix for some supplies:
			if supp.Symbol == "YAM" {
				supp.CirculatingSupply = float64(13907678)
			}
			if supp.Symbol == "CRO" {
				supp.CirculatingSupply = float64(20631963470)
			}

			ds.SetSupply(&supp)
			log.Info("set supply: ", supp)
		}
		timeFinal := time.Now()
		timeElapsed := timeFinal.Sub(timeInit)
		time.Sleep(24*60*60*time.Second - timeElapsed)
	}

}
