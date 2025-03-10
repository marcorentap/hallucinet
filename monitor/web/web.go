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

type ResponsePayload struct {
	HallucinetNetwork string
	Networks          map[string][]Entry
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
	listContainers := func(res http.ResponseWriter, req *http.Request) {
		var payload ResponsePayload
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Header().Set("Access-Control-Allow-Methods", "GET")
		res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if req.Method == "OPTIONS" {
			res.WriteHeader(http.StatusNoContent)
			return
		}

		networks := make(map[string][]Entry)
		networkNames := getNetworks(hctx)
		for _, network := range networkNames {
			networks[network] = getNetworkEntries(hctx, network)
		}

		payload.Networks = networks
		payload.HallucinetNetwork = hctx.Config.NetworkName
		jsonPayload, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			log.Panicf("Cannot marshal data %v: %v\n", payload, jsonPayload)
		}

		fmt.Fprintln(res, string(jsonPayload))
	}

	addr := fmt.Sprintf("%v:%v", hctx.Config.Host, hctx.Config.Port)
	http.HandleFunc("/containers", listContainers)
	http.ListenAndServe(addr, nil)
}
