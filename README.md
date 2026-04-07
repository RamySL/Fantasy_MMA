# Conception d'une application Web dans le cadre du cours de Programmation concurrente, répartie et réticulaire à Sorbonne Université.

## 1. Sujet de l’application

Application Web qui consiste en une plateforme de prédiction de résultats de combats de MMA, inspirée des applications de fantasy très répandues dans le football.

## 2. API Web utilisée

https://www.thesportsdb.com/

## 3. Fonctionnalités de l’application

* Création d'un profil utilisateur.

* Avoir accès aux prochaines cartes de combats.

  * Éventuellement, avoir plusieurs ligues (ligue UFC, PFL, ARES, selon ce que l'API offre).

* Faire des prédictions sur chaque combat.

  * Éventuellement, la prédiction est modifiable jusqu'à un certain délai.

  * Plusieurs détails pour la prédiction (KO, round, etc.).

* Un classement entre tous les pronostiqueurs selon leurs scores.

* Un système de récompenses pour les mieux classés.

  * Réinitialisation à intervalles réguliers.

## 4. Technologies qui seront utilisées

* React.js pour le front

* API REST

* PostgreSQL

* Serveur écrit en Go