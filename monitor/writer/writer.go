package writer

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/marcorentap/hallucinet/types"
	_ "github.com/mattn/go-sqlite3"
)

func AddEntry(db *sql.DB, e types.HallucinetEvent) {
	const query string = `
		INSERT INTO hallucinet (container_ip, container_id, container_name, network_id, network_name)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(container_id, network_id) 
		DO UPDATE SET 
			container_ip = excluded.container_ip,
		    container_name = excluded.container_name,
		    network_name = excluded.network_name;
	`
	_, err := db.Exec(query, e.ContainerIP, e.ContainerID, e.ContainerName, e.NetworkID, e.NetworkName)
	if err != nil {
		log.Panicf("Cannot insert entry: %v\n", err)
	}
}

func RemoveEntry(db *sql.DB, e types.HallucinetEvent) {
	const query string = `
		DELETE FROM hallucinet WHERE container_id = ? AND network_id = ?
	`
	_, err := db.Exec(query, e.ContainerID, e.NetworkID)
	if err != nil {
		log.Panicf("Cannot insert entry: %v\n", err)
	}
}

func InitializeDB(db *sql.DB) {
	const query string = `
		DROP TABLE IF EXISTS hallucinet;
		CREATE TABLE hallucinet (
		container_ip TEXT,
	    container_id TEXT,
	    container_name TEXT,
	    network_id TEXT,
	    network_name TEXT,
	    PRIMARY KEY (container_id, network_id)
		);
		`
	_, err := db.Exec(query)
	if err != nil {
		log.Panicf("Cannot insert entry: %v\n", err)
	}
}

func UpdateHosts(hctx types.HallucinetContext) {

	db := hctx.DB
	const query string = `
		SELECT container_ip, container_name FROM hallucinet WHERE network_name = ?
	`
	for {
		rows, err := db.Query(query, hctx.Config.NetworkName)
		if err != nil {
			log.Printf("Cannot query container IP and name: %v\n", err)
		}

		file, err := os.Create(hctx.Config.HostsPath)
		if err != nil {
			log.Panicf("Cannot open hosts file %v for writing: %v\n", hctx.Config.HostsPath, err)
		}
		defer file.Close()

		for rows.Next() {
			var containerIP string
			var containerName string
			if err := rows.Scan(&containerIP, &containerName); err != nil {
				log.Panicf("Cannot read container name in network %v: %v\n",
					hctx.Config.NetworkName, containerName)
			}
			fmt.Fprintf(file, "%v %v%v\n", containerIP, containerName, hctx.Config.DomainSuffix)
		}

		time.Sleep(5 * time.Second)
	}
}

func UpdateDatabase(hctx types.HallucinetContext) {
	db := hctx.DB
	eventChan := hctx.EventChan
	for event := range eventChan {
		switch event.Kind {
		case types.ContainerConnected:
			AddEntry(db, event)
		case types.ContainerDisconnected:
			RemoveEntry(db, event)
		}
	}
}
