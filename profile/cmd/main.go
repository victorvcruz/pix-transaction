package main

import (
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"profile/internal/account"
	"profile/internal/cfg"
	"profile/internal/event"
	"profile/internal/key"
	"profile/internal/transaction"
	"profile/internal/user"
	"profile/platform/kafka"
	"profile/platform/sqlserver"
	proto "profile/proto/v1"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config, err := cfg.Load()
	if err != nil {
		return
	}

	list, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalf("Failed to listen port 9080 %v", err)
	}

	db, err := sqlserver.Start(config)
	if err != nil {
		log.Fatalf("Failed to connect to database %v", err)
	}

	err = db.AutoMigrate(&user.User{}, &key.Key{}, &account.Account{})
	if err != nil {
		log.Fatalf("Failed to migrate tables %v", err)
	}

	// repositories
	userRepository := user.NewRepository(db, config)
	keyRepository := key.NewRepository(db, config)
	accountRepository := account.NewRepository(db, config)

	// services
	userService := user.NewService(userRepository)
	keyService := key.NewService(keyRepository)
	accountService := account.NewService(accountRepository)

	//kafka
	kafkaConn := kafka.NewClient(config).Connect()
	events := event.NewEvent(kafkaConn, "transaction_events_topic",
		event.WithAttempts(4), event.WithBroker("localhost:9092"))

	transactionService := transaction.NewService(events)

	//server
	profileServer := NewProfileService(userService, accountService, keyService, transactionService)
	server := grpc.NewServer()
	proto.RegisterUserServiceServer(server, profileServer)
	proto.RegisterAccountServiceServer(server, profileServer)
	proto.RegisterKeysServiceServer(server, profileServer)

	log.Printf("Serve is running  on port: %v", "9080")
	if err := server.Serve(list); err != nil {
		log.Fatalf("Failed to serve gRPC server on port 9080: %v", err)
	}
}
