package mortar

import (
	"context"

	"github.com/go-masonry/mortar/constructors/partial"
	serverInt "github.com/go-masonry/mortar/interfaces/http/server"
	workshop "github.com/go-masonry/tutorial/complete/api"
	"github.com/go-masonry/tutorial/complete/app/controllers"
	"github.com/go-masonry/tutorial/complete/app/data"
	"github.com/go-masonry/tutorial/complete/app/services"
	"github.com/go-masonry/tutorial/complete/app/validations"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type tutorialServiceDeps struct {
	fx.In

	// API Implementations
	Workshop    workshop.WorkshopServer
	SubWorkshop workshop.SubWorkshopServer
}

func TutorialAPIsAndOtherDependenciesFxOption() fx.Option {
	return fx.Options(
		// GRPC Service APIs registration
		fx.Provide(fx.Annotated{
			Group:  partial.FxGroupGRPCServerAPIs,
			Target: tutorialGRPCServiceAPIs,
		}),
		// GRPC Gateway Generated Handlers registration
		fx.Provide(fx.Annotated{
			Group:  partial.FxGroupGRPCGatewayGeneratedHandlers + ",flatten", // "flatten" does this [][]serverInt.GRPCGatewayGeneratedHandlers -> []serverInt.GRPCGatewayGeneratedHandlers
			Target: tutorialGRPCGatewayHandlers,
		}),
		// All other tutorial dependencies
		tutorialDependencies(),
	)
}

func tutorialGRPCServiceAPIs(deps tutorialServiceDeps) serverInt.GRPCServerAPI {
	return func(srv *grpc.Server) {
		workshop.RegisterWorkshopServer(srv, deps.Workshop)
		workshop.RegisterSubWorkshopServer(srv, deps.SubWorkshop)
		// Any additional gRPC Implementations should be called here
	}
}

func tutorialGRPCGatewayHandlers() []serverInt.GRPCGatewayGeneratedHandlers {
	return []serverInt.GRPCGatewayGeneratedHandlers{
		// Register workshop REST API
		func(mux *runtime.ServeMux, endpoint string) error {
			return workshop.RegisterWorkshopHandlerFromEndpoint(context.Background(), mux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
		},
		// Register sub workshop REST API
		func(mux *runtime.ServeMux, endpoint string) error {
			return workshop.RegisterSubWorkshopHandlerFromEndpoint(context.Background(), mux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
		},
		// Any additional gRPC gateway registrations should be called here
	}
}

func tutorialDependencies() fx.Option {
	return fx.Provide(
		services.CreateWorkshopService,
		services.CreateSubWorkshopService,
		controllers.CreateWorkshopController,
		controllers.CreateSubWorkshopController,
		data.CreateCarDB,
		validations.CreateWorkshopValidations,
		validations.CreateSubWorkshopValidations,
	)
}
