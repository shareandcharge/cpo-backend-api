package main

import "github.com/motionwerkGmbH/cpo-backend-api/handlers"

func InitializeRoutes() {

	v1 := router.Group("/api/v1")
	{
		//~~~~~~~~~~~~~~~~~~~~ GENERAL STUFF ~~~~~~~~~~~~~~~~~~~~~~

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

		//shows the history of an wallet
		v1.GET("/wallet/:addr/history/evcoin", handlers.GetWalletHistoryEVCoin)

		//view CDRs
		v1.GET("/view_cdrs/:reimbursement_id", handlers.ViewCDRs)

		//~~~~~~~~~~~~~~~~~~~~~~~~~~ CPO ~~~~~~~~~~~~~~~~~~~~~~~~~~

		//gets all the info about an cpo
		v1.GET("/cpo", handlers.CpoInfo)

		//creates in the database a new cpo
		v1.POST("/cpo", handlers.CpoCreate)

		//displays the mnemonic seed for the cpo
		v1.GET("/cpo/wallet/seed", handlers.CpoGetSeed)

		//get locations
		v1.GET("/cpo/locations", handlers.CpoGetLocations)

		//add locations
		v1.POST("/cpo/locations", handlers.CpoPostLocations)

		//update 1 location
		v1.PUT("/cpo/location/:scid", handlers.CpoPutLocation)

		//add 1 location
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
		v1.DELETE("/cpo/tariffs", handlers.CpoDeleteTariffs)

		//add one evse
		v1.POST("/cpo/evse", handlers.CpoPostEvse)

		//~~~~~~~~~~~~~~~~~~~~~~~~~~ CPO PAYMENT ~~~~~~~~~~~~~~~~~~~~~~~~~~

		//  my wallet page
		v1.GET("/cpo/payment/wallet", handlers.CpoPaymentWallet)

		// the records for the particular token
		v1.GET("/cpo/payment/cdr/:token", handlers.CpoPaymentCDR)

		//creates new reimbursement
		v1.POST("/cpo/payment/reimbursement/:msp_address", handlers.CpoCreateReimbursement)

		//lists all reimbursements filtered by
		v1.GET("/cpo/payment/reimbursements/:status", handlers.CpoGetAllReimbursements)

		//sets the reimbursement as completed (TODO: don't use it!)
		v1.PUT("/cpo/payment/reimbursement/:reimbursement_id/:status", handlers.CpoSetReimbursementStatus)

		//sets the reimbursement as completed
		v1.POST("/cpo/payment/send_tokens_to_msp/:reimbursement_id", handlers.CpoSendTokensToMsp)

	}

}
