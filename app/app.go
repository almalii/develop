package app

import (
	"context"
	"net/http"
	"notes-rew/internal/config"
	"notes-rew/internal/db/postgres"
	usersController "notes-rew/user/controller/handler"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"

	notesController "notes-rew/note/controller/handler"
	notesService "notes-rew/note/service"
	notesStorage "notes-rew/note/storage"
	notesUsecase "notes-rew/note/usecase"
	usersService "notes-rew/user/service"
	usersStorage "notes-rew/user/storage"
	usersUsecase "notes-rew/user/usecase"
)

type App struct {
	router chi.Router
}

func NewApp() *App {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	connectDB, err := postgres.NewConnectionDB(context.Background(), config.NewDbConfig())
	if err != nil {
		logrus.Print(err)
	}

	noteStorage := notesStorage.NewNoteStorage(connectDB)
	noteService := notesService.NewNoteService(noteStorage)
	noteUsecase := notesUsecase.NewNoteUsecase(noteService)
	noteController := notesController.NewNoteController(noteUsecase)
	noteController.Register(router)

	userStorage := usersStorage.NewPSQLUserStorage(connectDB)
	userService := usersService.NewUserService(userStorage)
	userUsecase := usersUsecase.NewUserUsecase(userService)
	userController := usersController.NewUserController(userUsecase)
	userController.Register(router)

	return &App{router: router}
}

func (a *App) Start(port string) error {
	httpServer := &http.Server{
		Addr:           ":" + port,
		Handler:        a.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if port == "" {
		logrus.Fatal("port is empty")
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			logrus.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	logrus.Println("server started on port " + port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return httpServer.Shutdown(ctx)
}
