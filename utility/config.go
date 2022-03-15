package utility

import (
	"encoding/json"
	"log"
	"os"
)

type configTemplate struct {
	RuleID         string `json:"rule_id"`
	Level          string `json:"level"`
	SoftHardID     string `json:"softhard"`
	Product        string `json:"product"`
	Company        string `json:"company"`
	Category       string `json:"category"`
	ParentCategory string `json:"parent_category"`
	Rules          [][]struct {
		Match   string `json:"match"`
		Content string `json:"content"`
	} `json:"rules"`
}

func OpenConfigFile(name string) {

	file, err := os.ReadFile(name)
	//config, err := os.ReadFile("fofa.json")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Print(string(config))
	err = json.Unmarshal(file, &config.FingerPrints)
	if err != nil {
		log.Fatal(err)
	}
	return
}
