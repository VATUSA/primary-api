package utils

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func Pprint(obj interface{}) {
	marshaled, err := json.MarshalIndent(obj, "", "   ")
	if err != nil {
		log.Fatalf("marshaling error: %s", err)
	}
	fmt.Println(string(marshaled))
}
