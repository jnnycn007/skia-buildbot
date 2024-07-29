package coverage

import (
	"go.skia.org/infra/go/coverage/coveragestore"
	"go.skia.org/infra/go/sklog"

	pb "go.skia.org/infra/go/coverage/proto/v1"
	"google.golang.org/grpc"
)

// coverageService implements backend.BackendService, provides a wrapper struct
// for the coverage service implementation.
type coverageService struct {
	pb.UnimplementedCoverageServiceServer
	coverageStore coveragestore.Store
}

// New returns a new instance of the coverage service.
func NewCoverageService(coverageStore coveragestore.Store) *coverageService {
	return &coverageService{
		coverageStore: coverageStore,
	}
	/*
	   ctx := context.Background()
	   db, err := pgxpool.Connect(ctx, config.GetConnectionString())

	   	if err != nil {
	   		sklog.Fatalf("Failed to connect to DB: %s", config.GetConnectionString())
	   		return nil
	   	}

	   sqlcoveragestore, err := sqlcoveragestore.New(db)

	   	if err != nil {
	   		return nil
	   	}

	   	return &coverageService{
	   		CoverageServiceServer: coverage_service.New(*sqlcoveragestore),
	   	}
	*/
}

// RegisterGrpc registers the grpc service with the server instance.
func (service *coverageService) RegisterGrpc(grpcServer *grpc.Server) {
	sklog.Infof("Register Coverage Service")
	pb.RegisterCoverageServiceServer(grpcServer, service)
}

// GetServiceDescriptor returns the service descriptor for the service.
func (service *coverageService) GetServiceDescriptor() grpc.ServiceDesc {
	return pb.CoverageService_ServiceDesc
}
