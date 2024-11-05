package grpc

import (
	"movieapp/gen"
	"movieapp/metadata/internal/controller/metadata"
)

// Handler definesa movie metadata gRPC handler.
type Handler struct {
	gen.UnimplementedMetadataServiceServer
	ctrl *metadata.Controller
}

// New creates a new movie metadata gRPC handler.
func New(ctrl *metadata.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}
