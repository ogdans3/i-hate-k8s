FROM nginx:alpine

# Install Certbot and necessary tools
RUN apk add --no-cache certbot certbot-nginx bash curl

ENV ACME_DIR=/var/www/certbot
RUN mkdir -p $ACME_DIR && chown -R nginx:nginx $ACME_DIR

COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
