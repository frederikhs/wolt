package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"frederikhs/wolt/wolt"
	"io/fs"
	"os"
)

const JsonFilename = "orders.json"

func WriteOrders(order *[]wolt.FullOrder) error {
	b, err := json.MarshalIndent(order, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf(JsonFilename), b, 0644)

	return err
}

func JsonExists() bool {
	if _, err := os.Stat(JsonFilename); errors.Is(err, fs.ErrNotExist) {
		return false
	}

	return true
}

func GetOrders() (*[]wolt.FullOrder, error) {
	b, err := os.ReadFile(JsonFilename)
	if err != nil {
		return nil, err
	}

	var orders []wolt.FullOrder
	err = json.Unmarshal(b, &orders)
	if err != nil {
		return nil, err
	}

	return &orders, nil
}
