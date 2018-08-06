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


		//shows the history of an wallet
		v1.GET("/wallet/:addr/history", handlers.GetWalletHistory)

		//shows the history of an wallet
		v1.GET("/wallet/:addr/history/evcoin", handlers.GetWalletHistoryEVCoin)

		//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		//~~~~~~~~~~~~~~~~~~~~~~~~~~ CPO ~~~~~~~~~~~~~~~~~~~~~~~~~~
		//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

		//gets all the info about an cpo
		v1.GET("/cpo", handlers.CpoInfo)

		//creates in the database a new cpo
		v1.POST("/cpo", handlers.CpoCreate)

		//generate a new wallet for the cpo
		v1.POST("/cpo/wallet/generate", handlers.CpoGenerateWallet)

		//displays the mnemonic seed for the cpo
		v1.GET("/cpo/wallet/seed", handlers.CpoGetSeed)

		//gets the CPO history of transactions
		v1.GET("/cpo/history", handlers.CpoHistory)


		//get locations
		v1.GET("/cpo/locations", handlers.CpoGetLocations)

		//update locations
		v1.PUT("/cpo/location", handlers.CpoPutLocation)

		//add locations
		v1.POST("/cpo/location", handlers.CpoPostLocation)

		//delete a location
		v1.DELETE("/cpo/location/:locationid", handlers.CpoDeleteLocation)


		//get tariffs
		v1.GET("/cpo/tariffs", handlers.CpoGetTariffs)

		//update tariff
		v1.PUT("/cpo/tariff", handlers.CpoPutTariff)

		//add tariff
		v1.POST("/cpo/tariff", handlers.CpoPostTariff)

		//delete a tariff
		v1.DELETE("/cpo/tariff/:tariffid", handlers.CpoDeleteTariff)


		//add one evse
		v1.POST("/cpo/evse", handlers.CpoPostEvse)

		//get history of transactions from msp
		v1.GET("/cpo/transactions/from_msp", handlers.CpoTransactionFromMsp)

		//~~~~ PAYMENTS SECTION ~~~~

		//  my wallet page
		v1.GET("/cpo/payment/wallet", handlers.CpoPaymentWallet)


		v1.GET("/cpo/payment/cdr/:token", handlers.CpoPaymentCDR)

		//creates new reimbursement
		v1.POST("/cpo/payment/reimbursement/:msp_address", handlers.CpoCreateReimbursement)

		//lists all reimbursements
		v1.GET("/cpo/payment/reimbursements/:status", handlers.CpoGetAllReimbursements)

		//sets the reimbursement as completed
		v1.PUT("/cpo/payment/reimbursement/:reimbursement_id/complete", handlers.CpoSetReimbursementComplete)

		//generates PDF for reimbursement id
		v1.GET("/cpo/payment/download_invoice/:reimbursement_id", handlers.CpoReimbursementGenPdf)
	}

}
