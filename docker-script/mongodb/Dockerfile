FROM mongo:4.4

RUN apt-get update && apt-get install -y wget

# install migration
RUN wget https://github.com/BlackMocca/migrate/releases/download/v5.0/migrate.linux-amd64
RUN chmod +x migrate.linux-amd64
RUN mv migrate.linux-amd64 /usr/local/bin/migrate

RUN pwd
ADD ./init /docker-entrypoint-initdb.d/
