import {HttpClient} from "../clients/HttpClient";
import {RepositoryContext} from "./RepositoryContext";
import {ErrorService} from "../ui/ErrorService";

export class RepositoryService {

    constructor(private httpClient: HttpClient, private errors: ErrorService) {
    }

    public getContext(): Promise<RepositoryContext> {
        return this.httpClient.get("/repository/context")
            .then(res => res.data)
            .catch(err => {
                this.errors.onError(err);
                return Promise.reject(err)
            })
    }

}