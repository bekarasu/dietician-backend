package healthcheck

var healthCheck *config

const (
	LivenessStatusOk       = "UP"
	LivenessStatusShutdown = "SHUTDOWN"
	ReadinessStatusOk      = "READY"
)

type IService interface {
	Name() string
	Check() bool
}

type config struct {
	serverUp bool
	services []IService
}

func InitHealthCheck(svs ...IService) {
	healthCheck = &config{
		serverUp: true,
		services: svs,
	}
}

func Readiness() map[string]bool {
	services := make(map[string]bool, len(healthCheck.services))
	for _, service := range healthCheck.services {
		services[service.Name()] = service.Check()
	}

	return services
}

func Liveness() bool {
	return healthCheck.serverUp
}

func ServerShutdown() {
	healthCheck.serverUp = false
}

func IsConnectionSuccessful(conn map[string]bool) bool {
	for _, status := range conn {
		if !status {
			return false
		}
	}

	return true
}
