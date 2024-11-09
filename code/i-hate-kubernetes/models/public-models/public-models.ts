export interface PublicProject {
    project: string;
    engine: PublicContainerEngine;

    logging: boolean;
    analytics: boolean;
    dashboard: boolean;
    registry: boolean;
    loadbalancer: boolean;

    services: PublicService[];
}

export enum PublicContainerEngine {
    docker,
}

export interface PublicService {
    name: string;
    image: string;
    www: boolean;
    https: boolean;
    ports: (string | PublicPort)[];
    autoscale: boolean | PublicAutoscale;
}

export interface PublicPort {
    protocol: string;
    hostPort: string;
    containerPort: string;
}

export interface PublicAutoscale {
    initial: number;
    autoscale: boolean;
}
