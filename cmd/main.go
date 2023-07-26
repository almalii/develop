package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/app/grpc_app"
	"notes-rew/internal/app/rest_app"
	"notes-rew/internal/config"
	"sync"
)

func main() {

	logrus.SetFormatter(&logrus.JSONFormatter{})
	//logrus.SetLevel(logrus.DebugLevel)

	cfg := config.InitConfig()
	ctx := context.Background()

	newApp := rest_app.NewApp(ctx, cfg)
	newAppGRPC := grpc_app.NewAppGRPC(ctx, cfg)

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := newAppGRPC.StartGRPC(); err != nil {
			logrus.Fatalf("Не удалось запустить GRPC-сервер: %+v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := newApp.Start(); err != nil {
			logrus.Fatalf("Не удалось запустить приложение: %+v", err)
		}
	}()

	wg.Wait()

}
