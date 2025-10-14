# API - GoBank

Esta é a API backend para o projeto "GoBank", desenvolvida em Go. O foco inicial deste repositório é implementar o fluxo de *onboarding* de novos clientes, desde o cadastro inicial até a autenticação.

## 🎯 Objetivo

O objetivo é construir o core bancário de um banco digital, fornecendo uma API robusta, segura e escalável. As funcionalidades atuais incluem:

- **Início do Onboarding**: Cadastro inicial do usuário com CPF, nome e e-mail.
- **Verificação de E-mail**: Envio de um link de verificação para confirmar a posse do e-mail.
- **Conclusão do Cadastro**: Definição de senha para ativar a conta.
- **Autenticação**: Login de usuário e geração de token JWT.

## 🛠️ Tecnologias Utilizadas

- **Linguagem**: [Go](https://go.dev/)
- **Framework Web**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/) para interação com o banco de dados.
- **Banco de Dados**: [PostgreSQL](https://www.postgresql.org/)
- **Testes**:
  - Testes de unidade com o framework padrão do Go.
  - Testes de integração com [Testcontainers](https://testcontainers.com/) para instanciar um banco de dados PostgreSQL em um ambiente Docker isolado.
  - Mocks e asserções com [Testify](https://github.com/stretchr/testify).
- **Envio de E-mail**: [Resend](https://resend.com/) para e-mails transacionais.
- **Documentação da API**: [Swagger/OpenAPI 3.0](https://swagger.io/).

## 🚀 Como Executar o Projeto

### Pré-requisitos

- **Go**: Versão 1.25 ou superior.
- **Docker**: Necessário para rodar os testes de integração.

### 1. Configuração do Ambiente

Clone o repositório e, na raiz do projeto, crie um arquivo chamado `.env`. Este arquivo armazenará as variáveis de ambiente necessárias para a aplicação.

```bash
git clone https://github.com/high-effort-low-stress/go-bank-api.git
cd go-bank-api
touch .env
```

Adicione as variáveis do `.example.env` ao seu arquivo `.env`, substituindo pelos respectivos valores.


### 2. Executando a Aplicação

Com o arquivo `.env` configurado, você pode iniciar o servidor da API:

```bash
go run cmd/api/main.go
```

O servidor estará disponível em `http://localhost:8080`.

### 3. Executando os Testes

Para rodar todos os testes, incluindo os de integração que utilizam Docker, execute o comando:

```bash
go test ./...
```

## 📄 Documentação da API

A documentação completa dos endpoints está disponível no formato OpenAPI 3.0 no arquivo `docs/api/swagger.yaml`.

Você pode usar ferramentas como o Swagger Editor para colar o conteúdo do arquivo e visualizar a documentação de forma interativa.