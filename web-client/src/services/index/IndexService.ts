import {SearchResult} from "./SearchResult";
import {ErrorService} from "../ui/ErrorService";
import {IndexOperationResult} from "./IndexOperationResult";
import {HttpClient} from "../clients/HttpClient";

export class IndexService {

    constructor(private http: HttpClient, private error: ErrorService) {
    }

    public search(query: string): Promise<SearchResult> {
        return this.http.get<SearchResult>(`/search?query=${query}`)
            .then(res => res.data)
            .catch(err => {
                this.error.onError(err);
                return Promise.reject(err);
            })
    }

    public async buildIndex(): Promise<IndexOperationResult> {
        return this.http.post<IndexOperationResult>(`/index/build`)
            .then(res => res.data)
            .catch(err => {
                this.error.onError(err);
                return Promise.reject(err);
            })
    }

    public async cleanIndex(): Promise<IndexOperationResult> {
        return this.http.post<IndexOperationResult>(`/index/clean`)
            .then(res => res.data)
            .catch(err => {
                this.error.onError(err);
                return Promise.reject(err);
            })
    }
}
