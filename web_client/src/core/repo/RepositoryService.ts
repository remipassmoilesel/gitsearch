import {HttpClient} from "../utils/HttpClient";
import {RepositoryContext} from "./RepositoryContext";
import {ErrorService} from "../ui/ErrorService";
import {IndexedFile} from "../index/SearchResult";

export class RepositoryService {

    constructor(private http: HttpClient, private error: ErrorService) {
    }

    public getContext(): Promise<RepositoryContext> {
        return this.http.get("/repository/context")
            .then(res => res.data)
            .catch(err => {
                this.error.onError(err);
                return Promise.reject(err)
            })
    }

    public getDocumentById(fileId: string): Promise<IndexedFile> {
        return this.http.get<IndexedFile>(`/repository/document?id=${fileId}`)
            .then(res => res.data)
            .catch(err => {
                this.error.onError(err);
                return Promise.reject(err);
            })
    }
}