package server

func BindRoutes() {
	MainServer.Get("/", IndexRoute)

	authRoute := MainServer.Group("/api/auth")
	authRoute.Post("/register", UserRegisterRoute)
	authRoute.Post("/login", UserLoginRoute)

	antiScamRoute := authRoute.Group("/anti-scam", UserPermissionMiddleware)
	antiScamRoute.Post("/upload-audio", AntiScamUploadAudioRoute)
	antiScamRoute.Post("/analyze", AntiScamAnalyzeRoute)

	userRoute := MainServer.Group("/user", UserPermissionMiddleware)
	userRoute.Post("/logout", UserLogoutRoute)
	userRoute.Put("/password/change", UserChangePasswordRoute)

	linkRoute := userRoute.Group("/link")
	linkRoute.Get("/list", LinkListRoute)
	linkRoute.Post("/add", LinkAddRoute)
	linkRoute.Delete("/remove/:id", LinkRemoveRoute)

	aiRoute := MainServer.Group("/ai", UserPermissionMiddleware)
	aiRoute.Post("/run", AIApiRoute)

	dataRoute := MainServer.Group("/data", UserPermissionMiddleware)
	dataRoute.Post("/add", DataAddRoute)
	dataRoute.Get("/get", DataGetRoute)
	dataRoute.Post("/cutget", DataCutGetRoute)
	dataRoute.Get("/count", DataCountRoute)
}
