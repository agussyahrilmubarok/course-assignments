FROM nginx:1.27-alpine
RUN rm -rf /usr/share/nginx/html/*
COPY dist /usr/share/nginx/html
# (Opsional) salin custom nginx.conf untuk Vue Router history mode
# COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
