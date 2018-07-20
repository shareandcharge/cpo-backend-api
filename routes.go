package main

import "github.com/motionwerkGmbH/cpo-backend-api/handlers"

func InitializeRoutes() {

	v1 := router.Group("/api/v1")
	{

		//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		//~~~~~~~~~~~~~~~~~~~~ GENERAL STUFF ~~~~~~~~~~~~~~~~~~~~~~
		//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

		v1.GET("/", handlers.Index)

		//used only to delete / reinit the database with default values.
		v1.DELETE("/s3cr3tReinitf32fdsfsdf98yu32jlkjfsd89yaf98j320j", handlers.Reinit)

		//shows the token info
		v1.GET("/token/info", handlers.TokenInfo)

		//shows the token balance
		v1.GET("/token/balance/:addr", handlers.TokenBalance)

		//Tops up the balance of the EV Driver
		v1.POST("/token/mint/:addr", handlers.TokenMint)

		//shows the balance in eth of a wallet
		v1.GET("/wallet/:addr", handlers.GetWalletBalance)

		// Stations // EVSEs // Connectors
		v1.GET("/locations", handlers.LocationsInfo)

		//returns a list of all EV Drivers with their details & balances
		v1.GET("/drivers", handlers.GetAllDrivers)


		//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		//~~~~~~~~~~~~~~~~~~~~~~~~~~ MSP ~~~~~~~~~~~~~~~~~~~~~~~~~~
		//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

		//gets all the info about an msp
		v1.GET("/msp", handlers.MspInfo)

		//creates in the database a new msp
		v1.POST("/msp", handlers.MspCreate)

		//generate a new wallet for the msp
		v1.POST("/msp/wallet/generate", handlers.MspGenerateWallet)

		//displays the mnemonic seed for the msp
		v1.GET("/msp/wallet/seed", handlers.MspGetSeed)

		//gets the MSP history of transactions
		v1.GET("/msp/history", handlers.MSPHistory)

		//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		//~~~~~~~~~~~~~~~~~~~~~~~~~~ CPO ~~~~~~~~~~~~~~~~~~~~~~~~~~
		//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

	}

}
