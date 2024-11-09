package models

const (
	nginx   = iota
	haproxy = iota
)

type LoadBalancer struct {
	Type    int     //Which underlying load balancer to use, e.g: nginx, haproxy //TODO: Find a better name
	Service Service //The service instance of the load balancer
}

func ParseLoadBalancer(loadbalancer bool) *LoadBalancer {
	//TODO: Fix hardcoding
	return &LoadBalancer{
		Type: nginx,
		Service: Service{
			Image:     "nginx:1.27.2-alpine",
			Autoscale: Autoscale{Initial: 1, Autoscale: false},
			Ports: []Port{
				{HostPort: "80", ContainerPort: "80"},
				{HostPort: "443", ContainerPort: "443"},
			},
			Https: true,
			Www:   true,
		},
	}
}
