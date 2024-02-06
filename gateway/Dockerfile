FROM node:16-alpine

WORKDIR /app

COPY . .

RUN npm install \
&& npm run build

CMD [ "npm", "start" ]
