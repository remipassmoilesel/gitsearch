# TODO

1. Améliorer les tests
- Enlever tous les champs *Interface -> Interface, * pas nécéssaire

- Afficher correctement les markup languages
- Répartir les commits analysés sur plusieurs jours (1 tous les deux jours ?)        
- Bouton copier dans le presse papier la commande checkout et le path
- Filtre tri chronologiue et anté-chrono

- Coverage
- Version de l'index dans les métadonnées, si la version est différente l'index est supprimé puis reconstruit
- Option --max-depth
- Test: tri des résultat par score
- Commande clean-all pour effacer toutes les données
- Améliorer les messages d'erreurs
- Vérifier si l'index d'un projet est déjà ouvert, erreur si c'est le cas (voir doc bbolt, ajouter un timeout à l'ouverture de connexion)
- Recherche de port libre pour ouverture de plusieurs instances
- Enlever tous les fmt.Println()
- Recherche "fuzzy find" sur master en ligne de commande
- Filtre des résultat par path, par période, etc ...
- Amélioration des perfs: Réglages fins taille de batch / nbr de shaards: https://www.philipotoole.com/increasing-bleve-performance-sharding/
