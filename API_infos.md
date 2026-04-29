# API ESPN utilisée dans Fantasy MMA

L'API ESPN sert à :

1. récupérer les cartes de combats avec la liste des combats ;
2. récupérer les résultats d'une carte une fois qu'elle est terminée.

---

## 1. Nature de l'API

L'API ESPN utilisée ici est une API JSON publique mais non officiellement documentée par ESPN.

Cependant elle possède une documentation non officiel très riche : 
https://github.com/pseudo-r/Public-ESPN-API/blob/main/docs/sports/mma.md
Elle ne nécessite pas de clé API.

Endpoint principal utilisé pour l'UFC :

```http
GET https://site.api.espn.com/apis/site/v2/sports/mma/ufc/scoreboard
```

Pour une autre ligue MMA exposée par ESPN, il faut remplacer le slug de ligue :

```http
GET https://site.api.espn.com/apis/site/v2/sports/mma/{league}/scoreboard
```

Exemples :

```http
GET https://site.api.espn.com/apis/site/v2/sports/mma/ufc/scoreboard
GET https://site.api.espn.com/apis/site/v2/sports/mma/pfl/scoreboard
```
---

## 2. Besoin 1 : lister les cartes avec tous les combats

### 2.1. Récupérer le calendrier des cartes

Endpoint :

```http
GET https://site.api.espn.com/apis/site/v2/sports/mma/ufc/scoreboard
```

Cet endpoint retourne notamment :

```json
{
  "leagues": [
    {
      "calendar": [
        {
          "label": "UFC Fight Night: Della Maddalena vs. Prates",
          "startDate": "2026-05-02T11:00Z",
          "endDate": "2026-05-03T06:59Z",
          "event": {
            "$ref": "http://sports.core.api.espn.pvt/v2/sports/mma/leagues/ufc/events/600058807?lang=en&region=us"
          }
        }
      ]
    }
  ]
}
```

Champs utiles à lire :

| Information voulue | Chemin JSON |
|---|---|
| Nom de la carte | `leagues[0].calendar[i].label` |
| Date de début de la carte | `leagues[0].calendar[i].startDate` |
| Date de fin estimée | `leagues[0].calendar[i].endDate` |
| Identifiant ESPN de la carte | à extraire depuis `leagues[0].calendar[i].event.$ref` |
| Ligue | `leagues[0].abbreviation` |
| Nom complet de la ligue | `leagues[0].name` |
| Saison | `leagues[0].season.year` |

L'identifiant de la carte n'est pas fourni directement dans un champ simple. Il faut l'extraire de l'URL `$ref`.

Exemple :

```text
http://sports.core.api.espn.pvt/v2/sports/mma/leagues/ufc/events/600058807?lang=en&region=us
```

Identifiant à garder :

```text
600058807
```

Dans ta table `Carte`, tu peux stocker cet identifiant dans un champ comme :

```text
external_id = "600058807"
```

---

### 2.2. Récupérer les combats d'une carte précise

Le calendrier donne les dates des cartes. Pour récupérer les combats d'une carte précise, il faut appeler le même endpoint avec le paramètre `dates`.

Format :

```http
GET https://site.api.espn.com/apis/site/v2/sports/mma/ufc/scoreboard?dates=YYYYMMDD
```

Exemple pour une carte du 2 mai 2026 :

```http
GET https://site.api.espn.com/apis/site/v2/sports/mma/ufc/scoreboard?dates=20260502
```

Exemple `curl` :

```bash
curl "https://site.api.espn.com/apis/site/v2/sports/mma/ufc/scoreboard?dates=20260502"
```

Dans la réponse, les cartes sont dans :

```text
events[]
```

Et les combats de chaque carte sont dans :

```text
events[i].competitions[]
```

Exemple de structure simplifiée :

```json
{
  "events": [
    {
      "id": "600058807",
      "name": "UFC Fight Night: Della Maddalena vs. Prates",
      "date": "2026-05-02T08:00Z",
      "competitions": [
        {
          "id": "401863149",
          "date": "2026-05-02T11:00Z",
          "type": {
            "id": "969",
            "abbreviation": "Welterweight"
          },
          "venue": {
            "fullName": "RAC Arena (AUS)",
            "address": {
              "city": "Perth",
              "state": "WA",
              "country": "Australia"
            }
          },
          "competitors": [
            {
              "id": "4294832",
              "order": 2,
              "winner": false,
              "athlete": {
                "fullName": "Carlos Prates",
                "displayName": "Carlos Prates",
                "shortName": "C. Prates"
              },
              "records": [
                {
                  "summary": "22-7-0"
                }
              ]
            }
          ],
          "status": {
            "type": {
              "name": "STATUS_SCHEDULED",
              "state": "pre",
              "completed": false,
              "description": "Scheduled"
            }
          }
        }
      ]
    }
  ]
}
```

Champs utiles pour créer une `Carte` :

| Information voulue | Chemin JSON |
|---|---|
| ID externe de la carte | `events[i].id` |
| Nom de la carte | `events[i].name` |
| Nom court de la carte | `events[i].shortName` |
| Date principale | `events[i].date` |
| Saison | `events[i].season.year` |

Champs utiles pour créer un `Combat` :

| Information voulue | Chemin JSON |
|---|---|
| ID externe du combat | `events[i].competitions[j].id` |
| Date du combat | `events[i].competitions[j].date` |
| Catégorie | `events[i].competitions[j].type.abbreviation` |
| ID catégorie ESPN | `events[i].competitions[j].type.id` |
| Nombre de rounds prévus | `events[i].competitions[j].format.regulation.periods` |
| Statut du combat | `events[i].competitions[j].status.type.name` |
| État du combat | `events[i].competitions[j].status.type.state` |
| Combat terminé ? | `events[i].competitions[j].status.type.completed` |
| Description du statut | `events[i].competitions[j].status.type.description` |

Champs utiles pour le lieu :

| Information voulue | Chemin JSON |
|---|---|
| Nom de la salle | `events[i].competitions[j].venue.fullName` |
| Ville | `events[i].competitions[j].venue.address.city` |
| État / région | `events[i].competitions[j].venue.address.state` |
| Pays | `events[i].competitions[j].venue.address.country` |

Champs utiles pour les combattants :

| Information voulue | Chemin JSON |
|---|---|
| ID externe du combattant | `events[i].competitions[j].competitors[k].id` |
| Ordre d'affichage | `events[i].competitions[j].competitors[k].order` |
| Nom complet | `events[i].competitions[j].competitors[k].athlete.fullName` |
| Nom d'affichage | `events[i].competitions[j].competitors[k].athlete.displayName` |
| Nom court | `events[i].competitions[j].competitors[k].athlete.shortName` |
| Pays du combattant | `events[i].competitions[j].competitors[k].athlete.flag.alt` |
| URL du drapeau | `events[i].competitions[j].competitors[k].athlete.flag.href` |
| Record MMA | `events[i].competitions[j].competitors[k].records[0].summary` |
| Gagnant du combat | `events[i].competitions[j].competitors[k].winner` |

---

## 3. Besoin 2 : récupérer les résultats d'une carte terminée

Pour récupérer les résultats, on réutilise le même endpoint que pour la carte précise après la fin du combat.

---

### 3.1. Vérifier qu'un combat est terminé

Pour chaque combat :

```text
events[i].competitions[j]
```

Lire :

| Information voulue | Chemin JSON |
|---|---|
| Combat terminé ? | `events[i].competitions[j].status.type.completed` |
| État du combat | `events[i].competitions[j].status.type.state` |
| Nom du statut | `events[i].competitions[j].status.type.name` |
| Description du statut | `events[i].competitions[j].status.type.description` |

Un combat terminé devrait avoir une structure du type :

```json
{
  "status": {
    "type": {
      "name": "STATUS_FINAL",
      "state": "post",
      "completed": true,
      "description": "Final"
    }
  }
}
```


---

### 4.2. Trouver le gagnant du combat

Dans chaque combat terminé, lire les deux combattants :

```text
events[i].competitions[j].competitors[]
```

Le gagnant est celui qui a :

```text
events[i].competitions[j].competitors[k].winner == true
```

Champs à utiliser :

| Information voulue | Chemin JSON |
|---|---|
| ID externe du gagnant | `events[i].competitions[j].competitors[k].id` |
| Nom du gagnant | `events[i].competitions[j].competitors[k].athlete.displayName` |
| Booléen gagnant | `events[i].competitions[j].competitors[k].winner` |

---

## 5. Mapping avec la base de données

### Table `Carte`

| Colonne dans ta base | Source ESPN |
|---|---|
| `id` | ID interne PostgreSQL |
| `external_id` | `events[i].id` ou ID extrait depuis `calendar[i].event.$ref` |
| `titre` | `events[i].name` ou `calendar[i].label` |
| `date` | `events[i].date` ou `calendar[i].startDate` |
| `status` | `events[i].status.type.completed` |
| `lieu_nom` | `events[i].competitions[0].venue.fullName` |
| `ville` | `events[i].competitions[0].venue.address.city` |
| `region` | `events[i].competitions[0].venue.address.state` |
| `pays` | `events[i].competitions[0].venue.address.country` |

---

### Table `Combat`

| Colonne dans ta base | Source ESPN |
|---|---|
| `id` | ID interne PostgreSQL |
| `external_id` | `events[i].competitions[j].id` |
| `carte_id` | ID interne de la carte liée |
| `combattant_id_1` | premier combattant dans `events[i].competitions[j].competitors[0].id` |
| `combattant_id_2` | deuxième combattant dans `events[i].competitions[j].competitors[1].id` |
| `combattant_id_gagnant` | combattant avec `competitors[k].winner == true` |
| `categorie` | `events[i].competitions[j].type.abbreviation` |
| `points_bonne_prediction` | défini par backend |

---

### Table `Combattant`

| Colonne dans ta base | Source ESPN |
|---|---|
| `id` | ID interne PostgreSQL |
| `external_id` | `competitors[k].id` |
| `nom_complet` | `competitors[k].athlete.fullName` |
| `record` | `competitors[k].records[0].summary` |

---