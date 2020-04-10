import {SearchService} from "./SearchService";
import {ErrorService} from "./ErrorService";

export interface ServiceMap {
    search: SearchService;
    error: ErrorService;
}

let services: ServiceMap | undefined;

export function servicesMap(): ServiceMap {
    if (!services) {
        services = {
            search: new SearchService(),
            error: new ErrorService()
        }
    }
    return services;
}