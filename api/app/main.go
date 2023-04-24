package main

import (
	"api/app/config"
	"api/app/utils"
)

func main() {
	config.LoadEnv()

	db := config.ConnectToDatabase()
	config.RunDatabaseMigrations(db)

	grpcServer := config.CreateGRPCServer(db)
	listener := config.GetGRPCListener()

	if err := grpcServer.Serve(listener); err != nil {
		utils.InfoLogger.Fatalf("failed to serve: %v", err)
	}
}
