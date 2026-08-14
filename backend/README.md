# LostEmp API

Backend REST (mockup/MVP) para gestão de empréstimos de equipamentos escolares
(notebooks/tablets), substituindo o controle manual em papel.

**Importante:** este projeto **não cria nem migra o banco de dados**. O arquivo
`data/lostemp.db` e suas tabelas são de responsabilidade do mantenedor do
schema (ver seção [Pendências de schema](#pendências-de-schema)). O código
apenas abre a conexão com `gorm.Open` e consome as tabelas existentes.

## Stack

- Go 1.22+
- Fiber v2 (HTTP)
- GORM + SQLite em arquivo local (`./data/lostemp.db`)
- `golang-jwt/jwt/v5` (autenticação)
- `golang.org/x/crypto` (bcrypt)
- `go-playground/validator` (validação)
- `gofiber/swagger` + `swaggo/swag` (documentação em `/docs`)

> **Decisão registrada:** o driver `gorm.io/driver/sqlite` (baseado em
> `mattn/go-sqlite3`) exige compilador C (CGo). Como o ambiente de
> desenvolvimento não possui `gcc`/CGo, foi usado o driver GORM puro-Go
> **`github.com/glebarez/sqlite`** (mesma API `sqlite.Open`, mesmo GORM,
> baseado em `modernc.org/sqlite`). Nenhuma outra parte da stack foi alterada.

## Requisitos

- Go 1.22+ (testado com Go 1.26)
- Nenhum compilador C é necessário

## Instalação, build e execução

```bash
cd backend
go mod tidy
go build ./...
go test ./...
```

O arquivo do banco precisa existir em `data/lostemp.db` com o schema correto
(ver [Pendências de schema](#pendências-de-schema)). Com ele presente:

```bash
go run cmd/server/main.go
# API em http://localhost:3000 — docs em http://localhost:3000/docs
```

### Seed (opcional — somente dados)

Insere **apenas linhas** (1 admin, 1 professor, 1 técnico e 5 equipamentos),
nunca tabelas. As tabelas precisam existir antes.

```bash
DB_PATH=./data/lostemp.db go run ./scripts/seed.go
```

Usuários criados (senha `123456` para todos, apenas desenvolvimento):

| Papel | E-mail | Matrícula |
|---|---|---|
| admin | admin@escola.edu | ADM-0001 |
| professor | professor@escola.edu | PROF-0001 |
| technician | ti@escola.edu | TEC-0001 |

### Variáveis de ambiente

| Variável | Default | Descrição |
|---|---|---|
| `PORT` | `3000` | Porta do servidor |
| `DB_PATH` | `./data/lostemp.db` | Caminho do banco SQLite existente |
| `JWT_SECRET` | `dev-secret-change-me` | Segredo de assinatura do JWT |
| `JWT_TTL_HOURS` | `12` | Expiração do token em horas |

## Papéis (RBAC)

- `admin` — Gestor/Admin: gerencia equipamentos, usuários, libera retiradas,
  devoluções e retiradas de emergência; vê toda a agenda.
- `professor` — consulta disponibilidade, reserva e vê os próprios empréstimos.
- `technician` — Técnico de TI: registra e conclui manutenções.

## Endpoints (prefixo `/api/v1`)

Autenticação: `Authorization: Bearer <token>`.

| Método | Rota | Papéis | Descrição |
|---|---|---|---|
| POST | `/auth/login` | público | Login → JWT |
| GET | `/auth/me` | autenticado | Usuário logado |
| POST | `/users` | admin | Criar usuário |
| GET | `/users` | admin | Listar (filtros `role`, `status`; paginado) |
| GET | `/users/:id` | admin | Buscar usuário |
| PATCH | `/users/:id` | admin | Atualizar usuário |
| DELETE | `/users/:id` | admin | Excluir usuário |
| POST | `/equipments` | admin | Cadastrar equipamento |
| GET | `/equipments` | autenticado | Listar (filtros `identifier`, `brand`, `model`, `status`; paginado) |
| GET | `/equipments/:id` | autenticado | Buscar equipamento |
| PATCH | `/equipments/:id` | admin | Atualizar equipamento |
| DELETE | `/equipments/:id` | admin | Excluir equipamento |
| POST | `/loans` | admin, professor | Registrar retirada |
| POST | `/loans/emergency` | admin | Retirada de emergência (sobrescreve agenda) |
| POST | `/loans/:id/return` | admin | Registrar devolução |
| GET | `/loans` | autenticado | Listar (admin vê todos; professor vê os próprios; paginado) |
| GET | `/loans/overdue` | autenticado | Listar empréstimos em atraso |
| POST | `/reservations` | admin, professor | Criar reserva |
| GET | `/reservations` | autenticado | Listar (agenda) |
| PATCH | `/reservations/:id/cancel` | admin, dono | Cancelar reserva |
| POST | `/maintenances` | technician, admin | Abrir manutenção |
| GET | `/maintenances` | technician, admin | Listar |
| GET | `/maintenances/:id` | technician, admin | Buscar |
| PATCH | `/maintenances/:id` | technician, admin | Atualizar (ex.: `status: FINISHED`) |
| DELETE | `/maintenances/:id` | technician, admin | Excluir |
| GET | `/health` | público | Health check |

Listagens são paginadas com `?page=1&page_size=20` (máx. 100) e respondem
`{ "data": [...], "meta": { "page", "page_size", "total" } }`.

## Regras de negócio implementadas

1. Equipamento só é retirado se estiver `AVAILABLE`.
2. Usuário com empréstimo vencido (`ACTIVE` com `due_at` no passado) fica
   bloqueado para novos empréstimos **e** reservas.
3. Retirada de emergência cancela reservas ativas do equipamento e registra
   auditoria (`reservation.cancel_emergency` + `loan.emergency`).
4. Devolução encerra o empréstimo (`RETURNED`) e devolve o equipamento para
   `AVAILABLE`.
5. Toda mutação relevante gera `AuditLog` (usuário, data/hora, entidade).
6. Reservas não podem conflitar no mesmo equipamento/horário (verificada com
   `starts_at < ? AND ends_at > ?`).
7. Manutenção em aberto (`OPEN`) coloca o equipamento como `MAINTENANCE`;
   concluir (`FINISHED`) o devolve para `AVAILABLE`.

## Testes

```bash
go test ./...          # inclui regra de bloqueio por atraso e conflito de agenda
```

## Swagger

Gerado a partir das anotações. Para regenerar após mudanças:

```bash
swag init -g cmd/server/main.go --output docs
```

## Decisões em aberto (assumidas da forma mais simples)

- **Status**: `active/inactive` (usuário), `AVAILABLE/LOANED/MAINTENANCE`
  (equipamento), `ACTIVE/RETURNED` (empréstimo),
  `ACTIVE/CANCELLED/FULFILLED` (reserva), `OPEN/FINISHED` (manutenção).
- **Titular do empréstimo/reserva**: `user_id` deve ser um usuário com papel
  `professor` e status `active`.
- **JWT**: claims `sub/role/email/name`, expiração padrão de 12h; nenhuma
  checagem de status do usuário por requisição (só no login).
- **Emergência**: cancela **todas** as reservas `ACTIVE` do equipamento.
- **Reserva** não valida o status do equipamento (`LOANED`/`MAINTENANCE`) —
  apenas conflito reserva-vs-reserva.
- **Auditoria**: `entity_id` armazenado como string (id numérico da entidade);
  payload em JSON.
- **`POST /users`** exige `password`; **`PATCH /users`** aceita campos
  parciais.

## Pendências de schema

Durante a implementação, verificou-se que o arquivo do banco e as tabelas
**não foram fornecidos**. Como este projeto não pode criar/migrar schema, segue
o **relatório do schema esperado** para que o mantenedor crie o arquivo
`data/lostemp.db` (as colunas abaixo são exatamente as que as structs GORM em
`internal/db/models.go` mapeiam):

**`users`**: `id` (PK), `name`, `email`, `password_hash`, `registration`,
`role`, `status`, `created_at`, `updated_at`

**`equipments`**: `id` (PK), `identifier` (único), `brand`, `model`, `status`,
`created_at`, `updated_at`

**`loans`**: `id` (PK), `equipment_id`, `user_id`, `withdrawn_by`,
`withdrawn_at`, `due_at`, `returned_at`, `returned_by`, `status`, `notes`,
`created_at`

**`reservations`**: `id` (PK), `equipment_id`, `user_id`, `starts_at`,
`ends_at`, `status`, `created_at`, `updated_at`

**`maintenances`**: `id` (PK), `equipment_id`, `technician_id`, `description`,
`status`, `started_at`, `finished_at`

**`audit_logs`**: `id` (PK), `user_id`, `action`, `entity_type`, `entity_id`,
`payload`, `created_at`

Recomendações (aplicadas pela aplicação, não pelo schema):

- Índice único em `users.email` e `equipments.identifier` — a aplicação já
  valida unicidade em nível de query, mas um índice único evita corrida.
- Índices em `loans(user_id)`, `loans(equipment_id)`, `loans(status)`,
  `reservations(equipment_id)`, `reservations(starts_at/ends_at)`.

> Nenhum arquivo de migration, `AutoMigrate` ou `CREATE TABLE` é gerado por
> este projeto. Se surgir a necessidade de nova coluna/tabela, deve ser
> **reportada** (não executada) aqui.
