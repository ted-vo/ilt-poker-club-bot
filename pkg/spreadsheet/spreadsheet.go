package spreadsheet

import (
	"io/ioutil"

	"github.com/apex/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

const SPREADSHEET_ID = "1Jxc5UY9GUrv5Q7SUOdMSDPJ2g2VEZ-lRthq-WUiifKo"

type SpreadsheetClub struct {
	Service     *spreadsheet.Service
	Spreadsheet *spreadsheet.Spreadsheet
}

func GetSheet() *SpreadsheetClub {
	data, err := ioutil.ReadFile("./config/client_secret.json")
	if err != nil {
		log.Error(err.Error())
	}

	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	if err != nil {
		log.Error(err.Error())
	}

	client := conf.Client(oauth2.NoContext)

	service := spreadsheet.NewServiceWithClient(client)

	spreadsheet, err := service.FetchSpreadsheet(SPREADSHEET_ID)
	if err != nil {
		log.Error(err.Error())
	}

	log.Infof("get spreadsheet success. ID=%s", spreadsheet.ID)

	return &SpreadsheetClub{
		Service:     service,
		Spreadsheet: &spreadsheet,
	}
}
