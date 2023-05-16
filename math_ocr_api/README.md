# Math OCR API

Service that will translate images of math equations to TeX. Uses services to compress the images and to convert them to text. Contains middleware to use the same auth system as the main api.

## Environment variables

Copy the example file:

```bash
cp .env.example .env
```

Or if using docker compose:

```bash
cp .env.example .env.docker
```

## Install dependencies

The project uses yarn. Please don't mix yarn and npm.

```bash
yarn
```

## Format code

Auto format code with prettier. Configuration in `.prettierrc.json`.

## Run the app in development mode

### With automatic restart on changes

```bash
yarn dev
```

### Without automatic restart on changes

```bash
yarn dev-start
```

## Run the app in production mode

The service is written using TypeScript. But the production version is compiled to JavaScript. The ts-node package has a higher footprint in production and uses more memory.

### Compile the app from TypeScript to JavaScript

```bash
yarn prod-build
```

### Run the app

```bash
yarn prod-start
```
