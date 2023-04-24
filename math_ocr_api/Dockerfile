FROM node:16-alpine

WORKDIR /usr/src/math_ocr_api

# pre-copy/cache package.json for pre-downloading dependencies and only redownloading
COPY package.json ./
COPY yarn.lock ./

# install dependencies
RUN yarn

# copy the source code
COPY . .
# copy the env file to the root of the project
COPY .env.docker .env

# build the app for production
RUN yarn prod-build

# run the app
CMD ["yarn", "prod-start"]