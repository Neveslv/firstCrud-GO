# 🚀 GoBlog — Aplicação Web em Golang

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)
![Gin](https://img.shields.io/badge/Gin-Framework-000000?style=flat&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-336791?style=flat&logo=postgresql)

Aplicação Web desenvolvida em **Go (Golang)** e **Gin Framework**, com persistência em **PostgreSQL**.

O **GoBlog** é um sistema completo de gerenciamento de conteúdo (CMS) que demonstra uma arquitetura híbrida poderosa: serve tanto páginas HTML renderizadas no servidor (SSR) quanto uma API JSON estruturada, tudo no mesmo backend.

---

## 📋 Funcionalidades Principais

* 👥 **Gerenciamento de Usuários** — Cadastro completo, listagem, edição e remoção de usuários do sistema.
* 📝 **Blog System (Posts)** — Criação de publicações, leitura detalhada e moderação de conteúdo.
* 💬 **Comentários e Interação** — Sistema de comentários vinculado aos posts, permitindo discussões entre usuários.
* 📊 **Dashboard Administrativo** — Interface moderna com sidebar fixa e cards de acesso rápido.

---

## 🧩 Como Clonar e Rodar o Projeto

Siga os passos abaixo para executar o projeto localmente em sua máquina:

### 1️⃣ Clonar o repositório

Abra o terminal e execute:

bash
git clone [https://github.com/seu-usuario/firstcrud-go.git](https://github.com/seu-usuario/firstcrud-go.git)
cd firstcrud-go

### 2️⃣ Configurar o Banco de Dados

Crie um banco PostgreSQL chamado crud_go e execute o script abaixo para criar as tabelas necessárias:

```SQL


CREATE TABLE usuario (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES usuario(id) ON DELETE CASCADE,
    titulo VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE comentarios (
    id SERIAL PRIMARY KEY,
    post_id INT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES usuario(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```
### 3️⃣ Configurar Variáveis de Ambiente

Crie um arquivo .env na raiz do projeto com as credenciais do seu banco, seguindo a estrutura:

```Snippet de código

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=sua_senha
DB_NAME=crud_go
```
### 4️⃣ Executar o Projeto

Instale as dependências e rode no terminal:
```Bash

go mod tidy
go run main.go
```
Acesse em seu navegador: http://localhost:8080
