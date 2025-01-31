package shutdown

type GracefulShutdownServiceInterface interface {
	Shutdown() error
}

type gracefulShutdownService struct {
	name    string
	service GracefulShutdownServiceInterface
}

func NewGracefulShutdownService(name string, service GracefulShutdownServiceInterface) *gracefulShutdownService {
	return &gracefulShutdownService{
		name,
		service,
	}
}

func (g *gracefulShutdownService) Shutdown() error {
	return g.service.Shutdown()
}

func (g *gracefulShutdownService) ServiceName() string {
	return g.name
}
