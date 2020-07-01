FROM node:10 as build-stage

WORKDIR /app
COPY web/package.json /app/

# Install app requirenments
RUN npm install

# Copy rest of app
COPY ./web/src /app/src
COPY ./web/angular.json /app
COPY ./web/tsconfig.app.json /app
COPY ./web/tsconfig.json /app

RUN npm run build

FROM nginx:1.17-alpine

EXPOSE 80

COPY --from=build-stage /app/dist/web /usr/share/nginx/html

# nginx Configuration
COPY ./build/package/harbour-web/nginx.conf /etc/nginx/conf.d/default.conf

WORKDIR /

COPY ./build/package/harbour-web/docker-entrypoint.sh docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]

CMD ["nginx", "-g", "daemon off;"]
