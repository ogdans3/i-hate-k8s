import { PublicAutoscale, PublicContainerEngine, PublicPort, PublicProject, PublicService } from "../public-models/public-models.ts";
export interface Project {
    project: string;
    engine: ContainerEngine;

    logging: boolean;
    analytics: boolean;
    dashboard: boolean;
    registry: boolean;
    loadbalancer: boolean;

    services: Service[];
}

export enum ContainerEngine {
    docker,
}

export interface Service {
    name: string;
    image: string;
    www: boolean;
    https: boolean;
    ports: Port[];
    autoscale: Autoscale;
}

export interface Port {
    protocol: string;
    hostPort: string;
    containerPort: string;
}

export interface Autoscale {
    initial: number;
    autoscale: boolean;
}

export function convertPublicModelToInternalModel(publicProject: PublicProject) {
    const project: Project = {
        project: publicProject.project,
        engine: convertEngine(publicProject.engine),

        logging: publicProject.logging,
        analytics: publicProject.analytics,
        dashboard: publicProject.dashboard,
        registry: publicProject.registry,
        loadbalancer: convertLoadbalancer(publicProject.loadbalancer),

        services: publicProject.services.map(convertService),
    };
    return project;
}

function convertEngine(_engine: PublicContainerEngine) {
    return _engine as unknown as ContainerEngine;
}
function convertLoadbalancer(_loadbalancer: boolean) {
    return _loadbalancer as boolean;
}
function convertPort(_port: PublicPort) {
    return _port as Port;
}
function convertService(_service: PublicService) {
    return <Service>{
        name: _service.name,
        image: _service.image,
        www: _service.www,
        https: _service.https,
        ports: _service.ports.map(convertPort),
        autoscale: convertAutoscale(_service.autoscale),
    }
}
function convertAutoscale(_autoscale: PublicAutoscale) {
    return _autoscale as Autoscale;
}
