FROM node:22-alpine

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run compile || true

EXPOSE 8545

CMD ["npm", "run", "node"]
