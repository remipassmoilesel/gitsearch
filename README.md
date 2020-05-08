# Gitsearch

[![pipeline status](https://gitlab.com/remipassmoilesel/gitsearch/badges/master/pipeline.svg)](https://gitlab.com/remipassmoilesel/gitsearch/-/commits/master)
[![go_report](https://goreportcard.com/badge/gitlab.com/remipassmoilesel/gitsearch)](https://goreportcard.com/report/gitlab.com/remipassmoilesel/gitsearch)
[![coverage report](https://gitlab.com/remipassmoilesel/gitsearch/badges/master/coverage.svg)](https://gitlab.com/remipassmoilesel/gitsearch/-/commits/master)


Search in your git repositories !   

Start gitsearch in a Git repository and it will index all committed files and search for you.   

**This is a learning project, paint is fresh, and contributions are welcome !**

Source code is hosted here:
- [https://gitlab.com/remipassmoilesel/gitsearch](https://gitlab.com/remipassmoilesel/gitsearch)
- [https://github.com/remipassmoilesel/gitsearch](https://github.com/remipassmoilesel/gitsearch)


## Installation

Gitsearch is a single executable. Drop it in a folder in your path and start use it.

TBC...


## Usage

Use the web user interface:

    $ cd /path/to/git-repository
    $ gitsearch
    

Search the current files and the history of a git repository:

    $ gitsearch search -q <query> -n <number-of-results>


You can use Bleve Query strings: https://blevesearch.com/docs/Query-String-Query/ 
    

Show help for commands or sub commands:

    $ gitsearch help
    $ gitsearch index help
    

## Development

For development and testing see:

- [https://github.com/blevesearch/bleve](https://github.com/blevesearch/bleve)
- [https://github.com/go-git/go-git](https://github.com/go-git/go-git)
- [https://github.com/golang/mock](https://github.com/golang/mock)


Build (POSIX OS only):

    $ ./scripts/ci.sh
