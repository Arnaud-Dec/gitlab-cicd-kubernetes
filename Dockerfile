FROM node:26.8.1-alpine3.24 AS builder
WORKDIR /app
COPY . .
RUN npm ci
RUN npm run build

FROM nginx:stable-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80

