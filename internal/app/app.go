package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/randnull/banner-service/internal/controllers"
	"github.com/randnull/banner-service/internal/repository"
	"github.com/randnull/banner-service/internal/service"
	"github.com/randnull/banner-service/internal/config"
)


type App struct {
	banner_repo   	*repository.Repository
	banner_service 	*service.BannerService
	banner_handlers *controllers.BannerHandlers

	user_repo 		*repository.UserRepository
	user_service 	*service.UserService
	user_handlers	*controllers.UserHandlers

	cfg 			*config.Config
}


func NewApp(cfg *config.Config) *App {
	repo_banner := repository.NewRepository(cfg)
	server_banner := service.NewBannerSevice(repo_banner)
	controller_banner := controllers.NewHandler(server_banner)

	repo_user := repository.NewUserRepository(cfg)
	serv_user := service.NewUserSevice(repo_user)
	controllers_user := controllers.NewUserHandler(serv_user, cfg)

	app := &App{
		banner_repo:     repo_banner,
		banner_service:  server_banner,
		banner_handlers: controller_banner,

		user_repo: 		 repo_user,
		user_service: 	 serv_user,
		user_handlers: 	 controllers_user,

		cfg: 			 cfg,
	}

	return app
}

func (a *App) Run() {
	router := mux.NewRouter()
	
	router.HandleFunc("/register/user", a.user_handlers.RegisterUser).Methods("POST")
	router.HandleFunc("/register/admin", a.user_handlers.RegisterAdmin).Methods("POST")
	router.HandleFunc("/login", a.user_handlers.Login).Methods("POST")

	banner_router := router.PathPrefix("/").Subrouter()
	banner_router.Use(a.user_handlers.Auth)

	banner_router.HandleFunc("/user_banner", a.banner_handlers.GetBanner).Methods("GET")
	banner_router.HandleFunc("/banner", a.banner_handlers.GetAllBanners).Methods("GET")
	banner_router.HandleFunc("/banner", a.banner_handlers.CreateBanner).Methods("POST")
	banner_router.HandleFunc("/banner/{id}", a.banner_handlers.DeleteBanner).Methods("DELETE")
	banner_router.HandleFunc("/banner/{id}", a.banner_handlers.UpdateBanner).Methods("PATCH")

	addr := fmt.Sprintf(":%v", a.cfg.ServerPort)

	log.Printf("Listen on %s\n", addr)

	http.ListenAndServe(addr, router)
}
