---
title: "WTF is Nimbus"
aliases:
  - Nimbus project plan
  - Что такое Nimbus
tags:
  - nimbus
  - architecture
  - roadmap
  - librus
status: draft
updated: 2026-09-06
---

# WTF is Nimbus?

> [!abstract] Коротко
> **Nimbus** — self-hosted open-source веб-приложение для учеников и родителей, которое превращает неудобный и нестабильный интерфейс Librus Synergia в нормальный, быстрый и последовательно спроектированный продукт.
>
> Go backend изолирует reverse-engineered Librus API, управляет авторизацией, нормализует данные и кэширует их в PostgreSQL. Nuxt/Vue frontend с SSR предоставляет удобный интерфейс электронного журнала.

## Зачем существует проект

Librus решает важную задачу, но пользоваться им неудобно. Nimbus создаётся сначала как личный инструмент, который должен сделать ежедневную работу со школьным журналом быстрее и приятнее. Одновременно это open-source и portfolio project: архитектура, безопасность, документация и качество кода должны выглядеть как у небольшого production-сервиса, а не как у одноразового скрипта.

Nimbus не должен быть простым UI proxy поверх Librus. Его задача — создать устойчивую границу между нестабильным upstream API и собственным публичным API, пригодным для frontend и будущих клиентов.

## Product principles

1. **User experience first.** Пользователь взаимодействует с Nimbus, а не с особенностями Librus.
2. **Upstream isolation.** Форматы, ошибки и странности Librus не должны протекать в public API Nimbus.
3. **Secure by default.** Пароли, cookies и session material считаются секретами.
4. **Fresh where it matters.** Изменения расписания и другие срочные данные обновляются чаще долговременных справочников.
5. **Graceful degradation.** Если Librus временно недоступен, Nimbus по возможности показывает последние сохранённые данные и честно сообщает время их обновления.
6. **Small now, scalable later.** Ожидается до 20 пользователей, но нельзя строить архитектуру, в которой данные и сессии разных пользователей смешиваются.
7. **Document discoveries.** Любая найденная особенность Librus API должна попасть в документацию или тест, а не оставаться только в голове разработчика.

## Главная цель: learning-first development

Nimbus создаётся не ради того, чтобы AI как можно быстрее сгенерировал готовое приложение. Главная цель проекта — чтобы автор на практике изучил:

- Go и устройство backend applications;
- HTTP, REST, cookies, sessions и authentication;
- PostgreSQL, SQL, migrations и data modelling;
- Vue, Nuxt и server-side rendering;
- Tailwind CSS и построение UI;
- Docker, VPS deployment, Git и нормальный engineering workflow;
- testing, security, observability и работу с нестабильным external API.

AI здесь — **наставник, pair programmer и reviewer**, но не автономный автор проекта. После каждого этапа автор должен понимать, что было сделано, почему выбран именно такой подход и как самостоятельно изменить или отладить получившийся код.

### Правила совместной работы с AI

1. Сначала объяснить понятия, цель шага и ожидаемый результат понятным языком.
2. Разбивать работу на маленькие проверяемые задачи, которые автор может выполнить самостоятельно.
3. Для нового учебного кода сначала предложить автору написать его самому; затем дать hints, разобрать ошибку и провести code review.
4. Не выдавать целый feature implementation без явной просьбы. Допустимы маленькие примеры, pseudocode и focused patches.
5. При code review не только исправлять, но и объяснять причину, trade-offs и способ проверить исправление.
6. Задавать короткие контрольные вопросы, чтобы проверить понимание, но не превращать работу в экзамен.
7. Отмечать, где используется временное упрощение для MVP, а где формируется production habit.
8. Не вводить библиотеку, pattern или infrastructure component без объяснения, какую конкретную проблему он решает.
9. Команды Git, Docker, Go, Nuxt и database migration по возможности выполняет автор под руководством AI.
10. Самостоятельно AI может делать уборку, механические изменения и явно порученные задачи, но после этого обязан кратко разобрать изменения с автором.

> [!important]
> Скорость разработки вторична. Если выбор стоит между «быстрее получить feature» и «понять технологию, написав её осмысленно», Nimbus выбирает второй вариант.

## Scope

### Пользователи

- ученик;
- родитель.

На первом этапе Nimbus не пытается связывать аккаунты ребёнка и родителя собственной family model. Каждый Librus account рассматривается как самостоятельный Nimbus account.

### MVP

- [ ] Вход через существующий аккаунт Librus Synergia.
- [ ] Безопасная собственная сессия Nimbus.
- [ ] План занятий / timetable с заменами и отменами.
- [ ] Просмотр оценок.
- [ ] Просмотр и отправка сообщений.
- [ ] Базовый SSR-интерфейс, удобный на desktop и mobile.
- [ ] Кэширование данных и отображение `last updated`.
- [ ] Docker deployment на VPS с PostgreSQL.

### После MVP

- посещаемость;
- домашние задания;
- контрольные и события;
- профиль, школа, учителя и предметы;
- замечания и оценки за поведение;
- push/email/другие уведомления;
- симулятор среднего балла и аналитика;
- полноценное покрытие функций Librus;
- дополнительные функции, список которых будет определён позже.

### Пока не входит в scope

- аккаунты учителей и школьных администраторов;
- собственная связь parent ↔ child;
- GraphQL;
- WebSocket/SSE без конкретного realtime use case;
- native mobile application;
- сохранение пароля Librus по умолчанию.

## High-level architecture

```text
Browser
  │  Secure HttpOnly session cookie
  ▼
Nuxt / Vue (SSR)
  │  Nimbus REST API
  ▼
Go backend
  ├── Auth & session service
  ├── Domain services
  ├── Librus adapter / anti-corruption layer
  ├── Cache & synchronization policy
  └── PostgreSQL repositories
          │
          ├──────── PostgreSQL
          │
          └──────── Librus Synergia (unofficial upstream API)
```

### Backend

- **Language:** Go.
- **HTTP API:** versioned REST, starting with `/api/v1`.
- **Database:** PostgreSQL.
- **Deployment:** Docker Compose initially; application and database run as separate containers.
- **Librus integration:** isolated adapter/client package. Domain and HTTP layers must not know about raw Librus JSON structures.

Suggested package boundaries (not a rigid final tree):

```text
backend/
  cmd/api/                 # entrypoint
  internal/
    auth/                  # Nimbus sessions and Librus login orchestration
    librus/                # transport, auth handshake, DTOs, endpoint adapters
    grades/                # grade domain/service/repository
    timetable/             # timetable domain/service/repository
    messages/              # message domain/service/repository
    user/                  # Nimbus user/account
    platform/              # config, database, logging, encryption, clock
    httpapi/               # Gin routes, middleware, handlers, response DTOs
  migrations/
```

Gin уже используется и может остаться для MVP. Бизнес-логика не должна зависеть от Gin, чтобы её можно было тестировать обычными unit tests.

### Frontend

- **Framework:** Nuxt + Vue.
- **Rendering:** SSR by default.
- Frontend обращается только к Nimbus API и никогда напрямую к Librus.
- DTO публичного API принадлежат Nimbus и не повторяют upstream format автоматически.
- UI должен показывать loading/error/stale states и время последней синхронизации.

## Authentication and sessions

### Решение для Nimbus

Для browser client основной вариант — **opaque server-side session**:

- браузер получает случайный session ID в cookie;
- cookie имеет `HttpOnly`, `Secure` и подходящий `SameSite`;
- сервер хранит сессию и её владельца в PostgreSQL (или позже в Redis, только если появится необходимость);
- CSRF protection обязателен для state-changing endpoints;
- session rotation выполняется после login и других security-sensitive событий;
- logout отзывает Nimbus session и по возможности завершает upstream session.

JWT в `localStorage` не нужен: для одного first-party web application он ухудшает отзыв сессий и увеличивает последствия XSS, не давая Nimbus заметной пользы.

### Работа с Librus credentials

Пароль используется только во время login handshake и **не сохраняется по умолчанию**. В логах, ошибках, traces и database dumps не должно быть ни пароля, ни полного набора Librus cookies.

После успешного входа Nimbus может хранить необходимый upstream session material в зашифрованном виде, отдельно для каждого пользователя. Encryption key передаётся через environment/secret storage и не лежит в БД или Git.

> [!question] Open decision
> Проверить экспериментально срок жизни и refresh flow Librus session. Если долгоживущую сессию безопасно восстановить невозможно, выбрать между повторным вводом пароля и opt-in encrypted credential storage. До отдельного решения пароль не хранить.

### Критический invariant

У каждого пользователя должен быть собственный Librus HTTP client/cookie jar или сериализованное session state. Никогда не использовать один глобальный cookie jar для всех запросов сервера.

## Data and caching strategy

PostgreSQL используется не только как технический cache, но и как источник последнего известного состояния. Каждая синхронизируемая запись или snapshot должна позволять определить:

- кому принадлежат данные (`user_id` / upstream account);
- когда данные получены (`fetched_at`);
- до какого момента они считаются свежими (`fresh_until`);
- источник и версию mapper-а, если это помогает миграциям;
- удалена ли сущность upstream или просто отсутствует в неполном ответе.

### Предварительная freshness policy

| Data type | Initial policy | Поведение |
|---|---:|---|
| Текущий/следующий timetable | 1–5 min | stale-while-revalidate; manual refresh; чаще в учебные часы |
| Далёкие недели timetable | 15–60 min | background refresh |
| Grades | 5–15 min | refresh при открытии, затем cache |
| Message list | 1–5 min | refresh при открытии |
| Message body | долго / immutable-ish | обновлять по необходимости |
| Dictionaries: subjects, teachers, classrooms | 12–24 h | редкое фоновое обновление |
| Profile/school | 24 h | обновление при login и периодически |

Значения являются starting point, а не контрактом. Их нужно сделать конфигурируемыми и уточнить по реальному поведению Librus.

При upstream failure API может вернуть stale data вместе с metadata (`stale: true`, `fetchedAt`) вместо полного отказа. Для отправки сообщений cache fallback невозможен: успех возвращается только после подтверждённого upstream write. Повтор write-запроса требует idempotency strategy, чтобы не отправить сообщение дважды.

## Public API guidelines

- Base path: `/api/v1`.
- JSON field naming — единообразный `camelCase`.
- Время — ISO 8601/RFC 3339 с timezone; школьные даты без времени моделировать отдельно.
- Pagination, filtering и sorting должны иметь общий формат.
- Ошибки возвращаются в стабильном Nimbus format, например `code`, `message`, `details`, `requestId`.
- Не отдавать пользователю raw cookies, upstream HTML или внутренние Librus errors.
- Не использовать HTTP 401 для любого сбоя upstream: различать invalid credentials, expired upstream session, rate limiting, validation error и temporary outage.
- Write endpoints должны иметь CSRF protection и при необходимости idempotency key.
- OpenAPI specification должна стать частью backend и обновляться вместе с handlers.

Черновые ресурсы MVP:

```text
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
GET    /api/v1/auth/session

GET    /api/v1/timetable?weekStart=YYYY-MM-DD
POST   /api/v1/timetable/refresh

GET    /api/v1/grades

GET    /api/v1/messages
GET    /api/v1/messages/{id}
POST   /api/v1/messages
```

Конкретные request/response schemas определяются после фиксации реальных Librus payloads и создания Nimbus domain models.

## Librus integration

Librus Synergia рассматривается как unstable external dependency без гарантированного публичного контракта.

Integration layer должен поддерживать:

- единый HTTP transport с timeout-ами;
- context cancellation;
- ограниченные retries только для безопасных/idempotent операций;
- per-user cookie/session state;
- rate limiting и backoff;
- typed upstream errors;
- redacted structured logging;
- DTO для Librus отдельно от domain models Nimbus;
- contract fixtures/tests на обезличенных ответах;
- обнаружение неожиданных изменений schema;
- endpoint-specific cache policy.

### Research reference: Librusek

Основной reference project: [SimonB50/librusek](https://github.com/SimonB50/librusek).

На момент написания у него полезны:

- `lib/auth.js` — последовательность Synergia authentication/activation и refresh flow;
- `lib/request.js` — selectors, pagination и cache semantics;
- `lib/timetable.js` — `Timetables`, `TimetableEntries`, homework/exam-related resources и teacher free days;
- `lib/grades.js` — ресурсы и преобразования оценок;
- `lib/messages.js` — работа с сообщениями;
- остальные модули `lib/` — user, school, class, lessons, notifications и связанные справочники.

Использовать Librusek как карту для исследования: сверять endpoints, HTTP methods, selectors, pagination и shape ответов с реальным аккаунтом. Затем реализовывать собственные Go adapters и собственные Nimbus domain models.

> [!warning] License boundary
> Librusek опубликован под **AGPL-3.0**. Изучение поведения API и независимая реализация протокола — не то же самое, что копирование кода. Если код или существенные части реализации переносятся/адаптируются напрямую, необходимо выполнить требования лицензии и сохранить корректную attribution. До публикации Nimbus отдельно выбрать лицензию проекта и проверить совместимость. Это project note, не юридическая консультация.

### Reverse-engineering log

Для каждого исследованного endpoint фиксировать:

```text
Purpose:
Method + URL:
Required cookies/headers:
Query/path parameters:
Example response fixture (redacted):
Pagination/selectors:
Observed errors:
Cache/freshness decision:
Nimbus domain mapping:
Last verified:
```

Никогда не коммитить реальные credentials, session cookies, имена учеников, сообщения, оценки или другие персональные данные. Fixtures должны быть синтетическими либо тщательно обезличенными.

## Security and privacy baseline

Nimbus будет обрабатывать школьные данные и сообщения, поэтому даже portfolio project требует серьёзной security hygiene:

- HTTPS на VPS;
- secrets только через environment/secret files вне Git;
- шифрование Librus session material at rest;
- password/session/cookie redaction во всех логах;
- минимальный срок хранения персональных данных;
- возможность удалить аккаунт и связанные cached data;
- authorization check на каждом user-owned resource;
- CSRF, secure cookies, CORS allowlist, request size limits;
- rate limits на login и message sending;
- database migrations и регулярные backups;
- dependency/container scanning;
- никакой production telemetry с содержимым оценок или сообщений;
- threat model перед публичным deployment.

Отдельно проверить условия использования Librus, применимое законодательство и обязанности оператора сервиса до предоставления публичного доступа другим людям.

## Observability and operations

- structured logs с `requestId`, но без PII/secrets;
- health endpoints: liveness отдельно от readiness;
- metrics: request latency, upstream latency/error rate, cache hit ratio, sync age;
- graceful shutdown;
- database connection pool limits;
- automated migrations с контролируемым rollout;
- backups и проверка восстановления;
- Docker healthchecks;
- dev/staging/prod configuration без hardcoded URLs и keys.

## Testing strategy

- **Unit tests:** mappers, cache policy, auth/session rules, domain services.
- **Contract tests:** сохранённые обезличенные Librus fixtures → ожидаемые Nimbus models.
- **HTTP integration tests:** handlers + test database + fake Librus server.
- **End-to-end tests:** login/session, timetable, grades, messages через Nimbus UI.
- **Live smoke tests:** только вручную/в закрытом окружении с test account; не в публичном CI.

Код Librus adapter должен тестироваться через injectable transport/fake server. Unit tests не должны ходить в реальный Librus.

## Roadmap

### Phase 0 — привести фундамент в порядок

- [ ] Зафиксировать текущий authentication handshake тестами.
- [ ] Удалить глобальный shared `LibrusClient` из request path.
- [ ] Перестать возвращать Librus cookies из `/login`.
- [ ] Добавить config, structured errors и безопасное логирование.
- [ ] Поднять PostgreSQL и migrations.
- [ ] Спроектировать `users`, `sessions`, `librus_accounts` и encrypted upstream sessions.

### Phase 1 — auth vertical slice

- [ ] Реализовать login → Nimbus session cookie → authenticated `/session` → logout.
- [ ] Обработать invalid credentials, expired account, CAPTCHA и upstream outage отдельно.
- [ ] Проверить refresh/expiration Librus session.
- [ ] Добавить rate limit и CSRF protection.

### Phase 2 — timetable vertical slice

- [ ] Исследовать timetable endpoints и fixtures.
- [ ] Создать domain model и mapper.
- [ ] Реализовать cache + stale-while-revalidate.
- [ ] Сделать SSR timetable page с `last updated` и manual refresh.

### Phase 3 — grades

- [ ] Исследовать обычную и points-based grading models.
- [ ] Реализовать subject/grade dictionaries и normalization.
- [ ] Добавить SSR grades page.

### Phase 4 — messages

- [ ] Исследовать list/read/send flows, recipients и attachments.
- [ ] Реализовать безопасное чтение сообщений.
- [ ] Реализовать отправку с validation и защитой от duplicate send.

### Phase 5 — deployment quality

- [ ] Dockerfiles и Docker Compose.
- [ ] Reverse proxy + TLS.
- [ ] CI: format, lint, tests, build, security checks.
- [ ] Backups, metrics и минимальный runbook.
- [ ] README, screenshots, architecture overview и contribution guide для portfolio/open-source вида.

## Current repository audit

Состояние на 2026-09-06:

- backend содержит Go/Gin entrypoint, router, login handler и Librus client;
- authentication handshake уже работает как proof of concept;
- frontend пока отсутствует;
- PostgreSQL/config/migrations/tests отсутствуют;
- `LoginRequest` продублирован в `api/router.go` и `api/handlers/auth.go`;
- один `LibrusClient` создаётся в `main` и передаётся всем handlers;
- внутри него один cookie jar — это смешает сессии разных пользователей и является blocker для multi-user deployment;
- login response раскрывает Librus cookies клиенту;
- часть ошибок `http.NewRequest` игнорируется;
- handshake, transport, session lifecycle и domain API находятся в одном client file;
- `math/rand` допустим только если `X-Baner` действительно не является security token; это следует задокументировать тестом протокола;
- timeouts есть, но пока нет context propagation, retry/backoff, rate limiting и observability.

Текущий код не нужно бездумно расширять endpoint за endpoint-ом. Сначала следует сохранить работающий handshake тестами, затем выделить per-user session model и устойчивый Librus adapter.

## AI cheat sheet

> [!important] Инструкция для AI coding agents
> Этот раздел является рабочим контекстом. Перед изменениями в Nimbus сначала прочитай весь документ и актуальный код. Не считай roadmap уже реализованным.

### Mission

Строить Nimbus как production-shaped open-source portfolio project: Go REST backend + PostgreSQL + Nuxt/Vue SSR frontend, проксирующий и нормализующий неофициальный Librus Synergia API для учеников и родителей.

### Неподвижные решения

- Project name: **Nimbus**.
- Backend: **Go**.
- Frontend: **Nuxt + Vue, SSR**.
- Database: **PostgreSQL**.
- Deployment target: **Docker on VPS**.
- API style: **REST `/api/v1`** до появления доказанного use case для другого транспорта.
- Browser auth: server-side opaque session + secure `HttpOnly` cookie.
- Librus sessions строго per-user; глобальный shared cookie jar запрещён.
- Frontend не обращается к Librus напрямую.
- Librus DTO ≠ Nimbus domain model ≠ public API DTO.
- Не хранить Librus password без отдельного явно принятого решения.
- Schedule freshness важнее большинства других данных.

### Workflow для каждой новой Librus feature

1. Найти релевантную реализацию/endpoint в Librusek и других источниках.
2. Проверить наблюдение реальным controlled request; не доверять reference code как официальной спецификации.
3. Записать endpoint в reverse-engineering log и сохранить redacted fixture.
4. Сначала определить Nimbus domain model и API contract.
5. Реализовать upstream DTO + mapper + adapter.
6. Назначить cache/freshness policy.
7. Добавить unit и contract tests.
8. Только затем подключать handler и frontend.

### Engineering rules

- Не переписывать рабочий auth handshake без regression coverage.
- Не возвращать raw upstream response, cookies или secrets наружу.
- Не логировать credentials, cookie headers, message bodies и grade data.
- Любой запрос пользователя фильтруется по authenticated `user_id`.
- Не делать retries для non-idempotent writes без idempotency design.
- Все outbound requests получают timeout и `context.Context`.
- Ошибки Librus преобразуются в typed internal errors и стабильные API error codes.
- При сомнении между быстрым endpoint wrapper и нормальной boundary abstraction выбрать boundary abstraction.
- Не добавлять Redis, queues, WebSocket, microservices или Kubernetes без измеримой необходимости.
- Не копировать код Librusek вслепую; учитывать AGPL-3.0 и attribution.
- После изменения запускать formatter, relevant tests и build; честно сообщать, что не было проверено.
- Обновлять этот документ, когда принято архитектурное решение или обнаружено важное поведение Librus.
- Соблюдать learning-first workflow: не забирать у автора учебную часть работы и не генерировать большие реализации без прямой просьбы.
- Перед новой задачей кратко объяснить, чему она научит и какое место занимает в общей архитектуре.
- После порученной AI правки перечислить существенные изменения так, чтобы автор мог их пересказать и воспроизвести.

### Первый следующий шаг

Не начинать с frontend. Первый meaningful milestone: безопасный multi-user auth vertical slice с PostgreSQL-backed Nimbus sessions и изолированным per-user Librus session state. После него — timetable vertical slice от Librus adapter до SSR page.

## Open questions

- [ ] Какую лицензию выбрать для Nimbus с учётом предполагаемого использования материалов Librusek?
- [ ] Как долго живёт Librus session и можно ли надёжно обновлять её без пароля?
- [ ] Требуется ли opt-in encrypted storage пароля для полностью фоновой синхронизации?
- [ ] Какие точные payloads и edge cases есть у отправки сообщений?
- [ ] Нужны ли attachments в MVP messages?
- [ ] Какой дизайн/API response format использовать для stale data?
- [ ] Какая retention policy нужна для сообщений, оценок и логов?
- [ ] Какой полный post-MVP feature list нужен Nimbus?

## References

- [Librusek repository](https://github.com/SimonB50/librusek) — основной research reference, AGPL-3.0.
- [Librusek authentication module](https://github.com/SimonB50/librusek/blob/main/lib/auth.js).
- [Librusek request/cache module](https://github.com/SimonB50/librusek/blob/main/lib/request.js).
- [Librusek timetable module](https://github.com/SimonB50/librusek/blob/main/lib/timetable.js).
- [Librusek grades module](https://github.com/SimonB50/librusek/blob/main/lib/grades.js).
- [Librusek messages module](https://github.com/SimonB50/librusek/blob/main/lib/messages.js).
