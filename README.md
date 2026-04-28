# Conception d'une application Web dans le cadre du cours de Programmation concurrente, répartie et réticulaire à Sorbonne Université.

## Sujet de l’application

Application Web qui consiste en une plateforme de prédiction de résultats de combats de MMA, inspirée des applications de fantasy très répandues dans le football.

## API Web utilisée

https://www.thesportsdb.com/

## Fonctionnalités de l’application

* Compte utilisateur.

* Liste de prochaines cartes de combats (à travers l'API).

  * Éventuellement, avoir plusieurs ligues (ligue UFC, PFL, ARES, selon ce que l'API offre).

* Prédictions sur chaque combat.

  * Éventuellement, la prédiction est modifiable jusqu'à un certain délai.

  * Plusieurs détails pour la prédiction (KO, round, etc.).

* Récupération des résultats réels de combats (à travers l'API).

* Comparaison des résultats réels aux prédictions.

* Un classement entre tous les pronostiqueurs selon leurs scores.

* Un système de récompenses pour les mieux classés.

  * Réinitialisation à intervalles réguliers.

## Shéma de base de données

- Utilisateur
  - id (clé prim)
  - pseudo
  - email
  - mot de passe (chiffré)

- Prédiction
  - id (clé prim)
  - utilisateur_id (Clé etr) 
  - combat_id (clé etr)
  - gagnant_prédit_id
  - points_obtenus (0 si mauvaise préd par exemple ou points_bonne_prédiction)

- Carte (req API)
  - id (clé prim)
  - titre
  - endroit (Pays, ville)
  - date
  - status

- Combat
  - ID (clé prim)
  - carte_id
  - combattant_id_1
  - combattant_id_2
  - combattant_id_gagnant (req API)
  - catégorie
  - points_bonne_prédiction (le score augmente de ce montant si la prédiction est bonne)

- Combattant
  - id
  - nom
  - prenom
  - catégorie

## Protocole 

#### Désignations des URIs

- **utilisateurs**
  - /users
  - /users/{id}
  - /users/{id}/predictions

- **prédictions**
  - /predictions
  - /predictions/{id}

- **combats**
  - /fights
  - /fights/{id}

- **cartes**
  - /cards
  - /cards/{id}
  - /cards/{cardId}/fights

- **combattant**
  - /fighters
  - /fighters/{id}

- **classement**
  - /ranking

## Technologies qui seront utilisées

* React.js pour le front

* API REST

* PostgreSQL

* Serveur écrit en Go