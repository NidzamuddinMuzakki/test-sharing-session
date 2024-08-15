package main

import (
	"context"
	"time"
	_ "time/tzdata"

	"github.com/NidzamuddinMuzakki/test-sharing-vision/common/util"
	"github.com/NidzamuddinMuzakki/test-sharing-vision/repository"

	"github.com/NidzamuddinMuzakki/test-sharing-vision/config"

	"github.com/NidzamuddinMuzakki/test-sharing-vision/service"
	// Import services here
	serviceHealth "github.com/NidzamuddinMuzakki/test-sharing-vision/service/health"

	// Import deliveries here

	httpDelivery "github.com/NidzamuddinMuzakki/test-sharing-vision/handler"
	httpDeliveryHealth "github.com/NidzamuddinMuzakki/test-sharing-vision/handler/health"

	// Import cmd here
	cmdHttp "github.com/NidzamuddinMuzakki/test-sharing-vision/cmd/http"

	// Import common lib here

	commonDs "github.com/NidzamuddinMuzakki/test-sharing-vision/go-lib-common/data_source"
	"github.com/NidzamuddinMuzakki/test-sharing-vision/go-lib-common/logger"

	commonPanicRecover "github.com/NidzamuddinMuzakki/test-sharing-vision/go-lib-common/middleware/gin/panic_recovery"

	commonRegistry "github.com/NidzamuddinMuzakki/test-sharing-vision/go-lib-common/registry"

	commonTime "github.com/NidzamuddinMuzakki/test-sharing-vision/go-lib-common/time"
	commonValidator "github.com/NidzamuddinMuzakki/test-sharing-vision/go-lib-common/validator"

	// Import third parties here
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	_ "github.com/spf13/viper/remote"
)

func main() {
	ctx := context.Background()

	// Start Init //
	loc, err := time.LoadLocation(commonTime.LoadTimeZoneFromEnv())
	if err != nil {
		panic(err)
	}
	time.Local = loc
	// Configuration
	config.Init()
	// Logger
	logger.Init(logger.Config{
		AppName: config.Cold.AppName,
		Debug:   config.Hot.AppDebug,
	})
	// Validator
	validator := commonValidator.New()
	// Sentry

	// Database
	// - Master
	master, err := commonDs.NewDB(&commonDs.Config{
		Driver:                config.Cold.DBMysqlMasterDriver,
		Host:                  config.Cold.DBMysqlMasterHost,
		Port:                  config.Cold.DBMysqlMasterPort,
		DBName:                config.Cold.DBMysqlMasterDBName,
		User:                  config.Cold.DBMysqlMasterUser,
		Password:              config.Cold.DBMysqlMasterPassword,
		SSLMode:               config.Cold.DBMysqlMasterSSLMode,
		MaxOpenConnections:    config.Cold.DBMysqlMasterMaxOpenConnections,
		MaxLifeTimeConnection: config.Cold.DBMysqlMasterMaxLifeTimeConnection,
		MaxIdleConnections:    config.Cold.DBMysqlMasterMaxIdleConnections,
		MaxIdleTimeConnection: config.Cold.DBMysqlMasterMaxIdleTimeConnection,
	})
	if err != nil {
		panic(err)
	}
	// - Slave
	slave, err := commonDs.NewDB(&commonDs.Config{
		Driver:                config.Cold.DBMysqlSlaveDriver,
		Host:                  config.Cold.DBMysqlSlaveHost,
		Port:                  config.Cold.DBMysqlSlavePort,
		DBName:                config.Cold.DBMysqlSlaveDBName,
		User:                  config.Cold.DBMysqlSlaveUser,
		Password:              config.Cold.DBMysqlSlavePassword,
		SSLMode:               config.Cold.DBMysqlSlaveSSLMode,
		MaxOpenConnections:    config.Cold.DBMysqlSlaveMaxOpenConnections,
		MaxLifeTimeConnection: config.Cold.DBMysqlSlaveMaxLifeTimeConnection,
		MaxIdleConnections:    config.Cold.DBMysqlSlaveMaxIdleConnections,
		MaxIdleTimeConnection: config.Cold.DBMysqlSlaveMaxIdleTimeConnection,
	})
	if err != nil {
		panic(err)
	}
	// Activity Log Client

	// Panic Recovery
	panicRecoveryMiddleware := commonPanicRecover.NewPanicRecovery(
		validator,
		commonPanicRecover.WithConfigEnv(config.Cold.AppEnv),
	)

	// Registry
	common := commonRegistry.NewRegistry(

		commonRegistry.WithValidator(validator),

		commonRegistry.WithPanicRecoveryMiddleware(panicRecoveryMiddleware),
	)
	// End Init //

	// Start Clients //
	// ...
	// End Clients //

	// Start Repositories //
	masterUtilTx := util.NewTransactionRunner(master)
	PostsRepository := repository.NewPostsRepository(common, master, slave)

	repoRegistry := repository.NewRegistryRepository(masterUtilTx, PostsRepository)
	// End Repositories //

	// Start Services //
	postService := service.NewPostsService(common, repoRegistry)
	healthService := serviceHealth.NewHealth(master, slave)
	serviceRegistry := service.NewRegistry(
		healthService,
		postService,
	)
	// End Deliveries //

	// Start Deliveries //
	healthDelivery := httpDeliveryHealth.NewHealth(common, healthService)
	postDelivery := httpDelivery.NewPosts(common, serviceRegistry)
	registryDelivery := httpDelivery.NewRegistry(healthDelivery, postDelivery)
	// End Deliveries //

	//

	// Start HTTP Server //
	httpServer := cmdHttp.NewServer(
		common,
		registryDelivery,
	)
	httpServer.Serve(ctx)
	// End HTTP Server //
}
