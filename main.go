package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"

	_authUsecase "belajar-bwa/auth/usecase"

	_userHandler "belajar-bwa/user/delivery/http"
	_userRepo "belajar-bwa/user/repository/mysql"
	_userUcase "belajar-bwa/user/usecase"

	_campaignHandler "belajar-bwa/campaign/delivery/http"
	_campaignRepo "belajar-bwa/campaign/repository/mysql"
	_campaignUsecase "belajar-bwa/campaign/usecase"

	_midtransUsecase "belajar-bwa/payment/usecase"

	_transactionHandler "belajar-bwa/transaction/delivery/http"
	_transactionRepo "belajar-bwa/transaction/repository/mysql"
	_transactionUsecase "belajar-bwa/transaction/usecase"
)

func init() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		gin.SetMode(gin.DebugMode)
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)

	dbConn, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		sqlDB, err := dbConn.DB()
		if err != nil {
			log.Panic(err)
		}
		sqlDB.Close()
	}()

	router := gin.Default()
	router.Use(cors.Default())
	router.Static("/images", "./images")
	api := router.Group("/api/v1")

	// Users
	authUsecase := _authUsecase.NewAuthUsecase()
	userRepository := _userRepo.NewMysqlUserRepository(dbConn)
	userUsecase := _userUcase.NewUserCase(userRepository)
	_userHandler.NewUserHandler(api, userUsecase, authUsecase)

	// Campaign
	campaignRepository := _campaignRepo.NewCampaignRepository(dbConn)
	campaignService := _campaignUsecase.NewCampaignUsecase(campaignRepository)
	_campaignHandler.NewCampaignHandler(api, campaignService, userUsecase, authUsecase)

	paymentService := _midtransUsecase.NewMidtransUsecase()

	transactionRepository := _transactionRepo.NewTransactionRepository(dbConn)
	transactionService := _transactionUsecase.NewTransactionUsecase(transactionRepository, campaignRepository, paymentService)
	_transactionHandler.NewTransactionHandler(api, transactionService, userUsecase, authUsecase)

	router.Run(viper.GetString("server.address"))
}
