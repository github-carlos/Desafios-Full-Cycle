
# Desafio 2 Módulo Docker

## Executando
Com o docker-compose instalado, execute:
```
docker-compose up -d
```
E acesse o localhost:80 para ver a aplicaçao em execuçao

## Descrição
Criar um ambiente com as seguintes características
- Um servidor Nginx que aponte para uma aplicaçao node
- Uma aplicaçao Node que depende um banco de dados MYSQL

Para fazer com que a aplicaçao inicie apenas quando o banco estiver disponível foi utilizado a ferramenta Dockerize.