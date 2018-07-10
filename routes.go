package main

func InitializeRoutes() {

	v1 := router.Group("/api/v1")
	{
		// Handle the index route
		v1.GET("/", HandleIndex)

		//CPO Management
		v1.POST("/cpo/create", HandleCpoCreate)

		//CPO info
		v1.GET("/cpo/info", HandleCpoInfo)

		//wallet query
		v1.GET("/wallet/info", HandleWalletInfo)


		//used only to delete / reinit the database with default values.
		v1.DELETE("/s3cr3tReinitf32fdsfsdf98yu32jlkjfsd89yaf98j320j", HandleReinit)

	}

}
