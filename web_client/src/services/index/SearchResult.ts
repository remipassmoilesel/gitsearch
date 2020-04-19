export interface SearchResult {
    Query: string;
    TookUs: number;
    Matches: SearchMatch[];
}

export interface SearchMatch {
    File: IndexedFile;
    Fragments: string[]
}

export interface IndexedFile {
    Id: string
    Hash: string
    Commit: string
    Content: string
    Path: string
}
