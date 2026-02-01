FROM node:22 AS build

WORKDIR /app

COPY ./smp/package.json ./

COPY ./smp .

RUN npm i
