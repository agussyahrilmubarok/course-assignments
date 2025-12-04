package controller

import (
	"os"
	"path/filepath"

	"example.com.backend/internal/config"
	"example.com.backend/internal/repos"
	"example.com.backend/internal/service"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(ginEngine *gin.Engine, cfg *config.Config, db *gorm.DB) {
	r := ginEngine

	userRepo := repos.NewUserRepository(db)
	campaignRepo := repos.NewCampaignRepository(db)
	campaignImageRepo := repos.NewCampaignImageRepository(db)
	transactionRepo := repos.NewTransactionRepository(db)

	uploadService := service.NewUploadService()
	userService := service.NewUserService(userRepo)
	campaignService := service.NewCampaignService(campaignRepo, campaignImageRepo)
	transactionService := service.NewTransactionService(transactionRepo)

	NewBaseController()
	homeController := NewHomeController()
	loginController := NewLoginController(userService)
	dashboardController := NewDashboardController()
	userController := NewUserController(userService, uploadService)
	campaignController := NewCampaignController(campaignService, userService, uploadService)
	transactionController := NewTransactionController(transactionService, userService, campaignService)

	r.Use(sessions.Sessions(cfg.Backend.Name, cookie.NewStore([]byte(cfg.Backend.Name))))

	r.HTMLRender = loadTemplate("./public/templates")
	r.Static("/assets", "./public/assets")
	r.Static("/uploads", "./public/uploads")
	r.StaticFile("/favicon.ico", "./public/assets/favicon.ico")

	r.GET("/", homeController.Index)
	r.GET("/login", loginController.Index)
	r.POST("/login", loginController.Login)

	adminDashboard := r.Group("/dashboard")
	{
		adminDashboard.GET("/", dashboardController.Index)

		users := adminDashboard.Group("/users")
		{
			users.GET("/", userController.Index)
			users.GET("/add", userController.Add)
			users.POST("/add", userController.Create)
			users.GET("/:id/edit", userController.Edit)
			users.POST("/:id/edit", userController.Update)
			users.GET("/:id/avatar", userController.Avatar)
			users.POST("/:id/avatar", userController.Upload)
			users.GET("/:id/delete", userController.Delete)
		}

		campaigns := adminDashboard.Group("/campaigns")
		{
			campaigns.GET("/", campaignController.Index)
			campaigns.GET("/add", campaignController.Add)
			campaigns.POST("/add", campaignController.Create)
			campaigns.GET("/:id/show", campaignController.Show)
			campaigns.GET("/:id/edit", campaignController.Edit)
			campaigns.POST("/:id/edit", campaignController.Update)
			campaigns.GET("/:id/image", campaignController.Image)
			campaigns.POST("/:id/image", campaignController.Upload)
			campaigns.GET("/:id/delete", campaignController.Delete)
		}

		transactions := adminDashboard.Group("/transactions")
		{
			transactions.GET("/", transactionController.Index)
			transactions.GET("/add", transactionController.Add)
			transactions.POST("/add", transactionController.Create)
			transactions.GET("/:id/show", transactionController.Show)
			transactions.GET("/:id/edit", transactionController.Edit)
			transactions.POST("/:id/edit", transactionController.Update)
			transactions.GET("/:id/delete", transactionController.Delete)
		}

		adminDashboard.GET("/logout", dashboardController.Logout)
	}
}

func loadTemplate(templateDir string) multitemplate.Renderer {
	renderer := multitemplate.NewRenderer()
	commons, err := filepath.Glob(templateDir + "/common/*.html")
	if err != nil {
		panic(err.Error())
	}

	homePages, err := filepath.Glob(templateDir + "/home/*.html")
	if err != nil {
		panic(err.Error())
	}

	for _, page := range homePages {
		if fileInfo, err := os.Stat(page); err == nil && !fileInfo.IsDir() {
			files := append([]string{filepath.Join(templateDir, "layouts", "home_layout.html")}, page)
			files = append(files, commons...)

			templateName := filepath.Base(page)
			// fmt.Println("Adding home template:", templateName)
			renderer.AddFromFiles(templateName, files...)
		}
	}

	authPages, err := filepath.Glob(templateDir + "/auth/*.html")
	if err != nil {
		panic(err.Error())
	}

	for _, page := range authPages {
		if fileInfo, err := os.Stat(page); err == nil && !fileInfo.IsDir() {
			files := append([]string{filepath.Join(templateDir, "layouts", "auth_layout.html")}, page)
			files = append(files, commons...)

			templateName := filepath.Base(page)
			// fmt.Println("Adding auth template:", templateName)
			renderer.AddFromFiles(templateName, files...)
		}
	}

	dashboardPages, err := filepath.Glob(templateDir + "/dashboard/**/*.html")
	if err != nil {
		panic(err.Error())
	}

	for _, page := range dashboardPages {
		if fileInfo, err := os.Stat(page); err == nil && !fileInfo.IsDir() {
			// Prepare the template files for dashboard layout
			files := append([]string{filepath.Join(templateDir, "layouts", "dashboard_layout.html")}, page)
			files = append(files, commons...)

			templateName := filepath.Base(page) // Gets just the file name
			// fmt.Println("Adding dashboard template:", templateName)
			renderer.AddFromFiles(templateName, files...)
		}
	}

	return renderer
}
