FROM node:18 AS build

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .

ARG VITE_API_URL
ARG VITE_WS_URL
ARG VITE_API_KEY
ENV VITE_API_URL=${VITE_API_URL}
ENV VITE_WS_URL=${VITE_WS_URL}
ENV VITE_API_KEY=${VITE_API_KEY}

RUN npm run build

FROM nginx:latest

COPY --from=build /app/dist /usr/share/nginx/html

COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
