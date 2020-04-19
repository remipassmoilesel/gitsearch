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


