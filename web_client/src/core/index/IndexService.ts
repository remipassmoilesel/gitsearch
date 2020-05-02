import {IndexedFile, SearchResult} from "./SearchResult";
import {ErrorService} from "../ui/ErrorService";
import {HttpClient} from "../utils/HttpClient";
import {BuildOperationResult, CleanOperationResult} from "./IndexOperationResult";

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

    public async buildIndex(): Promise<BuildOperationResult> {
        return this.http.post<BuildOperationResult>(`/index/build`)
            .then(res => res.data)
            .catch(err => {
                this.error.onError(err);
                return Promise.reject(err);
            })
    }

    public async cleanIndex(): Promise<CleanOperationResult> {
        return this.http.post<CleanOperationResult>(`/index/clean`)
            .then(res => res.data)
            .catch(err => {
                this.error.onError(err);
                return Promise.reject(err);
            })
    }

}
