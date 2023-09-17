package main

import (
	"api/configs"
	"api/internal/app/alfie/utils"
)

func main() {
	configs.LoadEnv()

	db := configs.ConnectToDatabase()
	configs.RunDatabaseMigrations(db)

	grpcServer := configs.CreateGRPCServer(db)
	listener := configs.GetGRPCListener()

	if err := grpcServer.Serve(listener); err != nil {
		utils.InfoLogger.Fatalf("failed to serve: %v", err)
	}
}
