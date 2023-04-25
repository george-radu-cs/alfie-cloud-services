# Alfie docker

## Docker for development

```bash
docker compose --env-file .env up -d --build
```

## Docker for production

```bash
docker compose --env-file .env -f compose.yaml -f compose.prod.yaml up -d --build
docker image prune -f
docker system prune -f
```
