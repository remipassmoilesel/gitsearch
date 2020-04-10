
export interface SearchResponse {
    query: string;
    matches: SearchMatch[];
}

export interface SearchMatch {
    path: string;
    commit: string;
    content: string;
}