FROM node:22.18.0 AS build

WORKDIR /app
COPY package*.json ./
RUN corepack enable && corepack prepare pnpm@latest --activate && pnpm install

COPY . .
RUN pnpm run build

FROM nginx:1.27.2-alpine3.20

RUN rm -rf /usr/share/nginx/html/*
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 8080
CMD ["nginx", "-g", "daemon off;"]
