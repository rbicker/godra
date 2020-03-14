package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/rbicker/godra/internal/hydraclient"

	"github.com/rbicker/godra/internal/db"
	"github.com/rbicker/godra/internal/godra"
	"github.com/rbicker/godra/internal/utils"
)

func main() {
	var dbOpts []func(*db.MGO) error
	dbOpts = append(dbOpts, db.SetURL(utils.LoadSetting("MONGO_URL", "mongodb://localhost:27017")))
	dbOpts = append(dbOpts, db.SetDBName(utils.LoadSetting("MONGO_DB", "db")))
	dbOpts = append(dbOpts, db.SetCollectionName(utils.LoadSetting("MONGO_COLLECTION", "users")))
	con, err := db.NewMongoConnection(dbOpts...)
	if err != nil {
		log.Fatalf("error while creating mongodb connection: %v\n", err)
	}
	err = con.Connect()
	if err != nil {
		log.Fatalf("could not connect to mongodb: %v\n", err)
	}
	var srvOpts []func(*godra.Server) error
	var client hydraclient.Client
	client.SetHydraPrivateURL(utils.LoadSetting("HYDRA_PRIVATE_URL", "http://localhost:4445"))
	srvOpts = append(srvOpts, godra.SetHydraClient(client))
	port := utils.LoadSetting("PORT", "5000")
	p, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("invalid port '%s' given, unable to convert to integer", port)
	}
	srvOpts = append(srvOpts, godra.SetPort(p))
	log.Printf("connected to mongodb")
	srvOpts = append(srvOpts, godra.SetDatabase(con))
	srv, err := godra.NewServer(srvOpts...)
	if err != nil {
		log.Fatalf("error while creating new godra server: %v", err)
	}
	go func() {
		log.Printf("starting godra server on port %s\n", port)
		if err := srv.Serve(); err != nil {
			log.Printf("http server startup encountered an error: %v\n", err)
			con.Disconnect()
			os.Exit(1)
		}
	}()
	// wait for control+c to exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// block until a signal is received
	<-c
	con.Disconnect()
	//srv.Shutdown()
}
