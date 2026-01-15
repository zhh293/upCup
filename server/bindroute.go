package server

func BindRoutes() {
	MainServer.Get("/", IndexRoute)

	userRoute := MainServer.Group("/user")
	userRoute.Post("/register", UserRegisterRoute)
	userRoute.Post("/login", UserLoginRoute)

	aiRoute := MainServer.Group("/ai", UserPermissionMiddleware)
	aiRoute.Post("/run", AIApiRoute)

	dataRoute := MainServer.Group("/data", UserPermissionMiddleware)
	dataRoute.Post("/add", DataAddRoute)
	dataRoute.Get("/get", DataGetRoute)
	dataRoute.Post("/cutget", DataCutGetRoute)
	dataRoute.Get("/count", DataCountRoute)

	linkRoute := MainServer.Group("/link", UserPermissionMiddleware)
	linkRoute.Post("/add", LinkAddRoute)
	linkRoute.Post("/exist", LinkExsitRoute)
	linkRoute.Post("/get", LinkGetRoute)
}
