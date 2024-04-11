<h4 align="center">Rate limit middleware by: <a href="https://www.linkedin.com/in/matheuslopes1999/" target="_blank">Matheus Lopes</a>.</h4>
<p align="center">Desafio proposto no módulo de desafio técnico do curso de pós graduação Pós Go Expert da Fullcycle</p>

<p align="center">
  <img src="https://img.shields.io/badge/Tests-Passing-2ea44f" alt="Tests - Passing">
  <img src="https://img.shields.io/badge/Go-1.21.x-2ea44f" alt="Go - 1.21.x">
</p>

<p align="center">
  <a href="#como-rodar-a-imagem-docker">Como rodar a imagem docker</a> •
  <a href="#interagindo-com-o-app">Interagindo com o app</a> •
  <a href="#consultando-a-api">Como Configurar as opções do rate limit</a> •
</p>

## Como rodar a imagem docker

Para clonar essa applicação você precisará ter instalado o [git](https://git-scm.com) e o [golang](https://go.dev/) em sua máquina. Insira os seguintes commandos em sua CLI para iniciar e rodar a instancia docker:

```bash
# Clone este repositório
$ git clone https://github.com/Nimbo1999/go-rate-limiter-middleware.git

# Navegue no repositório
$ cd go-rate-limiter-middleware

# Inicie o projeto com o docker
$ docker compose up -d --build
# Navegue na pasta /api para executar uma requisição HTTP ou realize um GET para o serviço que ficará disponível por padrão na porta :8080
```

## Interagindo com o app

Assim que o projeto estiver rodando com todos os container docker, você já poderá interagir com a applicação por meio de uma requisição http GET. Esse endpoint não espera receber nenhum parâmetro e caso você não tenha atingido o limite de requisições, ele retornará uma simples resposta com o status 200 OK e uma string em seu conteúdo com a seguinte mensagem "Welcome to the app!"

```bash
$ curl -X GET \
-H "API_KEY: token" \
http://localhost:8080
```

Ou se preferir, caso tenha instalado a extensão [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client), abra a pasta `/api` e fique avontade para alterar o cep e executar a requisição.

## Como Configurar as opções do rate limit

O projeto possui um arquivo `.env` na raiz do projeto. Esse arquivo contem todas as variáveis de ambiente utilizadas pelo projeto para configuração do mesmo.

> **Atenção** Com exceção da variável `REDIS_HOST`, todas as outras variáves vão ser refletidas na applicação. A variável `REDIS_HOST` é sobrescrita pelo arquivo `docker-compose.yaml` para injetar o container do redis como host.

Assim que você iniciar o build da aplicação, o Dockerfile disponibiliza esse arquivo de configuração para o container, e portanto, a aplicação consegue consultar os valores para utilizar em sua configuração.

> **Atenção** Caso quiser rodar o projeto com outra configuração, será necessário alterar o arquivo `.env` e rodar o comando `docker compose up -d --build` para fazer o build da imagem com os novos valores.
