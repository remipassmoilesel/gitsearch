import {SearchMatch, SearchResponse} from "./SearchResponse";

const sampleContent = `
            var DEFAULT_LISTEN_ADDRESS = "127.0.0.1:7777"

func LoadConfig() (Config, error) {
                    nearestGitRepo, err := getNearestRepoPath()
                    if err != nil {
                    return *new(Config), err
}`;

export class SearchService {

    public async search(query: string): Promise<SearchResponse> {
        const fakeMatches: SearchMatch[] = [
            {
                path: `/path/to/file-1`,
                commit: `commit-1`,
                content: sampleContent
            },
            {
                path: `/path/to/file-2`,
                commit: `commit-2`,
                content: sampleContent
            }
        ];

        return {
            query,
            matches: fakeMatches
        }
    }

}