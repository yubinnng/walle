# Stage 0, "build-stage", based on Node.js to build the frontend
FROM node:alpine as build
WORKDIR /app
COPY . .
RUN yarn
RUN yarn build

FROM nginx:stable-alpine
COPY --from=build /app/build /usr/share/nginx/html
COPY --from=build /app/nginx/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
