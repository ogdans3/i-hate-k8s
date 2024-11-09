interface Project {
    project: string;
    engine: ContainerEngine;

    logging: boolean;
    analytics: boolean;
    dashboard: boolean;
    registry: boolean;
    loadbalancer: boolean;

    services: Service[];
}

enum ContainerEngine {
    docker,
}

interface Service {
    name: string;
    image: string;
    www: boolean;
    https: boolean;
    ports: (string | Port)[];
    autoscale: boolean | Autoscale;
}

interface Port {
    protocol: string;
    hostPort: string;
    containerPort: string;
}

interface Autoscale {
    initial: number;
    autoscale: boolean;
}
