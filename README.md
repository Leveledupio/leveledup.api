# leveledup-api

docker build -t lvl-api:dev-latest . 

docker run -p 8080:8080 -v $PWD/config/dev-config.yaml:/root/config/dev-config.yaml:ro lvl-api:dev-latest

Api for Leveledup


Swagger ui tool for mock testing of the API

http://swagger.io/docs/swagger-tools/#swagger-ui-documentation-29

or if inclined you can test it in docker

docker pull swaggerapi/swagger-editor

docker run -p 80:8080 swaggerapi/swagger-editor
