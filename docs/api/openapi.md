# API — индекс

> **Единый контракт:** [API_Contract.md](./API_Contract.md)  
> **OpenAPI 3.1:** [../contracts/openapi.yaml](../contracts/openapi.yaml)

## Быстро

| | |
|--|--|
| Base URL | `http://localhost:8000/api/v1` |
| Auth | `Authorization: Bearer <jwt>` |
| Gateway | `back/services/core-api/api-gateway` |

## Critical для demo «Поток»

| Метод | Path | Назначение |
|-------|------|------------|
| POST | `/auth/login` | Вход |
| POST | `/expenses/manual` | Голос / ручной |
| GET | `/dashboard/sankey` | Главный экран |
| GET | `/credits/dashboard` | Кредитный светофор |
| POST | `/scenarios/simulate` | Ипотека / «что если» |
| POST | `/goals` | Цель |
| POST | `/fns/ticket` | ФНС (опционально) |

Полная таблица с приоритетами и JSON — в [API_Contract.md](./API_Contract.md).  
TypeScript (front): [typescript-types.md](./typescript-types.md).

## Связи

- [NAVI.md](../../NAVI.md)
- [../product/input-methods.md](../product/input-methods.md)
