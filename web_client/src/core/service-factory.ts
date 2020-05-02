import {IndexService} from "./index/IndexService";
import {ErrorService} from "./ui/ErrorService";
import {HttpClient} from "./utils/HttpClient";
import {UiMessageService} from "./ui/UiMessageService";
import {RepositoryService} from "./repo/RepositoryService";

export interface ServiceMap {
    error: ErrorService;
    repository: RepositoryService;
    index: IndexService;
    uiMessages: UiMessageService;
}

let services: ServiceMap | undefined;

export function servicesMap(): ServiceMap {
    if (!services) {
        const uiMessages = new UiMessageService();
        const error = new ErrorService(uiMessages);
        const index = new IndexService(HttpClient, error);
        const repository = new RepositoryService(HttpClient, error);
        services = {
            error,
            repository,
            index,
            uiMessages
        }
    }
    return services;
}
