package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/marcorentap/hallucinet/types"
)

func getNetworks(hctx types.HallucinetContext) []string {
	db := hctx.DB
	query := `SELECT DISTINCT network_name FROM hallucinet`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Cannot query distinct networks: %v\n", err)
	}

	var networks []string
	for rows.Next() {
		var networkName string
		if err := rows.Scan(&networkName); err != nil {
			log.Fatalf("Cannot scan network name: %v\n", err)
		}
		networks = append(networks, networkName)
	}

	return networks
}

type Entry struct {
	ContainerID   string
	ContainerName string
	ContainerIP   string
}

func getNetworkEntries(hctx types.HallucinetContext, networkName string) []Entry {
	db := hctx.DB
	query := `SELECT container_id, container_name, container_ip FROM hallucinet where network_name = ?`
	rows, err := db.Query(query, networkName)
	if err != nil {
		log.Fatalf("Cannot query distinct networks: %v\n", err)
	}

	var entries []Entry
	for rows.Next() {
		var containerID string
		var containerName string
		var containerIP string
		if err := rows.Scan(&containerID, &containerName, &containerIP); err != nil {
			log.Fatalf("Cannot scan network name: %v\n", err)
		}
		entries = append(entries, Entry{
			ContainerID:   containerID,
			ContainerName: containerName,
			ContainerIP:   containerIP,
		})
	}

	return entries
}

func Serve(hctx types.HallucinetContext) {
	handler := func(res http.ResponseWriter, req *http.Request) {
		ret := make(map[string][]Entry)
		networks := getNetworks(hctx)
		for _, network := range networks {
			ret[network] = getNetworkEntries(hctx, network)
		}

		jsonResponse, err := json.MarshalIndent(ret, "", "  ")
		if err != nil {
			log.Panicf("Cannot marshal data %v: %v\n", ret, jsonResponse)
		}

		fmt.Fprintln(res, string(jsonResponse))
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
