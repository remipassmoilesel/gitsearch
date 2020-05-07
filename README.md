# Git search

Search in your git repositories !

## Installation

TBD


## Usage

Search the current files and the history of a git repository:

    $ gitsearch search <query>


You can use Bleve Query strings: https://blevesearch.com/docs/Query-String-Query/ 
    

Show help for commands or sub commands:

    $ gitsearch help
    $ gitsearch index help
    

## Development

You will need pkger in your PATH:

    $ go get github.com/markbates/pkger/cmd/pkger


Build (Linux, Mac only):

    $ ./scripts/ci.sh
    

For development and testing see:

- https://github.com/blevesearch/bleve
- https://github.com/go-git/go-git
- https://github.com/golang/mock
