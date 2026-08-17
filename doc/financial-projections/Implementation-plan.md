# Proyecciones — Plan de implementación

Plan paso a paso derivado de [Technical-spec.md](./Technical-spec.md). El orden es de abajo hacia arriba: primero el motor de frecuencias (con su red de tests), después servicio → API → ruta, y luego frontend (datos → UI → navegación → i18n).

**Regla de trabajo:** cada fase termina con un criterio de aceptación verificable. No avanzar a la siguiente sin cumplirlo.

Comandos de verificación:

```bash
go build ./...                 # compilación backend
go test ./pkg/...              # tests backend
npm run lint                   # vue-tsc + eslint
npm run test                   # vitest
npm run build                  # build frontend
```

> **No hay toolchain de Go instalado en el entorno de desarrollo actual.** Los comandos de backend se corren en contenedor:
>
> ```bash
> docker volume create ezbk-gocache
> docker run --rm -v "$PWD":/src -w /src -v ezbk-gocache:/gomod \
>   -e GOFLAGS="-mod=mod -buildvcs=false" \
>   -e GOCACHE=/gomod/build -e GOMODCACHE=/gomod/mod \
>   golang:1.26 sh -c "go build ./... && go test ./pkg/..."
> ```
>
> `-buildvcs=false` es necesario porque el contenedor no puede leer el estado de git del host. El volumen `ezbk-gocache` evita volver a descargar los módulos en cada corrida.
>
> Fallo preexistente y ajeno a este trabajo: `pkg/exchangerates` →
> `TestExchangeRatesApiLatestExchangeRateHandler_NationalBankOfUkraineDataSource` hace una petición HTTP real a una API de terceros y falla sin acceso a esa red. El resto de la suite pasa.

---

## Fase 0 — Red de seguridad sobre el cron (antes de tocar nada) ✅ COMPLETADA

El motor de frecuencias se va a extraer de `CreateScheduledTransactions`, que hoy **no tiene ningún test**. Sin caracterizar el comportamiento actual primero, no hay forma de saber si el refactor rompió el cron.

**Resultado:** [pkg/services/transactions_scheduled_test.go](../../pkg/services/transactions_scheduled_test.go) — 24 tests, 59 aserciones, en verde.

> **Limitación asumida:** `CreateScheduledTransactions` no es invocable desde un test — consulta las bases de datos de usuario y necesita un `core.Context`, y el paquete `pkg/services` no tiene harness de base de datos (los tests existentes solo cubren helpers puros). Por eso la caracterización usa `legacyShouldCreateScheduledTransaction`, una **transcripción literal** del bloque de decisión de [transactions.go:820-902](../../pkg/services/transactions.go#L820-L902) sin logging, contadores ni creación de transacciones. El valor está en la tabla de expectativas: la Fase 1 la repunta a `GetScheduledOccurrences` y cualquier divergencia de comportamiento aparece ahí. Queda fuera el prefiltro SQL de [transactions.go:785-798](../../pkg/services/transactions.go#L785-L798), que decide qué plantillas llegan a esta lógica.

- [x] **0.1** Crear `pkg/services/transactions_scheduled_test.go` con tests de caracterización del comportamiento **actual** de la lógica de frecuencias ([pkg/services/transactions.go:845-890](../../pkg/services/transactions.go#L845-L890)), usando `testify` como el resto de los tests del paquete.

  Casos mínimos, uno por tipo de frecuencia:
  - `WEEKLY` con uno y con varios días de la semana.
  - `MONTHLY` con día normal (15).
  - `MONTHLY` con día 31 en un mes de 30 días y en febrero → **debe no disparar** (hoy no hay recorte al último día del mes).
  - `MONTHLY` con día negativo (-1 = último día del mes) en meses de 28/29/30/31 días.
  - `DAILY`.
  - `YEARLY` con clave `mes*100+día`, incluyendo 29 de febrero en año bisiesto y no bisiesto.
  - `EVERY_N_DAYS` con `ScheduledStartTime` en el pasado, justo en el borde, y en el futuro.
  - Respeto de `ScheduledStartTime` / `ScheduledEndTime`.
  - Plantillas con `ScheduledFrequencyType = DISABLED` o `ScheduledFrequency = ""` → descartadas.
  - Plantilla con `Hidden = true` → **sí dispara** (el cron no filtra por `Hidden`).

- [x] **0.2** Documentar en el propio test, como comentario, las dos rarezas que deben preservarse: no hay recorte de día 31 en meses cortos, y las plantillas ocultas disparan igual.

**Criterio de aceptación:** ✅ `go test ./pkg/services/...` pasa en verde contra el código sin modificar.

### Hallazgos de la Fase 0 que condicionan la Fase 1

Además de los casos previstos, la caracterización dejó fijados tres comportamientos que el motor extraído debe respetar o corregir explícitamente:

1. **Bug del día negativo, confirmado y acotado** (`TestLegacyScheduledFrequency_Monthly_NegativeDayResolvedAgainstServerLocalMonth`). El día negativo se resuelve con `GetMaxDayOfMonth(currentTime...)`, donde `currentTime` es el reloj **local del servidor**. Cuando la fecha local y la UTC caen en meses distintos —lo que ocurre en cada cambio de mes para cualquier servidor que no esté en UTC— el último día del mes no dispara. El test pinea hoy `false` y **debe pasar a `true` en la Fase 1**; ahí se mueve a la suite del motor corregido como test de regresión.
2. **La zona horaria de la plantilla corre el día evaluado** (`TestLegacyScheduledFrequency_Weekly_TemplateTimezoneShiftsTheEvaluatedDay`). Con `ScheduledAt = 60` y offset −180, el instante `2026-08-17T01:00Z` se evalúa como **domingo 16**, no lunes 17. Confirma que `GetScheduledOccurrences` no puede recibir una fecha suelta.
3. **`ScheduledStartTime` se compara contra el instante programado, no contra el inicio del día** (`TestLegacyScheduledFrequency_StartTimeIsComparedAgainstTheScheduledInstant`). Una plantilla que arranca al mediodía y está programada para las 08:00 **no dispara en su propia fecha de inicio**. No estaba previsto en el plan original y el motor debe replicarlo.

---

## Fase 1 — Backend: motor de frecuencias reutilizable

- [ ] **1.1** Crear `pkg/services/scheduled_frequency.go` con:

  ```go
  // Devuelve todos los instantes en que la plantilla dispararía dentro de (from, to],
  // resueltos en la zona horaria de la plantilla y respetando
  // ScheduledStartTime / ScheduledEndTime.
  func GetScheduledOccurrences(template *models.TransactionTemplate, from, to time.Time) ([]time.Time, error)

  // Predicado de una única ocurrencia. occurrenceTime ya debe venir en la zona
  // horaria de la plantilla (equivalente al transactionTime del cron).
  func matchesScheduledFrequency(template *models.TransactionTemplate, frequencyValues []int64, occurrenceTime time.Time) bool
  ```

  Decisiones de diseño obligatorias:

  - **No usar una firma por fecha suelta.** El cron no evalúa un `date`, sino el instante `inicio_del_día_UTC + ScheduledAt minutos` convertido a la zona de la plantilla ([transactions.go:856-858](../../pkg/services/transactions.go#L856-L858)). `GetScheduledOccurrences` debe construir ese mismo instante para cada día del rango.
  - **Corregir el bug del día negativo.** Hoy `MONTHLY` con día negativo resuelve `GetMaxDayOfMonth` sobre `currentTime`, que es la hora *local del servidor* ([transactions.go:846](../../pkg/services/transactions.go#L846)). Debe calcularse sobre el **mes de la ocurrencia evaluada** y en la zona de la plantilla. Sin esto, proyectar meses hacia adelante da resultados incorrectos.
  - Sin dependencias de `currentUnixTime`, del intervalo del cron ni de `core.Context`: función pura, testeable en aislamiento.

- [ ] **1.2** Portar los tests de la Fase 0 a `pkg/services/scheduled_frequency_test.go`, apuntando a `GetScheduledOccurrences`. Agregar:
  - Rango que abarca varios meses (verificar cantidad y fechas exactas de ocurrencias).
  - `from` estrictamente exclusivo y `to` inclusivo.
  - Plantilla con `ScheduledEndTime` en mitad del rango → las ocurrencias se cortan ahí.
  - Zonas horarias no triviales (`ScheduledTimezoneUtcOffset` = -180 para Argentina, y un offset positivo).
  - Caso del día negativo cruzando meses de distinta longitud, que es lo que corrige 1.1.

- [ ] **1.3** Modificar `CreateScheduledTransactions` ([pkg/services/transactions.go:845-890](../../pkg/services/transactions.go#L845-L890)) para delegar en `matchesScheduledFrequency`, eliminando el bloque duplicado. Mantener intactos el logging, los contadores (`skipCount`, etc.) y las validaciones previas de la función.

- [ ] **1.4** Verificar que los tests de la Fase 0 siguen pasando sin modificarlos.

**Criterio de aceptación:** `go build ./... && go test ./pkg/...` en verde, y los tests de caracterización de la Fase 0 pasan sin cambios respecto de su versión original — salvo el caso del día negativo, cuyo valor esperado se actualiza documentando que era un bug.

---

## Fase 2 — Backend: modelo de request

- [ ] **2.1** En [pkg/models/transaction.go](../../pkg/models/transaction.go), agregar junto a `TransactionStatisticTrendsRequest` (línea 304):

  ```go
  // TransactionProjectionRequest represents all parameters of transaction projection request
  type TransactionProjectionRequest struct {
      YearMonthRangeRequest
      UseTransactionTimezone bool `form:"use_transaction_timezone"`
  }
  ```

- [ ] **2.2** **No crear DTOs de respuesta nuevos.** Se reutilizan `TransactionStatisticTrendsResponseItem` (línea 475) y `TransactionStatisticResponseItem` (línea 836) tal cual. Verificar antes de seguir que:
  - `TransactionStatisticResponseItem` incluye `AccountId` — es lo que permite la conversión de moneda en el frontend.
  - `TransactionStatisticResponseItem` **no** tiene campo `Type` — ingreso vs. egreso se deriva del tipo de la categoría en el frontend, el backend no lo envía.

**Criterio de aceptación:** `go build ./...` compila.

---

## Fase 3 — Backend: servicio de proyección

- [ ] **3.1** Crear `pkg/services/transaction_projections.go` con un **método de la struct existente `TransactionService`** (no una struct nueva: crear una obligaría a cablear inyección de dependencias sin ganancia):

  ```go
  func (s *TransactionService) GetProjectedCategoryAmountsByMonth(
      c core.Context, uid int64,
      startYear, startMonth, endYear, endMonth int32,
      clientTimezone *time.Location, useTransactionTimezone bool,
      currentUnixTime int64,
  ) (map[int32][]*models.TransactionTotalAmount, error)
  ```

  `currentUnixTime` se recibe por parámetro en lugar de llamar a `time.Now()` internamente, para poder testear el corte.

- [ ] **3.2** Implementar la **validación del rango**:
  - `uid > 0`, `start` y `end` presentes, `start <= end`.
  - Tope máximo de meses (sugerido: 60). El lado real carga en memoria todas las transacciones del período, así que un rango abierto es un vector de abuso.
  - Devolver `errs.ErrIncompleteOrIncorrectSubmission` / el error apropiado del paquete `errs` en cada caso.

- [ ] **3.3** Implementar el **lado real**:
  - Llamar a `GetAccountsAndCategoriesMonthlyInflowAndOutflow` ([pkg/services/transactions.go:2551](../../pkg/services/transactions.go#L2551)) con el rango pedido, sin filtros de tag ni keyword.
  - **Filtrar transferencias**: descartar `TRANSACTION_DB_TYPE_TRANSFER_OUT` y `TRANSACTION_DB_TYPE_TRANSFER_IN`. Si no se hace, los meses pasados muestran transferencias y los futuros no.
  - **Filtrar por el corte en `now`**: descartar todo lo que tenga `transaction_time > currentUnixTime`. Esta función devuelve **todo** el rango de meses, incluidas transacciones manuales con fecha futura; sin este filtro los dos conjuntos se solapan y hay doble conteo.

  > Como `GetAccountsAndCategoriesMonthlyInflowAndOutflow` ya agrega y descarta el `transaction_time` individual, hará falta o bien un método hermano que acepte un `maxTransactionTime`, o bien reimplementar la consulta acotada. Preferir extender la función existente con un parámetro opcional de tiempo máximo antes que duplicar la query.

- [ ] **3.4** Implementar el **lado simulado**:
  - Traer las plantillas del usuario con `TemplateType = TRANSACTION_TEMPLATE_TYPE_SCHEDULE`, `Deleted = false` y `ScheduledFrequencyType != DISABLED`.
  - **No filtrar por `Hidden`**: el cron tampoco lo hace, y proyección y realidad deben coincidir.
  - Excluir las plantillas con `Type = TRANSACTION_TYPE_TRANSFER`.
  - Para cada plantilla, `GetScheduledOccurrences(template, now, finDelRango)`.
  - Cada ocurrencia aporta `template.Amount` a la clave `(Year, Month, CategoryId, AccountId)`.

- [ ] **3.5** Implementar el **bucketeo por mes de las ocurrencias simuladas usando la misma regla de zona horaria que el lado real**: `clientTimezone`, o la zona propia si `useTransactionTimezone = true`. Si el mes se calcula con la zona de la plantilla y los reales con la del cliente, la misma ocurrencia puede caer en meses distintos.

- [ ] **3.6** Implementar el **merge**: sumar ambos conjuntos por `(Year, Month, CategoryId, AccountId)`, devolviendo `map[yearMonth][]*models.TransactionTotalAmount`, el mismo tipo que devuelve el servicio de tendencias.

  > **La clave debe conservar `AccountId`.** El `Amount` de una plantilla está expresado en la moneda de su cuenta; agregar sin `AccountId` mezcla monedas y hace imposible la conversión en el frontend.

- [ ] **3.7** Tests en `pkg/services/transaction_projections_test.go`:
  - Mes completamente pasado → solo importes reales.
  - Mes completamente futuro → solo importes simulados.
  - Mes en curso → mezcla de ambos, sin doble conteo, con `currentUnixTime` fijo a mitad de mes.
  - Transacción manual con fecha futura dentro del rango → **no** se cuenta (queda del lado simulado del corte).
  - Plantilla de transferencia → no aporta; transacción real de transferencia → no aporta.
  - Dos plantillas de la misma categoría pero **cuentas con monedas distintas** → se devuelven como dos ítems separados, no sumados.
  - Rango inválido y rango que excede el tope → error.

**Criterio de aceptación:** `go test ./pkg/services/...` en verde, incluidos los casos de multi-moneda y de corte en `now`.

---

## Fase 4 — Backend: API y ruta

- [ ] **4.1** Crear `pkg/api/transaction_projections.go` con `TransactionProjectionsHandler` como **método de la struct existente `TransactionsApi`**, siguiendo el esqueleto de `TransactionStatisticsTrendsHandler` ([pkg/api/transactions.go:597-670](../../pkg/api/transactions.go#L597-L670)):
  1. `c.ShouldBindQuery(&projectionReq)`.
  2. `c.GetClientTimezone()` → error `errs.ErrClientTimezoneOffsetInvalid` si falla.
  3. `projectionReq.GetNumericYearMonthRange()`.
  4. `uid := c.GetCurrentUid()`.
  5. Llamar al servicio pasando `time.Now().Unix()` como `currentUnixTime`.
  6. Mapear el `map[int32][]*models.TransactionTotalAmount` a `[]*models.TransactionStatisticTrendsResponseItem`, incluyendo **siempre** `AccountId`, y ordenar con `sort.Sort` sobre `TransactionStatisticTrendsResponseItemSlice`.

- [ ] **4.2** Registrar la ruta en [cmd/webserver.go](../../cmd/webserver.go), inmediatamente después de la línea 411 (`asset_trends.json`):

  ```go
  apiV1Route.GET("/transactions/statistics/projections.json", bindApi(api.Transactions.TransactionProjectionsHandler, config))
  ```

  Se ubica bajo `/transactions/statistics/` por consistencia con `trends.json` y `asset_trends.json`, no en la raíz de `/transactions/`.

- [ ] **4.3** Prueba manual con `curl` sobre un usuario con plantillas programadas activas:
  - Rango solo pasado, solo futuro, y a caballo del mes en curso.
  - Verificar que cada ítem trae `accountId`.
  - Verificar que no aparecen ítems de categorías de transferencia.
  - Rango invertido y rango de 100 meses → error 400.

**Criterio de aceptación:** el endpoint responde con la misma forma que `trends.json` y los tres escenarios de rango dan valores coherentes con los datos de prueba.

---

## Fase 5 — Frontend: capa de datos

- [ ] **5.1** En [src/models/transaction.ts](../../src/models/transaction.ts), agregar junto a `TransactionStatisticTrendsRequest` (línea 693):

  ```ts
  export interface TransactionProjectionRequest extends YearMonthRangeRequest {
      readonly useTransactionTimezone: boolean;
  }
  ```

  Los tipos de respuesta se reutilizan (`TransactionStatisticTrendsResponseItem`, `TransactionStatisticResponseItem`) — no crear tipos paralelos.

- [ ] **5.2** En [src/lib/services.ts](../../src/lib/services.ts), agregar `getTransactionProjections(req)` con el mismo esqueleto que `getTransactionStatisticsTrends` (líneas 573-597): arma `start_year_month` / `end_year_month` y llama a `GET v1/transactions/statistics/projections.json`.

- [ ] **5.3** Crear `src/stores/projection.ts` siguiendo el patrón de [src/stores/statistics.ts](../../src/stores/statistics.ts):
  - **Estado**: `startYearMonth`, `endYearMonth`, datos crudos, `projectionDataLoaded`.
  - **Acción** `loadProjections({ force })`, análoga a `loadTrendAnalysis` ([statistics.ts:1899](../../src/stores/statistics.ts#L1899)).
  - **Enriquecimiento**: resolver `category` / `primaryCategory` / `account` y convertir a la moneda por defecto vía `exchangeRatesStore.getExchangedAmount(...)`, replicando `assembleAccountAndCategoryInfo` ([statistics.ts:1030-1088](../../src/stores/statistics.ts#L1030)). Los ítems sin cuenta resuelta quedan con `amountInDefaultCurrency = null` y deben poder señalarse en la UI.
  - **Clasificación** ingreso/egreso por `item.category.type === CategoryType.Income | Expense` ([src/core/category.ts:4](../../src/core/category.ts#L4)) — el backend no envía ese dato.

- [ ] **5.4** Computeds del store que arman las filas de la tabla:
  - Árbol `sección → categoría → subcategoría` con un importe por mes.
  - `Subtotal <categoría>` por mes.
  - `Total Ingresos` / `Total Egresos` por mes.
  - `Neto del mes` = Total Ingresos − Total Egresos.
  - `Acumulado` = suma corrida del Neto mes a mes.
  - Columna **Total**: suma horizontal para todas las filas, **excepto** `Acumulado`, cuya celda Total es el valor del **último mes** del período (ya es un acumulado, no se vuelve a sumar).

- [ ] **5.5** **Invalidación de caché**: el resultado depende de `now` **y** del conjunto de plantillas. Invalidar al crear/editar/borrar una plantilla programada y al cambiar de día. Cachear solo por período seleccionado es incorrecto.

- [ ] **5.6** Tests con vitest de los computeds de 5.4: subtotales, neto, acumulado, y el caso de la celda Total de la fila Acumulado.

**Criterio de aceptación:** `npm run lint && npm run test` en verde.

---

## Fase 6 — Frontend: componente de tabla

- [ ] **6.1** Crear `src/components/desktop/ProjectionTable.vue`: una columna por mes del período + columna **Total**.

- [ ] **6.2** Implementar la jerarquía de 3 niveles con `v-expansion-panels multiple` (mismo patrón que [PresetDialog.vue](../../src/views/desktop/categories/list/dialogs/PresetDialog.vue)):

  | Nivel | Fila | Colapsable | Al colapsar |
  |---|---|---|---|
  | 0 | Sección `INGRESOS` / `EGRESOS` | Sí | Se ocultan categorías y subtotales; **`Total Ingresos` / `Total Egresos` permanece visible** |
  | 1 | Categoría (*Trabajo*, *Casa*…) | Sí | Se ocultan sus subcategorías; **`Subtotal <categoría>` permanece visible** |
  | 2 | Subcategoría (*Salario*, *Recibos*…) | No | — |

  Los dos niveles se controlan de forma independiente. El estado de expansión vive en el componente, no en el store.

- [ ] **6.3** Filas de cierre, siempre visibles y fuera del acordeón: `Neto del mes` y `Acumulado`.

- [ ] **6.4** Formateo de importes reutilizando `BigDecimal` y los helpers de `numeral.ts` ya usados en el resto de la app. Manejar el caso de importe sin conversión de moneda disponible.

- [ ] **6.5** Marcar visualmente la separación entre meses ya transcurridos (datos reales) y meses proyectados, para que el usuario entienda por qué las columnas no son directamente comparables.

**Criterio de aceptación:** el componente renderiza el ejemplo de [Technical-spec.md](./Technical-spec.md) con los mismos números, y el acordeón se comporta según la tabla de 6.2.

---

## Fase 7 — Frontend: página, ruta y navegación

- [ ] **7.1** Crear `src/views/desktop/projections/TransactionPage.vue`: selector de período reutilizando [MonthRangeSelectionDialog.vue](../../src/components/desktop/MonthRangeSelectionDialog.vue) + `ProjectionTable`.

- [ ] **7.2** Definir el horizonte por defecto al entrar por primera vez. Sugerido: **mes en curso + 11 meses hacia adelante**, sin meses pasados, para que la tabla se lea como una proyección y no como un histórico.

- [ ] **7.3** En [src/router/desktop.ts](../../src/router/desktop.ts): agregar el import de la página (junto al de `StatisticsTransactionPage`, línea 19) y la ruta `/projections/transaction` con el mismo guard de login que el resto.

- [ ] **7.4** En [src/views/desktop/MainLayout.vue](../../src/views/desktop/MainLayout.vue): agregar un `<li class="nav-link">` en la sección "Transaction Data", después de "Insights Explorer" (líneas 43-48), con `router-link to="/projections/transaction"` y un ícono `mdi` acorde.

**Criterio de aceptación:** la pestaña aparece en la navegación, navega correctamente y carga datos reales del backend.

---

## Fase 8 — i18n

- [ ] **8.1** Agregar las claves nuevas en los **20** archivos de [src/locales/](../../src/locales/):
  `"Projections"`, `"Net"`, `"Accumulated"`, `"Subtotal"`, `"Projected"`.

- [ ] **8.2** **Reutilizar, no duplicar:** `"Total Income"` y `"Total Expense"` ya existen ([en.json:2307-2308](../../src/locales/en.json#L2307-L2308)). Verificar con `grep` cada clave nueva antes de agregarla.

- [ ] **8.3** Completar traducciones reales en `en.json` y `es.json`; en los idiomas restantes, seguir la convención del repo (clave inglesa → valor traducido) y traducir con el mismo criterio que las claves vecinas de Statistics.

**Criterio de aceptación:** `npm run build` sin warnings de claves faltantes, y la UI en español no muestra claves crudas.

---

## Fase 9 — Verificación final

- [ ] **9.1** Escenario end-to-end: crear plantillas programadas con las cinco frecuencias, esperar/forzar la materialización de al menos una, y verificar que la tabla muestra importes reales en meses pasados y proyectados en los futuros.
- [ ] **9.2** Caso multi-moneda: una plantilla sobre cuenta en ARS y otra sobre cuenta en USD, en la misma categoría. Verificar que el total se muestra convertido a la moneda por defecto y no como suma cruda.
- [ ] **9.3** Caso del mes en curso: verificar que lo ya transcurrido sale de transacciones reales y lo que falta del mes sale de la simulación, sin duplicados.
- [ ] **9.4** Caso de transferencias: una plantilla de transferencia y una transferencia real en el pasado → ninguna aparece en la tabla.
- [ ] **9.5** Caso de borde de calendario: plantilla `MONTHLY` día 31 → los meses cortos aparecen en cero; plantilla `MONTHLY` día -1 → dispara el último día de cada mes, correcto en febrero.
- [ ] **9.6** Acordeón: colapsar una categoría oculta solo sus subcategorías y mantiene su subtotal; colapsar Ingresos/Egresos oculta la sección y mantiene su total.
- [ ] **9.7** Neto y Acumulado contra el ejemplo de [Technical-spec.md](./Technical-spec.md), incluida la celda Total de la fila Acumulado.
- [ ] **9.8** `go test ./pkg/... && npm run lint && npm run test && npm run build`.
- [ ] **9.9** `/code-review` sobre la rama antes de abrir el PR.

---

## Fuera de alcance de este plan

- **Mobile** (`src/views/mobile/`): no se agrega página. Insights Explorer tampoco la tiene, así que no rompe la paridad existente.
- **Transferencias**: excluidas de ambos lados del cálculo.
- **Saldo inicial en el acumulado**: `Neto` y `Acumulado` reflejan solo el flujo del período. Sumar `getNetAssets()` ([src/stores/account.ts:472](../../src/stores/account.ts#L472)) para ver patrimonio absoluto queda para v2.
- **Huecos del cron**: si el servidor estuvo caído, una ocurrencia pasada no se materializó y tampoco se proyecta (ya es pasado). Limitación conocida, sin mitigación en v1.
