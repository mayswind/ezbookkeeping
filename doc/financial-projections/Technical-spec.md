# Proyecciones — Especificación técnica de cambios

Documento breve de alcance técnico: qué se toca en backend, qué se toca en frontend.
Complementa [Proposal.md](./Proposal.md) (diseño funcional) y [Implementation-plan.md](./Implementation-plan.md) (orden de trabajo).

## Objetivo

Nueva pestaña **Proyecciones** dentro de la sección de navegación **Datos de Transacción**, que muestra en formato tabla la evolución de ingresos, egresos, neto y acumulado del período seleccionado, calculada a partir de las transacciones ya materializadas (pasado) y de las plantillas programadas `SCHEDULE` (futuro).

## Regla de negocio (resumen)

- **Corte en `now`.** Los meses/días anteriores o iguales a `now` se leen de transacciones **reales**; los posteriores se **simulan** con el motor de frecuencias de las plantillas.
- El corte hay que aplicarlo **explícitamente** sobre el lado real: `GetAccountsAndCategoriesMonthlyInflowAndOutflow` devuelve todo el rango de meses, incluidas transacciones manuales con fecha futura. Sin ese filtro los dos conjuntos se solapan y hay doble conteo.
- **Ambos lados se filtran igual**: se excluyen transferencias (`TRANSFER_OUT` / `TRANSFER_IN` en el lado real, plantillas con `Type = TRANSFER` en el simulado). De lo contrario los meses pasados muestran transferencias y los futuros no.
- **La moneda se preserva vía `AccountId`.** El `Amount` de una plantilla está en la moneda de su cuenta; agregar sin `AccountId` impide la conversión a moneda por defecto en el frontend.

---

# Cambios en Backend

## Archivos a modificar / crear

| Archivo | Ubicación actual | Acción |
|---|---|---|
| `transactions.go` (servicio) | [pkg/services/transactions.go](../../pkg/services/transactions.go) | Modificar: extraer el motor de frecuencias de `CreateScheduledTransactions` (líneas 845-890) |
| `scheduled_frequency.go` | `pkg/services/` | **Crear**: motor de frecuencias puro y reutilizable |
| `transaction_projections.go` (servicio) | `pkg/services/` | **Crear**: método `GetProjectedCategoryAmountsByMonth` sobre `TransactionService` |
| `transaction.go` (modelos) | [pkg/models/transaction.go](../../pkg/models/transaction.go) | Modificar: agregar `TransactionProjectionRequest` junto a `TransactionStatisticTrendsRequest` (línea 304) |
| `transaction_projections.go` (api) | `pkg/api/` | **Crear**: `TransactionProjectionsHandler`, método de `TransactionsApi` |
| `webserver.go` | [cmd/webserver.go](../../cmd/webserver.go) | Modificar: registrar la ruta después de la línea 411 (`asset_trends.json`) |
| `transaction_projections_test.go` | `pkg/services/` | **Crear**: tests de frecuencias y de merge |

> El servicio y el handler se implementan como **métodos de las structs existentes** (`TransactionService`, `TransactionsApi`) en archivos separados. Crear structs nuevas obligaría a cablear inyección de dependencias sin ninguna ganancia.

## Models que se usarán

### Existentes (sin cambios)

| Model | Ubicación | Uso |
|---|---|---|
| `TransactionTemplate` | [pkg/models/transaction_template.go:33](../../pkg/models/transaction_template.go#L33) | Fuente de las ocurrencias futuras: `Type`, `CategoryId`, `AccountId`, `Amount`, `ScheduledFrequencyType`, `ScheduledFrequency`, `ScheduledStartTime`, `ScheduledEndTime`, `ScheduledAt`, `ScheduledTimezoneUtcOffset` |
| `Transaction` | [pkg/models/transaction.go](../../pkg/models/transaction.go) | Lado real, ya materializado |
| `TransactionTotalAmount` | [pkg/models/transaction.go](../../pkg/models/transaction.go) | Estructura intermedia de agregación (`Type`, `CategoryId`, `AccountId`, `Amount`) |
| `TransactionStatisticResponseItem` | [pkg/models/transaction.go:836](../../pkg/models/transaction.go) | Item de respuesta: `CategoryId`, `AccountId`, `TotalAmount` |
| `TransactionStatisticTrendsResponseItem` | [pkg/models/transaction.go:475](../../pkg/models/transaction.go) | Contenedor mensual: `Year`, `Month`, `Items[]` |
| `YearMonthRangeRequest` | [pkg/models/transaction.go:391](../../pkg/models/transaction.go#L391) | Parseo de `start_year_month` / `end_year_month` |

**Se reutiliza el DTO de tendencias tal cual.** El frontend ya sabe consumir esa forma (`assembleAccountAndCategoryInfo` en [src/stores/statistics.ts:1030](../../src/stores/statistics.ts#L1030)), incluida la conversión de moneda por cuenta. No se crea un DTO paralelo.

> Nota: `TransactionStatisticResponseItem` **no tiene campo `Type`**. Ingreso vs. egreso se deriva del tipo de la *categoría* en el frontend (`CategoryType.Income` / `CategoryType.Expense`). El backend no lo envía y no hace falta que lo envíe.

### Nuevo

```go
// pkg/models/transaction.go
type TransactionProjectionRequest struct {
    YearMonthRangeRequest
    UseTransactionTimezone bool `form:"use_transaction_timezone"`
}
```

### Firma del motor de frecuencias

```go
// pkg/services/scheduled_frequency.go
// Devuelve todos los instantes en que la plantilla dispararía dentro de (from, to],
// ya resueltos en la zona horaria de la plantilla.
func GetScheduledOccurrences(template *models.TransactionTemplate, from, to time.Time) ([]time.Time, error)
```

No alcanza con un `MatchesScheduledFrequency(template, date)` de fecha suelta:

- El cron no evalúa una "fecha", sino el instante `inicio_del_día_UTC + ScheduledAt minutos` convertido a la zona de la plantilla ([transactions.go:856-858](../../pkg/services/transactions.go#L856-L858)). Una firma por fecha desalinea las plantillas cuyo `ScheduledAt` cruza el límite de día respecto de su zona.
- `EVERY_N_DAYS` necesita `ScheduledStartTime`; `MONTHLY` con día negativo necesita el mes evaluado.
- **Bug a corregir en la extracción**: hoy el día negativo se resuelve con `GetMaxDayOfMonth(currentTime...)`, donde `currentTime` es la hora *local del servidor* ([transactions.go:846](../../pkg/services/transactions.go#L846)). Al proyectar meses hacia adelante hay que calcularlo sobre el mes evaluado y en la zona de la plantilla.

Semántica que debe replicarse **exactamente** para que proyección y cron coincidan:

- `MONTHLY` día 31 en un mes corto → **no hay ocurrencia** (no se recorta al último día).
- Plantillas `Hidden` → **sí disparan** (el cron no las excluye).
- `ScheduledFrequencyType = DISABLED` o `ScheduledFrequency` vacío → se descartan.

## Endpoint

```
GET /api/v1/transactions/statistics/projections.json
    ?start_year_month=YYYY-MM
    &end_year_month=YYYY-MM
    &use_transaction_timezone=false
```

Se ubica bajo `/transactions/statistics/` por consistencia con `trends.json` y `asset_trends.json` ([cmd/webserver.go:409-411](../../cmd/webserver.go#L409-L411)).

**Respuesta:** `[]TransactionStatisticTrendsResponseItem` — idéntica a la de tendencias.

```json
[
  { "year": 2026, "month": 8,
    "items": [ { "categoryId": "123", "accountId": "456", "amount": "150000" } ] }
]
```

### Lógica de negocio — `GetProjectedCategoryAmountsByMonth(uid, startYM, endYM, clientTz, useTransactionTimezone)`

1. **Validar el rango**: `start <= end`, ambos presentes, y tope máximo de meses (p. ej. 60) para evitar rangos abusivos — el lado real carga en memoria todas las transacciones del período.
2. **Lado real** — llamar a `GetAccountsAndCategoriesMonthlyInflowAndOutflow` ([pkg/services/transactions.go:2551](../../pkg/services/transactions.go#L2551)) con el rango pedido, y sobre el resultado:
   - descartar `TRANSACTION_DB_TYPE_TRANSFER_OUT` y `TRANSACTION_DB_TYPE_TRANSFER_IN`;
   - **descartar lo que tenga `transaction_time > now`** (transacciones manuales con fecha futura), para que no se solape con lo simulado.
3. **Lado simulado** — traer las plantillas del usuario con `TemplateType = SCHEDULE`, `Deleted = false` y `ScheduledFrequencyType != DISABLED`; excluir las de `Type = TRANSFER`.
4. Para cada plantilla, `GetScheduledOccurrences(template, now, finDelRango)`, respetando `ScheduledStartTime` / `ScheduledEndTime`. Cada ocurrencia aporta `template.Amount` a la clave `(Year, Month, CategoryId, AccountId)`.
5. **Bucketear por mes con la misma regla de zona horaria que el lado real**: `clientTimezone`, o la zona de la transacción si `use_transaction_timezone=true`. Si el mes se calcula con la zona de la plantilla y los reales con la del cliente, la misma ocurrencia puede caer en meses distintos.
6. **Fusionar** ambos conjuntos sumando por `(Year, Month, CategoryId, AccountId)`.
7. Ordenar por año-mes y serializar a `TransactionStatisticTrendsResponseItem`.

### Limitaciones conocidas (v1)

- Si el servidor estuvo caído, una ocurrencia pasada no se materializó y tampoco se proyecta (ya es pasado) → subconteo silencioso.
- No se aplican filtros de tag / keyword (a diferencia de `trends.json`); el lado real va sin filtrar.

---

# Cambios en Frontend

## Ventanas a implementar

### 1. Página de Proyecciones (escritorio) — `/projections/transaction`

Estructura:

- **Encabezado**: selector de período (mes inicio – mes fin), reutilizando [MonthRangeSelectionDialog.vue](../../src/components/desktop/MonthRangeSelectionDialog.vue).
- **Cuerpo**: la tabla de proyecciones (ver abajo).

### 2. Tabla de proyecciones

Una columna por mes del período, más una columna **Total** al final.

```
                              | Mes 1 | Mes 2 | Mes 3 | Mes 4 | Total |
▼ INGRESOS
  ▼ Trabajo
      Salario                 |  1000 |  1000 |  1000 |  1000 |  4000 |
      Bonos                   |   500 |   500 |   500 |   500 |  2000 |
    Subtotal Trabajo          |  1500 |  1500 |  1500 |  1500 |  6000 |
  Total Ingresos              |  1500 |  1500 |  1500 |  1500 |  6000 |
▼ EGRESOS
  ▼ Entretenimiento
      Suscripciones           |   200 |   200 |   200 |   200 |   800 |
    Subtotal Entretenimiento  |   200 |   200 |   200 |   200 |   800 |
  ▼ Trabajo
      Suscripciones           |   200 |   200 |   200 |   200 |   800 |
    Subtotal Trabajo          |   200 |   200 |   200 |   200 |   800 |
  ▼ Casa
      Recibos                 |   600 |   600 |   600 |   600 |  2400 |
    Subtotal Casa             |   600 |   600 |   600 |   600 |  2400 |
  Total Egresos               |  1000 |  1000 |  1000 |  1000 |  4000 |
─────────────────────────────────────────────────────────────────────
  Neto del mes                |   500 |   500 |   500 |   500 |  2000 |
  Acumulado                   |   500 |  1000 |  1500 |  2000 |  2000 |
```

**Jerarquía de filas (3 niveles + 2 filas de cierre):**

| Nivel | Fila | Colapsable | Al colapsar |
|---|---|---|---|
| 0 | Sección `INGRESOS` / `EGRESOS` | Sí | Se ocultan sus categorías y subtotales; **la fila `Total Ingresos` / `Total Egresos` permanece visible** |
| 1 | Categoría (*Trabajo*, *Casa*…) | Sí | Se ocultan sus subcategorías; **la fila `Subtotal <categoría>` permanece visible** |
| 2 | Subcategoría (*Salario*, *Recibos*…) | No | — |
| — | `Neto del mes` = Total Ingresos − Total Egresos | No | Siempre visible |
| — | `Acumulado` = suma corrida del Neto | No | Siempre visible |

**Celdas de la columna Total:**

- Filas de categoría / subcategoría / totales / neto → suma horizontal de los meses.
- Fila `Acumulado` → el valor del **último mes** del período (no la suma), porque ya es un acumulado.

**Acordeón:** `v-expansion-panels multiple` de Vuetify, mismo patrón que [PresetDialog.vue](../../src/views/desktop/categories/list/dialogs/PresetDialog.vue). Cada nivel se controla de forma independiente; el estado de expansión vive en el componente, no en el store.

**Formateo de importes:** montos convertidos a la moneda por defecto del usuario vía `exchangeRatesStore.getExchangedAmount(...)` a partir de la moneda de `item.account`, replicando [statistics.ts:1074-1080](../../src/stores/statistics.ts#L1074-L1080). Los ítems sin cuenta conocida quedan sin convertir y deben señalarse en la UI.

## Archivos a modificar / crear

| Archivo | Ubicación | Acción |
|---|---|---|
| `transaction.ts` (modelos) | [src/models/transaction.ts](../../src/models/transaction.ts) | Modificar: agregar `TransactionProjectionRequest` junto a `TransactionStatisticTrendsRequest` (línea 693). Los tipos de respuesta se reutilizan tal cual |
| `services.ts` | [src/lib/services.ts](../../src/lib/services.ts) | Modificar: agregar `getTransactionProjections(req)`, mismo esqueleto que `getTransactionStatisticsTrends` (líneas 573-597) |
| `projection.ts` (store) | `src/stores/` | **Crear**: acción de carga + computeds que arman las filas de la tabla |
| `ProjectionTable.vue` | `src/components/desktop/` | **Crear**: tabla + acordeón de 3 niveles |
| `TransactionPage.vue` | `src/views/desktop/projections/` | **Crear**: página (selector de período + tabla) |
| `desktop.ts` (router) | [src/router/desktop.ts](../../src/router/desktop.ts) | Modificar: import de la página + ruta `/projections/transaction` |
| `MainLayout.vue` | [src/views/desktop/MainLayout.vue](../../src/views/desktop/MainLayout.vue) | Modificar: `<li class="nav-link">` en la sección "Transaction Data", después de "Insights Explorer" (líneas 43-48) |
| `*.json` (i18n) | [src/locales/](../../src/locales/) | Modificar: **20 archivos** de idioma |

### Store `src/stores/projection.ts`

Patrón de [src/stores/statistics.ts](../../src/stores/statistics.ts):

- **Estado**: rango de meses seleccionado, datos crudos de la API, flag de carga.
- **Acción** `loadProjections({ force })`, análoga a `loadTrendAnalysis` ([statistics.ts:1899](../../src/stores/statistics.ts#L1899)).
- **Enriquecimiento**: resolver `category` / `primaryCategory` / `account` y convertir moneda, reutilizando la lógica de `assembleAccountAndCategoryInfo`.
- **Computeds**: árbol sección → categoría → subcategoría con importes por mes, subtotales por categoría, totales por sección, fila Neto y fila Acumulado.
- **Invalidación de caché**: el resultado depende de `now` **y** del conjunto de plantillas. Hay que invalidar al crear/editar/borrar una plantilla programada y al cambiar de día — no alcanza con cachear por período.

### i18n

Claves nuevas a agregar en los 20 archivos de `src/locales/`:

- `"Projections"`, `"Net"`, `"Accumulated"`, `"Subtotal"`, `"Projected"`

`"Total Income"` y `"Total Expense"` **ya existen** ([en.json:2307-2308](../../src/locales/en.json#L2307-L2308)) — reutilizarlas, no duplicarlas.

## Fuera de alcance (v1)

- **Mobile.** No se agrega página en `src/views/mobile/`. Insights Explorer tampoco la tiene, así que no rompe la paridad existente del proyecto.
- **Transferencias.** Excluidas de ambos lados del cálculo.
- **Saldo inicial en el acumulado.** `Neto` y `Acumulado` reflejan solo el *flujo* del período, no el patrimonio absoluto. Sumar el saldo actual (`getNetAssets()`, [src/stores/account.ts:472](../../src/stores/account.ts#L472)) queda para v2.
