# Git search

Projet d'apprentissage de Go.

## TODO

- Vérifier si l'index d'un projet est déjà ouvert, erreur si c'est le cas
- Améliorer l'indexation:
    - Une goroutine par commit
    - Batch index
- Inclure tous les champs dans les résultats de recherche bleve: https://stackoverflow.com/questions/50572998/how-to-use-golang-bleve-search-results
- Recherche de port libre pour ouverture de plusieurs instances
- Enlever tous les fmt.Println()
- Syntaxe colorées avec highlight.js
- Toutes les données dans ~/.gitsearch, y compris le client web
- Variable d'environnement optionnelle pour ~/.gitsearch 
- Recherche "fuzzy find" sur master en ligne de commande
- Recherche textuelle dans les fichiers actuels
- Recherche textuelle dans l'historique
- Interface Web minimaliste
- Le client web doit être embarqué dans un tar, extrait dans ~/.gitsearch si absent. Voir: https://github.com/markbates/pkger
