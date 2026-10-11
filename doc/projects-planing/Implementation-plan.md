# Planificación de Proyectos — Plan de implementación

Plan paso a paso derivado de [Technical-spec.md](./Technical-spec.md). El orden es de abajo hacia arriba: primero las tablas y la lógica pura que se puede testear sin base de datos, después los servicios que la usan, la integración con proyecciones, la API, y por último el frontend (datos → pantallas → diálogos → navegación → i18n).

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

> **Limitación del entorno de tests, heredada de la feature de proyecciones.** `pkg/services` no tiene harness de base de datos: los tests existentes solo cubren helpers puros. Todo lo que consulte la base queda fuera del alcance de los tests automáticos y se verifica a mano contra una instancia levantada. Por eso este plan empuja deliberadamente la lógica de negocio hacia funciones puras (estado derivado, cálculo del ahorro, imputación por mes), que sí son testeables, y deja en los servicios solo el acceso a datos.

---

## Fase 1 — Modelos y tablas

Los cuatro modelos de la [sección 3](./Technical-spec.md#3-modelos-nuevos), sin lógica todavía.

- [ ] **1.1** Crear `pkg/models/project.go` con la struct `Project` y el enum `ProjectStatus`:

  ```go
  type ProjectStatus byte

  const (
      PROJECT_STATUS_DRAFT       ProjectStatus = 1
      PROJECT_STATUS_PLANNED     ProjectStatus = 2
      PROJECT_STATUS_IN_PROGRESS ProjectStatus = 3
      PROJECT_STATUS_DELAYED     ProjectStatus = 4
      PROJECT_STATUS_COMPLETED   ProjectStatus = 5
  )
  ```

- [ ] **1.2** Crear `pkg/models/project_item.go`, `pkg/models/project_item_transaction.go` y `pkg/models/project_scheduled_transaction.go`.

  **Verificar antes de seguir** que las dos tablas de relación usan como clave primaria el id de la entidad enlazada (`TransactionId` y `TemplateId`) y **no** un uuid propio. Es lo que mantiene el consumo de tipos de uuid en dos (paso 1.3).

- [ ] **1.3** En [pkg/uuid/uuid_type.go](../../pkg/uuid/uuid_type.go) agregar **exactamente dos** tipos:

  ```go
  UUID_TYPE_PROJECT      UuidType = 12
  UUID_TYPE_PROJECT_ITEM UuidType = 13
  ```

  > ⚠️ `UuidType` se empaqueta en 4 bits ([internal_generator.go:64](../../pkg/uuid/internal_generator.go#L64)): el rango válido es 0–15 y no es ampliable sin cambiar el formato de todos los ids existentes. Después de esta fase quedan libres **14 y 15**. Mencionarlo en el PR.

- [ ] **1.4** Registrar las cuatro tablas en `updateAllDatabaseTablesStructure` ([cmd/database.go:45](../../cmd/database.go#L45)), al final, con su `log.BootInfof` correspondiente.

**Criterio de aceptación:** `go build ./...` compila y `ezbookkeeping database update` crea las cuatro tablas con sus índices en una base limpia y en una base existente.

---

## Fase 2 — Estado derivado y errores

Lógica pura: es la parte del backend que sí admite tests.

- [ ] **2.1** Implementar en `pkg/models/project.go`:

  ```go
  // GetStatus devuelve el estado del proyecto en el instante indicado. Las tres fechas se comparan
  // normalizadas al inicio del día en la zona horaria recibida, para que un proyecto cuya fecha de
  // ejecución es hoy esté en ejecución durante todo el día.
  func (p *Project) GetStatus(nowUnixTime int64, tz *time.Location) ProjectStatus
  ```

  Tabla de decisión completa en la [sección 2](./Technical-spec.md#2-estado-del-proyecto). El orden de evaluación importa: `Completed` primero, ausencia de `ExecutionTime` después.

- [ ] **2.2** Crear `pkg/services/projects_test.go` con los casos de estado:
  - Completado con y sin fechas → `COMPLETED` en todos.
  - Sin fecha de ejecución, con y sin fecha límite → `DRAFT` en ambos.
  - Antes / el mismo día / después de la fecha de ejecución.
  - El mismo día de la fecha límite → todavía `IN_PROGRESS`, no `DELAYED`.
  - Un día después de la fecha límite → `DELAYED`.
  - Con fecha de ejecución pasada y sin fecha límite → `IN_PROGRESS`, nunca `DELAYED`.
  - Fecha de ejecución "hoy" evaluada a las 00:30 y a las 23:30 en zonas `-0300` y `+0900` → mismo resultado. Este es el caso que rompe una implementación que compare instantes en vez de días.

- [ ] **2.3** Crear `pkg/errs/project.go` con los códigos de la [sección 8](./Technical-spec.md#errores-nuevos) y agregar `NormalSubcategoryProject = 21` en [pkg/errs/error.go](../../pkg/errs/error.go#L46).

**Criterio de aceptación:** `go test ./pkg/services/...` en verde, con los casos de zona horaria incluidos.

---

## Fase 3 — Servicio de proyectos

- [ ] **3.1** Crear `pkg/services/projects.go` con `ProjectService`, siguiendo el patrón de [pkg/services/explorer.go](../../pkg/services/explorer.go): singleton con `ServiceUsingDB` + `ServiceUsingUuid`.

  Métodos de lectura: `GetAllProjectsByUid`, `GetProjectByProjectId`, `GetMaxDisplayOrder`, `GetItemsByProjectId`.

- [ ] **3.2** Métodos de escritura: `CreateProject`, `ModifyProject`, `CompleteProject`, `HideProject`, `ModifyProjectDisplayOrders`, `DeleteProject`, `DeleteAllProjects`, y los cuatro equivalentes de ítems.

  `DeleteProject` borra lógicamente el proyecto, sus ítems y sus filas de relación en una sola transacción de base de datos. **No borra las transacciones de las compras**: son gasto real ya ocurrido y borrarlas alteraría saldos de cuentas por un cambio de planificación.

- [ ] **3.3** Implementar las validaciones de la [sección 3.1](./Technical-spec.md#31-tabla-project): categoría de gasto de segundo nivel y no oculta, cuenta existente, no oculta y no `MULTI_SUB_ACCOUNTS`, presupuesto en rango, fecha límite no anterior a la de ejecución.

  Las dos primeras replican lo que `isCategoryValid` ([transactions.go:3363](../../pkg/services/transactions.go#L3363)) y `doCreateTransaction` ([transactions.go:2700](../../pkg/services/transactions.go#L2700)) exigen a una transacción. Si divergen, se puede crear un proyecto cuyas compras fallan al registrarse.

- [ ] **3.4** Implementar `GetSpentAmountsByProject` con la regla de la [sección 4](./Technical-spec.md#4-costo-real-se-calcula-no-se-guarda): por ítem, el importe de la transferencia si existe, si no el del gasto; unido contra `transaction` con `deleted = false`; con corte opcional por `transaction_time`.

- [ ] **3.5** Agregar a `projects_test.go` los tests de la función pura que resuelve el importe por ítem a partir de sus filas de relación (transferencia gana sobre gasto, relación sin transacción se ignora, ítem sin compras aporta cero).

**Criterio de aceptación:** `go build ./...` y `go test ./pkg/services/...` en verde. Alta, edición y borrado de un proyecto verificados a mano contra una instancia levantada.

---

## Fase 4 — Registrar compra

La fase con más riesgo: crea transacciones reales y mueve saldos.

- [ ] **4.1** Extraer `createTransactionInSession` de `CreateTransaction` ([transactions.go:567-636](../../pkg/services/transactions.go#L567-L636)).

  Es un movimiento literal: todo el cuerpo actual salvo el `userDataDb.DoTransaction(...)` final pasa al método nuevo, y `CreateTransaction` queda como la envoltura que abre la transacción de base de datos y lo llama.

  > **No hay forma de cubrir esto con tests** (ver la nota del encabezado). La verificación es que sea un movimiento sin cambios de lógica —revisar el diff línea por línea— y que el comportamiento observable de `CreateTransaction` sea idéntico: mismos uuids pedidos, misma normalización de `TransactionTime`, mismos índices de tags, mismo modelo de imágenes.

- [ ] **4.2** Crear `pkg/services/project_purchases.go` con:

  ```go
  func (s *ProjectService) RegisterItemPurchase(c core.Context, uid int64, req *models.ProjectItemPurchaseRequest) (*models.ProjectItemPurchaseResponse, error)
  func (s *ProjectService) RevokeItemPurchase(c core.Context, uid int64, itemId int64) error
  ```

  Todo dentro de un solo `DoTransaction`: transferencia (si la cuenta de compra difiere de la del proyecto) → gasto → filas de relación. Orden y detalles en la [sección 7](./Technical-spec.md#7-registrar-compra).

- [ ] **4.3** Rechazar el segundo intento sobre un ítem que ya tiene compra activa con `ErrProjectItemAlreadyPurchased`.

- [ ] **4.4** Verificar a mano, contra una instancia levantada, los cinco escenarios:
  1. Compra con la cuenta del proyecto → una sola transacción, saldo de la cuenta baja el costo real.
  2. Compra con otra cuenta de la misma moneda → transferencia + gasto; el saldo de la cuenta de compra vuelve a su valor original y el de la cuenta del proyecto baja.
  3. Compra con otra cuenta de distinta moneda → los dos importes de la transferencia se respetan.
  4. Anular una compra → las transacciones desaparecen, los saldos vuelven, el ítem queda disponible.
  5. Dos compras en el mismo segundo → ninguna colisión de `transaction_time`, gracias al corrimiento de [transactions.go:2822](../../pkg/services/transactions.go#L2822).

**Criterio de aceptación:** los cinco escenarios de 4.4 verificados, y `go test ./pkg/...` sin regresiones respecto de la línea base.

---

## Fase 5 — Ahorro programado y propagación

- [ ] **5.1** Implementar en `pkg/services/projects.go` el cálculo del importe por ocurrencia como **función pura**, para poder testearla:

  ```go
  // Devuelve el importe de cada ocurrencia y cuántas ocurrencias quedan hasta la fecha de ejecución.
  func calculateSavingOccurrenceAmount(template *models.TransactionTemplate, pendingAmount int64, from time.Time, executionTime time.Time) (amount int64, occurrences int, err error)
  ```

  Usa `GetScheduledOccurrences` ([pkg/services/scheduled_frequency.go](../../pkg/services/scheduled_frequency.go)) para contar, de modo que la cuenta coincida exactamente con lo que el cron va a ejecutar. Redondeo hacia arriba.

- [ ] **5.2** Tests del cálculo: división exacta, división con resto (redondea hacia arriba), pendiente cero, cero ocurrencias → error, frecuencia mensual día 31 en un rango de meses cortos (no dispara, y la cuenta lo refleja).

- [ ] **5.3** Implementar `CreateProjectScheduledSaving` y `DeleteProjectScheduledSaving`: plantilla + fila de relación en un solo `DoTransaction`. Campos fijados por el proyecto según la tabla de la [sección 5](./Technical-spec.md#creación).

- [ ] **5.4** Implementar la propagación dentro de `ModifyProject`, sobre todas las plantillas enlazadas y en la misma transacción de base de datos: `Name`, `RelatedAccountId`, `ScheduledEndTime` y recálculo del importe.

  **La categoría no se propaga.** Dejar un comentario en el código explicando por qué —`isCategoryValid` rechaza una transferencia con categoría de gasto ([transactions.go:3386](../../pkg/services/transactions.go#L3386))—, o alguien va a "arreglarlo" más adelante.

- [ ] **5.5** Si el recálculo da cero ocurrencias, **no tocar el importe** y devolver la advertencia en la respuesta. Nunca fallar la modificación del proyecto por esto.

- [ ] **5.6** `CompleteProject` oculta las plantillas enlazadas (`Hidden = true`), no las borra.

**Criterio de aceptación:** `go test ./pkg/services/...` en verde. Verificado a mano: crear un ahorro, cambiar el presupuesto del proyecto y comprobar que el importe de la plantilla se recalcula; completar el proyecto y comprobar que la plantilla queda oculta y el cron deja de mostrarla en la lista.

---

## Fase 6 — Integración con las proyecciones

- [ ] **6.1** Crear `pkg/services/project_projections.go` con `GetProjectedProjectAmountsByMonth`, siguiendo los seis pasos de la [sección 6](./Technical-spec.md#el-tercer-sumando).

  La parte pura —dado un proyecto, su pendiente y el mes actual, decidir si aporta y en qué mes— va en una función separada, testeable sin base de datos.

- [ ] **6.2** Modificar `GetProjectedCategoryAmountsByMonth` ([transaction_projections.go:41](../../pkg/services/transaction_projections.go#L41)) para fusionar el tercer conjunto. `mergeMonthlyTotalAmounts` ya agrupa por `(mes, categoría, cuenta)`; no hace falta tocarlo.

- [ ] **6.3** Crear `pkg/services/project_projections_test.go`:
  - Proyecto planificado a tres meses → aporta el presupuesto completo en su mes.
  - Proyecto con la mitad comprada → aporta la mitad.
  - Proyecto totalmente comprado → no aporta fila.
  - Proyecto completado → no aporta, con o sin pendiente.
  - Proyecto en borrador → no aporta.
  - Proyecto **retrasado** (fecha de ejecución el mes pasado, pendiente > 0) → aporta en el **mes actual**, no en el pasado. Es el caso del `max` del paso 4 de la sección 6.
  - Proyecto cuyo mes de ejecución queda fuera del rango pedido → no aporta.
  - Compra con fecha futura → **no** se descuenta del pendiente, porque el lado real tampoco la cuenta. Es el caso que justifica el parámetro de corte de `GetSpentAmountsByProject`.
  - Proyecto y plantilla programada con la misma categoría y cuenta → una sola celda con la suma, no dos filas.

**Criterio de aceptación:** `go test ./pkg/services/...` en verde, incluidos los nueve casos. Verificado a mano en la tabla de proyecciones: el presupuesto pendiente aparece bajo EGRESOS en la subcategoría del proyecto, y al registrar una compra el total de esa celda **no cambia** (baja el pendiente, sube el gasto real).

---

## Fase 7 — API y rutas

- [ ] **7.1** Agregar los requests y `ProjectInfoResponse` / `ProjectItemInfoResponse` a `pkg/models/project.go`, con los `binding` del patrón del proyecto (`required`, `notBlank`, `max=64`, `min=1`, ids como `json:"...,string"`).

  `ActualAmount` y `Status` son solo de salida.

- [ ] **7.2** Crear `pkg/api/projects.go` con `ProjectsApi` y los quince handlers, siguiendo el patrón de [pkg/api/explorers.go](../../pkg/api/explorers.go): bind, `c.GetCurrentUid()`, llamada al servicio, log con el prefijo `[projects.XxxHandler]`, `errs.Or(err, errs.ErrOperationFailed)`.

  `ProjectListHandler` resuelve estado y costo real de todos los proyectos **en dos consultas** (proyectos + costos agregados), no una por proyecto.

- [ ] **7.3** Registrar las rutas en [cmd/webserver.go](../../cmd/webserver.go) después del bloque de exploradores (línea 483).

- [ ] **7.4** Agregar `DeleteAllProjects` a `ClearAllDataHandler` ([data_managements.go:170](../../pkg/api/data_managements.go#L170)), **antes** de `DeleteAllTransactions`.

- [ ] **7.5** Verificar cada endpoint con `curl` contra una instancia levantada, incluidos los caminos de error: categoría de ingreso, categoría de primer nivel, cuenta oculta, fecha límite anterior a la de ejecución, ahorro sobre un proyecto en borrador, compra duplicada.

**Criterio de aceptación:** los quince endpoints responden lo esperado en el camino feliz y en los seis casos de error de 7.5.

---

## Fase 8 — Frontend: capa de datos

- [ ] **8.1** Crear `src/core/project.ts` con el enum `ProjectStatus` espejo del backend y su información de presentación (clave i18n, color, orden de listado), siguiendo el patrón de [src/core/projection.ts](../../src/core/projection.ts).

  **El estado no se recalcula en el cliente.** Este archivo solo lo presenta.

- [ ] **8.2** Crear `src/models/project.ts` con las clases `Project` y `ProjectItem` y las interfaces de request/response, siguiendo [src/models/explorer.ts](../../src/models/explorer.ts).

- [ ] **8.3** Crear `src/lib/project.ts` con las funciones puras: porcentaje de progreso, pendiente, suma de presupuestos de ítems, y la comparación entre el presupuesto del proyecto y el de sus ítems.

- [ ] **8.4** Crear `src/lib/__tests__/project.test.ts` con los tests de 8.3, incluidos presupuesto cero (no dividir por cero) y costo real mayor que el presupuesto (progreso tope 100%, pendiente cero).

- [ ] **8.5** Agregar los quince métodos a [src/lib/services.ts](../../src/lib/services.ts), con el esqueleto de los de exploradores (líneas 854-874).

- [ ] **8.6** Crear `src/stores/project.ts` según la [sección 10](./Technical-spec.md#store-srcstoresprojectts).

**Criterio de aceptación:** `npm run lint` y `npm run test` en verde.

---

## Fase 9 — Frontend: lista y detalle

- [ ] **9.1** Crear `src/components/desktop/ProjectProgressBar.vue` y registrarlo en [src/desktop-main.ts](../../src/desktop-main.ts#L594).

- [ ] **9.2** Crear `src/views/desktop/projects/ListPage.vue` con la maqueta de la [sección 9.1](./Technical-spec.md#91-lista-de-proyectos--projectslist): tarjetas con estado, categoría, progreso y cuenta; reordenamiento por arrastrar y soltar como en [templates/ListPage.vue](../../src/views/desktop/templates/ListPage.vue); filtro por estado y conmutador de ocultos.

- [ ] **9.3** Crear `src/views/desktop/projects/DetailPage.vue` con las tres zonas de la [sección 9.2](./Technical-spec.md#92-detalle-del-proyecto--projectsid): encabezado con estado y fechas, bloque de ahorro programado, tabla de ítems con la fila de cierre que compara presupuestos.

- [ ] **9.4** Un ítem ya comprado muestra costo real, vendedor y fecha, con enlace a la transacción en `/transaction/list`.

**Criterio de aceptación:** `npm run lint` y `npm run build` en verde; ambas pantallas renderizan datos reales de la API con los cinco estados representados.

---

## Fase 10 — Frontend: diálogos

- [ ] **10.1** `dialogs/EditDialog.vue` — alta y edición del proyecto. Categoría con `two-column-select` filtrado a categorías de gasto de segundo nivel. Al editar un proyecto con ahorro enlazado, mostrar "el ahorro programado pasará de X a Y" **antes** de guardar.

- [ ] **10.2** `dialogs/ItemEditDialog.vue` — nombre, descripción, presupuesto, URL.

- [ ] **10.3** `dialogs/RegisterPurchaseDialog.vue` — costo real precargado con el presupuesto del ítem, vendedor, fecha, cuenta de compra. Si la cuenta difiere de la del proyecto, desplegar el bloque de transferencia con categoría de transferencia e importe convertido (sugerido con `exchangeRatesStore.getExchangedAmount`, editable).

- [ ] **10.4** `dialogs/ScheduleSavingDialog.vue` — cuenta de origen, categoría de transferencia, frecuencia con `ScheduleFrequencySelect`, fecha de inicio y hora. Mostrar en vivo el importe calculado y las ocurrencias restantes; deshabilitar el guardado si son cero.

  **No reutilizar [EditDialog.vue](../../src/views/desktop/transactions/list/dialogs/EditDialog.vue) de transacciones**: guardaría contra el endpoint de plantillas y rompería la atomicidad plantilla+relación.

- [ ] **10.5** `dialogs/CompleteDialog.vue` — confirmación, avisando que las plantillas de ahorro enlazadas quedarán ocultas.

**Criterio de aceptación:** `npm run lint` y `npm run build` en verde; los cinco diálogos completan su flujo de punta a punta contra la API.

---

## Fase 11 — Ruta, navegación e invalidación cruzada

- [ ] **11.1** Agregar `/projects/list` y `/projects/:id` a [src/router/desktop.ts](../../src/router/desktop.ts#L160), con `beforeEnter: checkLogin`.

- [ ] **11.2** Agregar el `<li class="nav-link">` en la sección "Transaction Data" de [MainLayout.vue](../../src/views/desktop/MainLayout.vue#L49-L54), después de "Projections".

- [ ] **11.3** Cablear la invalidación cruzada desde el store de proyectos:
  - Toda acción que altere presupuesto, categoría, fecha de ejecución, estado o compras → `projectionsStore.updateProjectionInvalidState(true)`.
  - Registrar o anular una compra → invalidar también el store de transacciones y el de cuentas, porque se crearon transacciones y se movieron saldos.

- [ ] **11.4** Verificar la cadena completa: crear un proyecto con fecha de ejecución a tres meses → abrir Proyecciones → el pendiente aparece en su subcategoría → registrar una compra → volver a Proyecciones **sin recargar la página** → el total de la celda no cambió.

**Criterio de aceptación:** 11.4 verificado en el navegador.

---

## Fase 12 — i18n

- [ ] **12.1** Agregar las claves de la [sección 10](./Technical-spec.md#i18n) a los **20 archivos** de `src/locales/`.

- [ ] **12.2** Verificar que no se duplican claves existentes: `"Category"`, `"Account"`, `"Description"`, `"Amount"` y `"Transfer"` ya están.

**Criterio de aceptación:** `npm run lint` en verde y ningún literal sin traducir en las pantallas nuevas.

---

## Fase 13 — Verificación final

- [ ] **13.1** `go build ./...`, `go test ./pkg/...`, `npm run lint`, `npm run test`, `npm run build`, todos en verde.

- [ ] **13.2** Recorrido completo contra una instancia limpia:
  1. Crear un proyecto sin fecha → aparece como **borrador**, sin botón de programar ahorro.
  2. Ponerle fecha de ejecución a tres meses → pasa a **planificado**.
  3. Programar un ahorro mensual → el importe por ocurrencia es el presupuesto dividido en tres, redondeado hacia arriba.
  4. Duplicar el presupuesto → el importe del ahorro se duplica, y el diálogo lo anticipó antes de guardar.
  5. Agregar cuatro ítems → la fila de cierre muestra la diferencia contra el presupuesto del proyecto.
  6. Registrar una compra con la cuenta del proyecto → baja el saldo, sube el costo real, baja el pendiente.
  7. Registrar una compra con otra cuenta → se crean transferencia y gasto; el saldo de la cuenta de compra queda igual que antes.
  8. Anular esa compra → ambas transacciones desaparecen y los saldos vuelven.
  9. Mover la fecha de ejecución al mes pasado → pasa a **retrasado** y en Proyecciones el pendiente se imputa al **mes actual**.
  10. Completar el proyecto → **completado**, las plantillas de ahorro quedan ocultas y desaparece de Proyecciones.
  11. Borrar el proyecto → desaparece, y las transacciones de sus compras **siguen existiendo** en la lista de transacciones.

- [ ] **13.3** `ezbookkeeping database update` sobre una base con datos anteriores a esta feature: las cuatro tablas se crean y nada existente se altera.

- [ ] **13.4** Borrar todos los datos del usuario desde la sección de datos → los proyectos, ítems y relaciones desaparecen sin dejar filas huérfanas.

**Criterio de aceptación:** los once pasos de 13.2 y los dos de 13.3-13.4 verificados.

---

## Fuera de alcance de este plan

Detallado en las [secciones 11 y 12](./Technical-spec.md#11-limitaciones-conocidas-de-v1) de la especificación. En resumen: mobile, ahorro programado de tipo gasto, propagación de categoría al ahorro, compras parciales y devoluciones, exportación/importación de proyectos, adjuntar imágenes, y fila propia del proyecto en la tabla de proyecciones.

## Nota para el PR

Tres cosas que hay que mencionar explícitamente al abrir el PR:

1. **Se extrajo `createTransactionInSession` de `CreateTransaction`** (fase 4.1). Es código en el camino crítico de toda alta de transacción y no está cubierto por tests, por la falta de harness de base de datos en `pkg/services`. Merece revisión línea por línea.
2. **Se consumieron dos de los cuatro tipos de uuid libres.** Quedan 14 y 15, y el espacio no es ampliable.
3. **La propagación de categoría del proyecto a su ahorro programado no se implementó**, por incompatibilidad entre categorías de gasto y de transferencia. Es la única desviación respecto del pedido original y está justificada en la [sección 1](./Technical-spec.md#decisión-que-se-aparta-del-pedido-original-y-por-qué) de la especificación.
