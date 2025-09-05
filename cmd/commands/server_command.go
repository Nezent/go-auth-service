package commands

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Nezent/auth-service/internal/application/services"
	"github.com/Nezent/auth-service/internal/infrastructure/config"
	"github.com/Nezent/auth-service/internal/infrastructure/persistence"
	"github.com/Nezent/auth-service/internal/infrastructure/repository"
	"github.com/Nezent/auth-service/internal/interfaces/handlers"
	"github.com/Nezent/auth-service/internal/interfaces/routes"
	"github.com/Nezent/auth-service/pkg/router"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var serverCommand = &cobra.Command{
	Use:   "server",
	Short: "api server",
	Long:  "start the api server",
	RunE: func(cmd *cobra.Command, args []string) error {
		app := fx.New(
			router.Module,
			routes.Module,
			config.Module,
			persistence.Module,
			handlers.Module,
			services.Module,
			repository.Module,
			fx.Invoke(func(
				cfg *config.Config,
				router *chi.Mux,
				lc fx.Lifecycle,
			) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						log.Printf("Server started on port %v", cfg.Service.Port)
						go func() {
							if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Service.Port), router); err != nil {
								log.Fatalf("failed to start server: %v", err)
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						return nil
					},
				})
			}),
		)
		app.Run()
		return nil
	},
}
