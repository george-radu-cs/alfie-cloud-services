# Alfie Cloud Services

Services for the Alfie project. This repo uses docker compose V2.
The compose file contains definitions for the following services:

- [Alfie API](./api/README.md)
- [Math OCR API](./math_ocr_api/README.md)
- [A PostgreSQL database](./db/README.md)

## Client app

The client app is a separate repo. It can be found [here](https://github.com/george-radu-cs/alfie-client).

## Environment variables

Example from `.env.example`

```bash
cp .env.example .env
```

## Docker for development

```bash
docker compose --env-file .env up -d --build
```

## Docker for production

In addition to the `compose.yaml` file, the production version uses one more file: `compose.prod.yaml`. This file contains the configuration for the services to use a certificate for secure connections.

```bash
docker compose --env-file .env -f compose.yaml -f compose.prod.yaml up -d --build
docker image prune -f
docker system prune -f
```

## GitHub pipeline for production deploy

Deploy to production on push to `main` branch. Deployment is done with docker on an ec2 instance via ssh. Requieres the following secrets to be set in GitHub's repo secrets: `AWS_PRIVATE_KEY`, `AWS_USERNAME`, `AWS_HOSTNAME`
