# CI/CD — déploiement de la Todo App

Ce dépôt ne contient pas de code applicatif : il décrit **comment déployer** l'application Todo
(dépôts [`backend`](https://git.granux.fr/decourt-demaria/tp/backend) et
[`frontend`](https://git.granux.fr/decourt-demaria/tp/frontend)), en local avec Docker Compose et
sur Kubernetes avec Helm, ainsi que le **pipeline de CD** qui crée et supprime un environnement de
test éphémère dans le cluster.

```
docker-compose.yaml   pile complète en local (postgres, backend, frontend, swagger)
k8s/                  manifests Kubernetes « natifs », à appliquer avec kubectl
helm/                 la même application packagée en chart Helm (utilisée par la CD)
.gitlab-ci.yml        le pipeline de CD
```

`k8s/` et `helm/` décrivent la même application. On garde les deux volontairement : les manifests
natifs se lisent sans connaître Helm, la chart permet de déployer plusieurs fois l'application dans
un même namespace. **Toute modification doit être faite dans les deux.**

## Prérequis

| Outil | Version | Pourquoi |
|---|---|---|
| Docker + Compose v2 | 24+ | `docker compose` (sans tiret) pour la pile locale |
| kubectl | **1.35** | le cluster est en 1.35.2, l'écart toléré est de ±1 version mineure |
| Helm | **4.x** | la chart utilise `apiVersion: v2` ; testée avec 4.3.0 |
| Un kubeconfig | — | fourni par l'enseignant, un par namespace |

Les namespaces utilisés :

- `decourt-demaria-a` — déploiements manuels (kubeconfig par défaut `~/.kube/config`) ;
- `decourt-demaria-ci-a` — environnements créés par la CD (`~/.kube/config-ci`).

Vérifier l'accès avant toute chose :

```bash
kubectl version
kubectl get pods
```

## Lancer l'application en local

```bash
docker compose up --build
```

Les images sont construites depuis les dépôts voisins (`../backend`, `../frontend`), qui doivent
donc être clonés à côté de celui-ci. Les services démarrent dans l'ordre grâce aux `healthcheck`
et aux `depends_on: condition: service_healthy`.

| Service | URL | Remarque |
|---|---|---|
| Frontend | <http://localhost:5173> | nginx, écoute sur le port 80 dans le conteneur |
| API | <http://localhost:8080/api/todos> | pas de route `/` : tout est sous `/api` |
| Santé | <http://localhost:9090/healthz>, `/readyz`, `/metrics` | port séparé de l'API |
| Swagger UI | <http://localhost:8081> | cible `swagger` du Dockerfile backend |
| Postgres | `localhost:5432` | user/password/base : `postgres` / `postgres` / `todos` |

Arrêt et suppression des données : `docker compose down -v`.

## Déployer sur Kubernetes

### Secret d'accès à la registry (une fois par namespace)

Les images sont privées : chaque namespace a besoin d'un secret `gitlab-registry`, référencé par
`imagePullSecrets` dans tous les Deployments. Il se crée avec un **deploy token** du groupe
(GitLab → groupe `tp` → Settings → Repository → Deploy tokens, portée `read_registry`) :

```bash
read -rs TOKEN    # colle le token puis Entrée : rien ne s'affiche
kubectl create secret docker-registry gitlab-registry \
  --docker-server=registry.granux.fr \
  --docker-username=gitlab+deploy-token-4 \
  --docker-password="$TOKEN"
```

Le nom d'utilisateur est celui affiché par GitLab (`gitlab+deploy-token-<n>`), pas une adresse
e-mail. Le token ne doit jamais être écrit dans un fichier du dépôt ni collé dans un terminal
partagé.

### Version native (kubectl)

```bash
kubectl apply -f k8s/
```

Attention : `k8s/` contient **deux** descriptions de la base de données.

- `pg.yaml` — la ressource `postgresql` de l'opérateur Zalando, c'est celle qui est utilisée ;
- `pg-deployment.yaml` — l'ancien `Deployment` `postgres:16` fait à la main, **conservé pour
  mémoire et à ne pas appliquer** : il n'entre pas en conflit (il s'appelle `pg`, l'opérateur crée
  `todo-pg`), mais il ferait tourner une seconde base que personne n'utilise, sur un cluster où le
  CPU réservable est rare.

`kubectl apply -f k8s/` applique tout le dossier, y compris le fichier à ne pas appliquer :
pour un déploiement propre, appliquer les fichiers un par un, ou supprimer temporairement
`pg-deployment.yaml`.

Suppression : `kubectl delete -f k8s/`.

### Version Helm

```bash
helm install app helm/ --wait --timeout 5m
```

`app` est le nom de la release : il **préfixe toutes les ressources** (`app-backend`,
`app-frontend`, `app-pg`…), ce qui permet d'installer plusieurs fois l'application dans le même
namespace. `--wait` fait échouer la commande si les pods ne deviennent pas prêts.

Valeurs surchargeables (voir `helm/values.yaml`) :

| Valeur | Défaut | Rôle |
|---|---|---|
| `backend.imageTag` / `frontend.imageTag` | `latest` | tag des images à déployer (`backend.imageTag` sert aussi à Swagger) |
| `host.frontend` / `host.backend` / `host.swagger` | `decourt-demaria-*.k8s.granux.fr` | hosts des Ingress |
| `backend.replicas` / `frontend.replicas` / `swagger.replicas` | `1` | nombre de pods |

Exemple avec d'autres hosts et un tag précis :

```bash
helm install test helm/ \
  --set host.frontend=decourt-demaria-test-frontend.k8s.granux.fr \
  --set backend.imageTag=a1b2c3d4 --wait
```

Mise à jour : `helm upgrade app helm/`. Retour arrière : `helm rollback app`. Suppression :
`helm uninstall app`.

### Vérifier le déploiement

```bash
kubectl get pods
curl -s https://decourt-demaria-backend.k8s.granux.fr/api/todos   # -> []
curl -s -o /dev/null -w '%{http_code}\n' https://decourt-demaria-frontend.k8s.granux.fr
```

Les hosts doivent respecter le format imposé par le cours, `<groupe>-*.k8s.granux.fr`, et être
uniques dans tout le cluster. Le HTTPS fonctionne sans section `tls:` : l'ingress-nginx du cluster
sert un certificat générique pour `*.k8s.granux.fr`.

Le pod backend peut redémarrer une ou deux fois au premier démarrage : il attend la base, que
l'opérateur met une vingtaine de secondes à créer, et s'arrête après dix tentatives. Kubernetes le
relance, et il se stabilise dès que Postgres répond.

## La base de données

Elle est gérée par l'**opérateur Zalando** déjà installé sur le cluster : le dépôt ne décrit qu'une
ressource `postgresql` (`helm/templates/pg.yaml`, `k8s/pg.yaml`), et l'opérateur crée lui-même le
StatefulSet, le volume persistant, les Services et un secret de connexion par utilisateur.

Conséquences utiles à connaître :

- **aucun mot de passe dans le dépôt** : le backend lit `DB_USER` et `DB_PASSWORD` dans le secret
  `todo.<cluster>.credentials.postgresql.acid.zalan.do` généré par l'opérateur ;
- `DB_SSLMODE` vaut `require` : l'image de l'opérateur refuse les connexions non chiffrées ;
- la version est `"15"` : la CRD installée sur ce cluster n'accepte pas la 16 ;
- le nom de la ressource doit commencer par la valeur de `teamId` ;
- les données survivent à la suppression du pod (`kubectl delete pod app-pg-0`) grâce au PVC de
  1 Gi, mais le PVC est supprimé avec la ressource lors d'un `helm uninstall`.

## Le pipeline de CD

`.gitlab-ci.yml` déploie un **environnement de test éphémère** dans `decourt-demaria-ci-a` à chaque
pipeline, et sait le supprimer.

| Job | Stage | Déclenchement | Rôle |
|---|---|---|---|
| `deploy` | `deploy` | automatique | `helm install ci-<iid> helm/ --wait --timeout 5m` avec des hosts et des tags propres au pipeline |
| `cleanup` | `cleanup` | **manuel** (bouton Stop de l'environnement) | `helm uninstall ci-<iid>` |

Points de conception :

- **une release par pipeline** (`RELEASE: ci-$CI_PIPELINE_IID`, hosts préfixés de la même façon) :
  deux pipelines simultanés ne se gênent pas ;
- `CI_PIPELINE_IID` et non `CI_PIPELINE_ID` : `CI_PIPELINE_ID` n'existe qu'au lancement du job,
  alors que GitLab résout `environment:name` à la création du pipeline — le nom de l'environnement
  serait tronqué en `ci-` ;
- l'environnement **reste disponible** après le pipeline (`cleanup` en `when: manual`), pour pouvoir
  tester l'application déployée ; il s'arrête tout seul au bout d'un jour (`auto_stop_in`) ;
- si `deploy` échoue, l'`after_script` désinstalle la release : **le cluster n'est jamais laissé
  sale**, même quand l'image n'existe pas ou qu'un pod ne démarre pas ;
- l'image du pipeline est `alpine/helm` avec `entrypoint: [""]`, sinon GitLab ne peut pas lancer de
  `script` (le point d'entrée de l'image est `helm` lui-même).

Les environnements sont visibles dans GitLab sous **Operate → Environments** : bouton *Open* vers le
frontend déployé, bouton *Stop* pour tout supprimer.

### Déclenchement depuis backend et frontend

Les pipelines de `backend` et `frontend` finissent par un job `deploy` de type `trigger` qui lance
ce pipeline-ci :

```yaml
deploy:
  stage: deploy
  variables:
    BACKEND_TAG: $CI_COMMIT_SHORT_SHA   # FRONTEND_TAG côté frontend
  trigger:
    project: decourt-demaria/tp/ci-cd
    branch: main
    strategy: depend
```

La variable est transmise au pipeline aval, qui déploie donc **l'image que la CI vient de
construire** et non `latest`. `strategy: depend` fait prendre au job amont le statut du
déploiement : si le déploiement casse, le pipeline du backend est rouge.

### Configuration GitLab nécessaire

À refaire si le projet est recréé ailleurs :

1. **Variable `KUBECONFIG`** (Settings → CI/CD → Variables) : type **File**, visibilité
   **Visible**, **Protected** — elle contient le kubeconfig du namespace CI. Helm et kubectl lisent
   automatiquement cette variable, le pipeline n'a rien à faire. Elle ne peut pas être *masquée*
   (GitLab refuse les valeurs contenant des espaces), donc **ne jamais l'afficher** dans un job
   (`cat $KUBECONFIG` est interdit). Comme elle est protégée, elle n'existe que sur une branche
   protégée.
2. **Branche `main` protégée** dans ce projet, sinon le job `deploy` ne reçoit pas la variable et
   tente de joindre `localhost:8080`.
3. **Secret `gitlab-registry`** dans les deux namespaces (voir plus haut), avec un deploy token
   par namespace.

## Choix et écarts assumés

- **`requests` = 100m CPU / 100Mi de mémoire** sur tous les conteneurs, et `limits` à 100Mi de
  mémoire (identique à la requête) et 500m de CPU. La mémoire se partage mal, d'où une limite égale
  à la réservation ; le CPU se partage bien, d'où une limite plus large pour absorber le démarrage.
- Les **probes** interrogent les ports de santé directement sur le pod (`/readyz` et `/healthz` sur
  9090 pour le backend), ce qui n'oblige pas à exposer ce port dans le Service.
- Le **HPA** a été étudié mais n'est pas déployé : le sujet demande seulement de « regarder pour le
  mettre en place », et sur un cluster partagé où le CPU réservable est rare, autoriser des réplicas
  supplémentaires revient à réserver du CPU au détriment des autres groupes. Les conditions
  techniques sont réunies (`metrics-server` installé, `requests` CPU définies sur tous les
  conteneurs) : il suffirait d'un `HorizontalPodAutoscaler` en `autoscaling/v2` ciblant le
  Deployment du backend, de 1 à 3 réplicas, sur une utilisation CPU moyenne de 70 %.
- Le bouton *Try it out* de Swagger UI ne fonctionne pas : la spécification générée ne contient pas
  de section `servers`, donc l'interface ne sait pas à quelle URL envoyer les requêtes.
