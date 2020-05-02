export interface SearchResult {
    Query: string;
    TookMs: number;
    Matches: SearchMatch[];
}

export interface SearchMatch {
    File: IndexedFile;
    Fragments: string[]
}

export interface IndexedFile {
    Id: string
    Hash: string
    Date: string;
    Commit: string
    Name: string
    Content: string
    Path: string
}
