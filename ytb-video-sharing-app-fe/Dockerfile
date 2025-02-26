FROM node:18 as builder

WORKDIR /app
COPY . .

# Cài đặt dependencies và build
RUN npm install && npm run build

# Sử dụng nginx để serve React app
FROM nginx:alpine
COPY --from=builder /app/build /usr/share/nginx/html
COPY --from=builder /app/build /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
