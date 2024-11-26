package utils

import (
	"encoding/json"
	"fmt"
)

func DisplayList(list []string, asJson bool) {
	if asJson {
		data, _ := json.Marshal(&list)
		fmt.Println(string(data))
		return
	}
	for _, item := range list {
		fmt.Println(item)
	}
}

func DisplayMap(m map[string]string, asJson bool) {
	if asJson {
		data, _ := json.Marshal(&m)
		fmt.Println(string(data))
		return
	}
	for key, value := range m {
		fmt.Printf("%s: %s\n", key, value)
	}
}
