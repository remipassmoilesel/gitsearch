export interface Language {
    extensions: string[]
    name: string
}

export const Languages: Language[] = [
    {
        extensions: ["js"],
        name: "javascript"
    },
    {
        extensions: ["ts"],
        name: "typescript"
    },
    {
        extensions: ["go"],
        name: "go"
    },
    {
        extensions: ["sh"],
        name: "bash"
    }
]

export function languageFromFileName(path: string): string {
    const extension = path.split('.').pop();
    if (!extension) {
        return '';
    }

    const language = Languages.find(lang => lang.extensions.indexOf(extension) !== -1);
    return language?.name || '';
}
