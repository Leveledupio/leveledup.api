# leveledup.api

API for Leveledup.io

API configuration file 

```yaml
port: "8080"
dsn: "USERNAME:PASSWORD@@tcp(HOST:3306)/DATABASE"
http_addr: ":8888"
http_drain_interval: 1s
MYSQL_TCP_PORT: "3306"
MYSQL_USER: "USERNAME"
MYSQL_PWD: "PASSWORD"

```

You need to mount that config file for docker. The API will panic if it doesnt find the config file. Then will try to connect to database, and also will exit if it can't connect. 

[How to connect to the Database](https://github.com/strongjz/leveledup.api/blob/master/Database.md)

Then run the api in docker

``` bash
docker build -t lvl-api:dev-latest . 

docker run -p 8080:8080 -v $PWD/config/config.yaml:/root/config/config.yaml:ro lvl-api:dev-latest
```

Swagger ui tool for mock testing of the API

http://swagger.io/docs/swagger-tools/#swagger-ui-documentation-29

or if inclined you can test it in docker
```bash
cd swagger/

docker pull swaggerapi/swagger-editor

docker run -p 80:8080 swaggerapi/swagger-editor
```