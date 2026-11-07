FROM node:26.8.1-alpine3.24 AS base
WORKDIR /app
COPY package.json package-lock.json ./
# --ignore-scripts : n'exécute aucun script postinstall de dépendance (vecteur de
# l'attaque axios de mars 2026) ; aucune dépendance du projet n'en a besoin
RUN npm ci --ignore-scripts
COPY . .

# Check du README, lancé avec `docker build --target test .`
# Un build normal (image finale, docker compose) ne passe pas par cette étape.
FROM base AS test
RUN npm test

FROM base AS builder
RUN npm run build

FROM nginx:stable-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
