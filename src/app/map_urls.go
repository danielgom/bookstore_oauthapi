package app

func mapUrls() {
	router.GET("/oauth/accessToken/:atId", atHandler.GetById)
	router.POST("/oauth/accessToken", atHandler.Create)
	router.GET("/health", atHandler.Health)
}
