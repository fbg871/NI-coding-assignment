FROM node:latest

WORKDIR komp-registry
COPY . .

RUN npm install

ENV NEXT_PUBLIC_IS_COMPOSE=true
RUN npm run build

ENTRYPOINT ["npm", "run", "start"]