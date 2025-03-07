package grpc

import (
	"context"

	pagu "github.com/pagu-project/pagu/internal/delivery/grpc/gen/go"
	"github.com/pagu-project/pagu/internal/entity"
)

type paguServer struct {
	*Server
}

func newPaguServer(server *Server) *paguServer {
	return &paguServer{
		Server: server,
	}
}

func (ps *paguServer) Run(_ context.Context, req *pagu.RunRequest) (*pagu.RunResponse, error) {
	res := ps.engine.ParseAndExecute(entity.PlatformIDWeb, req.Id, req.Command)

	return &pagu.RunResponse{
		Response: res.Message,
	}, nil
}
