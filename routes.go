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

		// Stations // EVSEs // Connectors
		v1.GET("/msp/locations", handlers.MspGetLocations)

		//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		//~~~~~~~~~~~~~~~~~~~~~~~~~~ CPO ~~~~~~~~~~~~~~~~~~~~~~~~~~
		//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

		//gets all the info about an cpo
		v1.GET("/cpo", handlers.CpoInfo)

		//creates in the database a new cpo
		v1.POST("/cpo", handlers.CpoCreate)

		//generate a new wallet for the cpo
		v1.POST("/cpo/wallet/generate", handlers.CpoGenerateWallet)

		//displays the mnemonic seed for the msp
		v1.GET("/cpo/wallet/seed", handlers.CpoGetSeed)

		//gets the MSP history of transactions
		v1.GET("/cpo/history", handlers.CpoHistory)

		//upload new locations
		v1.GET("/cpo/locations", handlers.CpoGetLocations)

		//upload new locations
		v1.PUT("/cpo/locations", handlers.CpoPutLocations)

		//add one location
		v1.POST("/cpo/location", handlers.CpoPostLocation)

		//delete a location
		v1.DELETE("/cpo/location/:locationid", handlers.CpoDeleteLocation)

		//add one evse
		v1.POST("/cpo/evse", handlers.CpoPostEvse)
	}

}
