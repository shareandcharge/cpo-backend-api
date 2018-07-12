package main

import "github.com/motionwerkGmbH/cpo-backend-api/handlers"

func InitializeRoutes() {

	v1 := router.Group("/api/v1")
	{
		// General
		v1.GET("/", handlers.Index)

		//used only to delete / reinit the database with default values.
		v1.DELETE("/s3cr3tReinitf32fdsfsdf98yu32jlkjfsd89yaf98j320j", handlers.Reinit)


		//CPO Management
		v1.POST("/cpo/create", handlers.CpoCreate)

		//CPO info
		v1.GET("/cpo/info", handlers.CpoInfo)
		v1.GET("/token/info", handlers.TokenInfo)




		// Account Related
		v1.GET("/account/info", handlers.AccountInfo)
		v1.GET("/account/wallet", handlers.WalletInfo)
		v1.GET("/account/history", handlers.AccountHistory)
		v1.GET("/account/mnemonic", handlers.AccountMnemonic)

		// Stations // EVSEs // Connectors
		v1.GET("/stations", handlers.StationsInfo)

	}

}
