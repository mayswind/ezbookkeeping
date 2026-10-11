# Planificación de Proyectos — Especificación técnica de cambios

Documento de alcance técnico: qué se toca en backend, qué se toca en frontend.
Complementa [Implementation-plan.md](./Implementation-plan.md) (orden de trabajo).

Sigue el mismo formato que la especificación de [Proyecciones](../financial-projections/Technical-spec.md), y se apoya en la infraestructura que esa feature dejó montada: la integración con proyecciones de la sección 6 es un tercer sumando dentro del motor que ya existe.

## Objetivo

Nueva sección **Proyectos** que permite planificar un gasto grande antes de que ocurra: se define un presupuesto y una cuenta donde se junta la plata, se programa el ahorro que llena esa cuenta, se desglosa el gasto en ítems, y se registran las compras reales contra esos ítems. El presupuesto pendiente del proyecto aparece en la tabla de **Proyecciones** dentro de EGRESOS, en la categoría y subcategoría del proyecto.

---

## 1. Modelo conceptual

Tres piezas y una regla que las mantiene coherentes.

| Pieza | Qué es | Cómo se refleja en la contabilidad |
|---|---|---|
| **Proyecto** | Una meta de gasto futuro: presupuesto, categoría de gasto, cuenta donde se junta la plata, fechas | No es una transacción. Existe solo como plan hasta que se ejecuta |
| **Ahorro programado** | Un `TransactionTemplate` de tipo `SCHEDULE` **y de tipo transferencia**, desde una cuenta de origen hacia la cuenta del proyecto | Mueve dinero entre cuentas del usuario. **No es gasto** |
| **Ítem + compra** | Un renglón presupuestado del proyecto y el gasto real que lo materializa | Transacción de **gasto** normal, con la categoría del proyecto |

### Regla de oro: el gasto del proyecto se cuenta una sola vez

> **Un proyecto es gasto futuro. El ahorro es movimiento de dinero, no gasto.**
>
> - El ahorro programado es siempre una **transferencia**. Las transferencias ya están excluidas de las proyecciones en ambos lados ([transaction_projections.go:34-36](../../pkg/services/transaction_projections.go#L34-L36)), así que juntar plata nunca infla los egresos proyectados.
> - El proyecto aporta a la proyección **el presupuesto que todavía no se gastó**, en su mes de ejecución.
> - Las compras de los ítems son gastos reales y entran por el lado real de la proyección. Como el aporte del proyecto se descuenta exactamente en lo ya gastado, la suma de ambos nunca supera el presupuesto ni lo cuenta dos veces.

Esta regla es lo que hace que la sección 6 funcione. Cualquier cambio que rompa una de las tres viñetas produce doble conteo silencioso en la tabla de proyecciones.

### Decisión que se aparta del pedido original, y por qué

El pedido dice que un cambio de **categoría** del proyecto debe reflejarse en la transacción programada. **No se puede hacer, y no por una decisión de diseño.** La categoría del proyecto es de tipo *gasto*; el ahorro programado es una *transferencia*, y `isCategoryValid` rechaza una transferencia cuya categoría no sea de tipo transferencia ([transactions.go:3386-3390](../../pkg/services/transactions.go#L3386-L3390)). Propagar la categoría del proyecto a la plantilla la dejaría inválida y el cron dejaría de crear la transacción.

Lo que sí se propaga del proyecto a su plantilla de ahorro (sección 5): **cuenta destino, fecha límite del ahorro, importe recalculado desde el presupuesto y nombre**. Y la categoría del proyecto sí se propaga a donde importa de verdad y donde es válida: a **las transacciones de compra de sus ítems** y a la **fila de la proyección**.

La alternativa —permitir que el ahorro programado sea un gasto en vez de una transferencia, y así poder propagarle la categoría— rompe la regla de oro: ese gasto programado y el presupuesto pendiente del proyecto son la misma plata contada dos veces. Queda fuera de alcance de v1. Si se prefiere esa semántica hay que rediseñar la sección 6 completa, no agregar una opción.

---

## 2. Estado del proyecto

**El estado no se persiste.** Depende del reloj: guardarlo lo deja desactualizado sin que nadie toque el registro. Se guarda únicamente el flag manual `Completed`, y el resto se deriva.

```
Completed == true                              → COMPLETED   ("completado")
ExecutionTime == nil                           → DRAFT       ("borrador")
now < ExecutionTime                            → PLANNED     ("planificado")
ExecutionTime <= now <= DeadlineTime           → IN_PROGRESS ("en ejecución")
now > DeadlineTime                             → DELAYED     ("retrasado")
```

Casos borde, resueltos explícitamente:

| Situación | Estado | Por qué |
|---|---|---|
| `Completed` con fechas de cualquier tipo | `COMPLETED` | El estado manual gana sobre todo lo demás |
| Sin `ExecutionTime` pero con `DeadlineTime` | `DRAFT` | El pedido define borrador por ausencia de fecha de ejecución; la fecha límite sola no alcanza para planificar |
| Con `ExecutionTime` y sin `DeadlineTime`, y `now >= ExecutionTime` | `IN_PROGRESS` indefinidamente | Sin fecha límite no hay nada contra qué estar retrasado. No se inventa una |
| `DeadlineTime < ExecutionTime` | — | Rechazado en validación (`ErrProjectDeadlineBeforeExecutionTime`) |

**Dónde se calcula:** en el backend, en `(*Project).GetStatus(nowUnixTime, tz)`, y se devuelve ya resuelto en el DTO. Una segunda implementación en TypeScript se desincronizaría; el frontend trata `status` como campo de solo lectura.

**Comparación por día, no por instante.** Un proyecto con fecha de ejecución "hoy" debe verse `IN_PROGRESS` durante todo el día, no a partir de la hora exacta que quedó guardada. Las tres fechas se normalizan al inicio del día en la zona horaria del cliente antes de comparar.

**Consecuencia asumida:** una pestaña abierta que cruza la medianoche muestra el estado del día anterior. El store del frontend invalida por cambio de día, igual que hace el de proyecciones ([projection.ts:87-89](../../src/stores/projection.ts#L87-L89)).

---

## 3. Modelos nuevos

Cuatro tablas en `UserDataStore`, todas con el patrón del resto del proyecto: `Uid` + `Deleted` en los índices, tres timestamps al final, borrado lógico.

### 3.1 Tabla `project`

Archivo nuevo: `pkg/models/project.go`.

```go
// Project represents a planned expense project stored in database
type Project struct {
    ProjectId       int64  `xorm:"PK"`
    Uid             int64  `xorm:"INDEX(IDX_project_uid_deleted_order) NOT NULL"`
    Deleted         bool   `xorm:"INDEX(IDX_project_uid_deleted_order) NOT NULL"`
    Title           string `xorm:"VARCHAR(64) NOT NULL"`
    Description     string `xorm:"VARCHAR(255) NOT NULL"`
    CategoryId      int64  `xorm:"NOT NULL"`   // categoría de gasto de segundo nivel
    AccountId       int64  `xorm:"NOT NULL"`   // cuenta del proyecto
    ExecutionTime   *int64                     // opcional, unix time
    DeadlineTime    *int64                     // opcional, unix time
    BudgetAmount    int64  `xorm:"NOT NULL"`   // en la moneda de AccountId, unidad mínima
    Completed       bool   `xorm:"NOT NULL"`
    CompletedTime   *int64
    DisplayOrder    int32  `xorm:"INDEX(IDX_project_uid_deleted_order) NOT NULL"`
    Hidden          bool   `xorm:"NOT NULL"`
    CreatedUnixTime int64
    UpdatedUnixTime int64
    DeletedUnixTime int64
}
```

**No hay columna de costo real.** Se calcula al leer (sección 4). Una columna de caché se desincroniza en cuanto el usuario edita o borra la transacción de una compra desde la pantalla de transacciones, que es una ruta que este código no controla.

**Validaciones al crear/modificar:**

- `CategoryId` debe existir, no estar oculta, ser de tipo gasto y **de segundo nivel** — las mismas tres condiciones que `isCategoryValid` exige a una transacción ([transactions.go:3378-3390](../../pkg/services/transactions.go#L3378-L3390)). Si el proyecto admitiera una categoría que las compras no pueden usar, registrar una compra fallaría después de que el proyecto ya existe.
- `AccountId` debe existir, no estar oculta y no ser de tipo `ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS` — mismas condiciones que `doCreateTransaction` ([transactions.go:2700-2706](../../pkg/services/transactions.go#L2700-L2706)), por la misma razón.
- `BudgetAmount >= 0` y dentro de `MinimumTransactionAmount` / `MaximumTransactionAmount`.
- `DeadlineTime >= ExecutionTime` cuando ambas están presentes.

### 3.2 Tabla `project_item`

Archivo nuevo: `pkg/models/project_item.go`.

```go
// ProjectItem represents a budgeted line of a project stored in database
type ProjectItem struct {
    ItemId          int64  `xorm:"PK"`
    Uid             int64  `xorm:"INDEX(IDX_project_item_uid_deleted_project_id_order) NOT NULL"`
    Deleted         bool   `xorm:"INDEX(IDX_project_item_uid_deleted_project_id_order) NOT NULL"`
    ProjectId       int64  `xorm:"INDEX(IDX_project_item_uid_deleted_project_id_order) NOT NULL"`
    Name            string `xorm:"VARCHAR(64) NOT NULL"`
    Description     string `xorm:"VARCHAR(255) NOT NULL"`
    BudgetAmount    int64  `xorm:"NOT NULL"`   // en la moneda de la cuenta del proyecto
    Url             string `xorm:"VARCHAR(255) NOT NULL"`
    DisplayOrder    int32  `xorm:"INDEX(IDX_project_item_uid_deleted_project_id_order) NOT NULL"`
    CreatedUnixTime int64
    UpdatedUnixTime int64
    DeletedUnixTime int64
}
```

El ítem **no guarda costo real, vendedor ni fecha de compra**: esos datos pertenecen a la compra, no al renglón presupuestado, y viven en la tabla de relación.

La suma de los presupuestos de los ítems **no está obligada** a coincidir con el presupuesto del proyecto. La UI muestra la diferencia; el backend no la valida. Forzar la igualdad haría imposible empezar a desglosar un proyecto.

### 3.3 Tabla `project_item_transaction`

Archivo nuevo: `pkg/models/project_item_transaction.go`.

```go
// ProjectItemTransactionType represents the role a transaction plays in a project item purchase
type ProjectItemTransactionType byte

const (
    PROJECT_ITEM_TRANSACTION_TYPE_PURCHASE_EXPENSE ProjectItemTransactionType = 1
    PROJECT_ITEM_TRANSACTION_TYPE_FUNDING_TRANSFER ProjectItemTransactionType = 2
)

// ProjectItemTransaction represents project item and transaction relation stored in database
type ProjectItemTransaction struct {
    TransactionId   int64                      `xorm:"PK"`
    Uid             int64                      `xorm:"INDEX(IDX_project_item_transaction_uid_deleted_project_id) INDEX(IDX_project_item_transaction_uid_deleted_item_id) NOT NULL"`
    Deleted         bool                       `xorm:"INDEX(IDX_project_item_transaction_uid_deleted_project_id) INDEX(IDX_project_item_transaction_uid_deleted_item_id) NOT NULL"`
    ProjectId       int64                      `xorm:"INDEX(IDX_project_item_transaction_uid_deleted_project_id) NOT NULL"`
    ItemId          int64                      `xorm:"INDEX(IDX_project_item_transaction_uid_deleted_item_id) NOT NULL"`
    RelationType    ProjectItemTransactionType `xorm:"NOT NULL"`
    Vendor          string                     `xorm:"VARCHAR(64) NOT NULL"`
    CreatedUnixTime int64
    UpdatedUnixTime int64
    DeletedUnixTime int64
}
```

Tres decisiones:

- **La clave primaria es `TransactionId`, no un uuid nuevo.** Una transacción pertenece a lo sumo a un ítem, los ids de transacción ya son uuids globalmente únicos, y nunca se reutilizan — así que ni siquiera el borrado lógico puede colisionar. Ahorra un tipo de uuid, que es un recurso escaso (sección 3.5).
- **`ProjectId` está desnormalizado.** Permite calcular el costo real de todos los proyectos del usuario con una consulta sola, que es lo que necesita la proyección (sección 6).
- **`Vendor` vive acá y no en `Transaction`.** El modelo de transacción no tiene campo de vendedor y agregárselo obligaría a migrar una tabla que puede tener cientos de miles de filas por un dato que solo usa esta feature.

### 3.4 Tabla `project_scheduled_transaction`

Archivo nuevo: `pkg/models/project_scheduled_transaction.go`.

```go
// ProjectScheduledTransaction represents project and scheduled transaction template relation
type ProjectScheduledTransaction struct {
    TemplateId      int64 `xorm:"PK"`
    Uid             int64 `xorm:"INDEX(IDX_project_scheduled_transaction_uid_deleted_project_id) NOT NULL"`
    Deleted         bool  `xorm:"INDEX(IDX_project_scheduled_transaction_uid_deleted_project_id) NOT NULL"`
    ProjectId       int64 `xorm:"INDEX(IDX_project_scheduled_transaction_uid_deleted_project_id) NOT NULL"`
    CreatedUnixTime int64
    UpdatedUnixTime int64
    DeletedUnixTime int64
}
```

Misma decisión de clave primaria, por la misma razón. Un proyecto puede tener varias plantillas de ahorro (por ejemplo, un aporte quincenal fijo más un extra anual); una plantilla pertenece a un solo proyecto.

### 3.5 El espacio de tipos de uuid tiene 4 bits

⚠️ `UuidType` va de 0 a 15 ([uuid_type.go:3](../../pkg/uuid/uuid_type.go#L3), y el generador lo empaqueta en 4 bits en [internal_generator.go:64](../../pkg/uuid/internal_generator.go#L64)). Hay **11 tipos usados y 4 libres**: 12, 13, 14, 15.

Esta feature consume **dos**: `UUID_TYPE_PROJECT = 12` y `UUID_TYPE_PROJECT_ITEM = 13`. Quedan 14 y 15 para el resto de la vida del proyecto.

Ese es exactamente el motivo por el que las dos tablas de relación reutilizan como clave primaria el id de la entidad que enlazan. Un diseño con uuid propio en cada relación habría consumido los cuatro valores libres.

### 3.6 Registro de las tablas

Cuatro bloques nuevos al final de `updateAllDatabaseTablesStructure` ([cmd/database.go:45-185](../../cmd/database.go#L45-L185)), con el mismo `log.BootInfof` de confirmación que los demás:

```go
err = datastore.Container.UserDataStore.SyncStructs(new(models.Project))
// ... ProjectItem, ProjectItemTransaction, ProjectScheduledTransaction
```

---

## 4. Costo real: se calcula, no se guarda

**Definición: el costo real de un ítem es la plata que salió de la cuenta del proyecto.**

| Caso | Importe que cuenta |
|---|---|
| La compra se pagó con la cuenta del proyecto | `Amount` de la transacción de gasto |
| La compra se pagó con otra cuenta | `Amount` de la transferencia (el lado origen, que está en la moneda de la cuenta del proyecto) |

Con esta definición el costo real está **siempre en la moneda de la cuenta del proyecto**, sin conversiones, y es directamente comparable con `BudgetAmount`. Si en cambio se sumara el importe del gasto, un proyecto con compras en varias monedas daría un total sin sentido.

Traducido a consulta: por cada ítem se toma la fila de relación de tipo `FUNDING_TRANSFER` si existe, y si no la de tipo `PURCHASE_EXPENSE`; se une contra `transaction` filtrando `deleted = false`, y se suma `transaction.amount`.

**La transacción es la fuente de verdad, no la relación.** Si el usuario edita el importe desde la pantalla de transacciones, el costo real del proyecto lo refleja en la siguiente lectura. Si la borra, el `join` la descarta y el costo real se corrige solo; la fila de relación queda huérfana y sin efecto (`ErrProjectPurchaseTransactionMissing` solo se levanta si alguien intenta operar sobre ella).

Firma:

```go
// pkg/services/projects.go
// Devuelve el costo real por proyecto, opcionalmente acotado a compras anteriores a un instante.
// maxTransactionUnixTime <= 0 desactiva el corte.
func (s *ProjectService) GetSpentAmountsByProject(c core.Context, uid int64, projectIds []int64, maxTransactionUnixTime int64) (map[int64]int64, error)
```

El parámetro de corte no es decorativo: la proyección lo necesita para no descontar del presupuesto pendiente una compra con fecha futura que el lado real todavía no cuenta (sección 6).

---

## 5. Ahorro programado y propagación

### Creación

Endpoint dedicado, no reutiliza `/transaction/templates/add.json`: la plantilla y su fila de relación tienen que nacer en la **misma transacción de base de datos**, o un fallo intermedio deja una plantilla de ahorro que ningún proyecto reconoce y que el usuario no va a encontrar.

La plantilla se arma con estos campos fijados por el proyecto:

| Campo de la plantilla | Valor |
|---|---|
| `TemplateType` | `TRANSACTION_TEMPLATE_TYPE_SCHEDULE` |
| `Type` | `TRANSACTION_TYPE_TRANSFER` |
| `Name` | Título del proyecto |
| `RelatedAccountId` | `project.AccountId` (destino: la cuenta del proyecto) |
| `AccountId` | Cuenta de origen, elegida por el usuario |
| `CategoryId` | Categoría de **transferencia**, elegida por el usuario |
| `ScheduledEndTime` | `project.ExecutionTime` |
| `Amount` / `RelatedAccountAmount` | Ver abajo |

El resto —frecuencia, fecha de inicio, hora— lo elige el usuario.

**Precondición:** el proyecto debe tener `ExecutionTime`. Sin fecha de ejecución no hay hasta cuándo ahorrar. Un proyecto en borrador no puede programar ahorro (`ErrProjectExecutionTimeRequired`).

### Cálculo del importe por ocurrencia

```
ocurrencias   = len(GetScheduledOccurrences(template, now, project.ExecutionTime))
pendiente     = max(project.BudgetAmount - spent(project), 0)
Amount        = ceilDiv(pendiente, ocurrencias)
```

`GetScheduledOccurrences` es el motor que dejó la feature de proyecciones ([scheduled_frequency.go](../../pkg/services/scheduled_frequency.go)); usarlo garantiza que la cuenta de ocurrencias sea la misma que el cron va a ejecutar de verdad, incluidos los casos raros (día 31 en meses cortos no dispara, días negativos, zonas horarias).

Se redondea **hacia arriba**: es preferible juntar de más que llegar corto a la fecha de ejecución.

Si `ocurrencias == 0` —la frecuencia no dispara ni una vez antes de la fecha de ejecución— se rechaza con `ErrProjectSavingPlanHasNoOccurrences` en vez de crear una plantilla que nunca va a aportar nada.

**Monedas distintas.** Si la cuenta de origen y la del proyecto no comparten moneda, `RelatedAccountAmount` (lo que llega a la cuenta del proyecto) es el importe calculado arriba, y `Amount` (lo que sale de la cuenta de origen) lo aporta el usuario, con una conversión sugerida en la UI. Si comparten moneda, `Amount == RelatedAccountAmount`; `doCreateTransaction` rechaza lo contrario ([transactions.go:2708-2711](../../pkg/services/transactions.go#L2708-L2711)).

### Propagación al modificar el proyecto

Al modificar un proyecto, sobre **cada** plantilla enlazada, dentro de la misma transacción de base de datos:

| Cambio en el proyecto | Efecto en la plantilla |
|---|---|
| Título | `Name` |
| Cuenta del proyecto | `RelatedAccountId` |
| Fecha de ejecución | `ScheduledEndTime`, y **recálculo del importe** con la fórmula de arriba |
| Presupuesto | **Recálculo del importe** |
| Categoría | Nada — ver la advertencia de la sección 1 |
| Marcar como completado | Las plantillas enlazadas se **ocultan** (`Hidden = true`), no se borran |

Dos salvaguardas sobre el recálculo:

- **Es visible antes de confirmar.** El diálogo de edición del proyecto muestra "el ahorro programado pasará de X a Y" antes de guardar. Reescribir importes en silencio es la clase de comportamiento que hace que el usuario deje de confiar en la cifra.
- **Si el recálculo da `ocurrencias == 0`** (por ejemplo, se adelantó la fecha de ejecución a antes de la próxima ocurrencia), el importe **no se toca** y la respuesta incluye una advertencia. Es preferible un plan de ahorro desactualizado y visible que uno silenciosamente roto.

Completar un proyecto oculta sus plantillas en vez de borrarlas porque el ahorro ya ocurrido es historia contable; borrar la plantilla no borra las transacciones que creó, pero sí la explicación de por qué existen.

---

## 6. Integración con las proyecciones financieras

Punto de entrada: `GetProjectedCategoryAmountsByMonth` ([transaction_projections.go:41](../../pkg/services/transaction_projections.go#L41)), que hoy fusiona dos conjuntos. Pasa a fusionar **tres**:

```go
return mergeMonthlyTotalAmounts(
    mergeMonthlyTotalAmounts(removeTransferTotalAmounts(actualAmounts), projectedAmounts),
    projectAmounts,
), nil
```

`mergeMonthlyTotalAmounts` ya suma por `(yearMonth, categoryId, accountId)` ([transaction_projections.go:273](../../pkg/services/transaction_projections.go#L273)), así que si un proyecto y una plantilla programada comparten categoría y cuenta, caen en la misma celda de la tabla sin trabajo extra.

### El tercer sumando

```go
// pkg/services/project_projections.go (crear)
func (s *ProjectService) GetProjectedProjectAmountsByMonth(c core.Context, uid int64,
    startYear, startMonth, endYear, endMonth int32,
    clientTimezone *time.Location, currentUnixTime int64) (map[int32][]*models.TransactionTotalAmount, error)
```

Para cada proyecto no borrado del usuario:

1. **Se descarta** si está `COMPLETED` (su gasto real ya está en el lado real) o en `DRAFT` (sin fecha de ejecución no hay mes donde ponerlo).
2. **Importe** = `max(BudgetAmount - spent, 0)`, donde `spent` viene de `GetSpentAmountsByProject(..., maxTransactionUnixTime: currentUnixTime)`. El corte es el mismo instante que separa las dos mitades de la proyección: una compra con fecha futura no está contada en el lado real, así que tampoco puede descontarse del pendiente. Sin ese corte, una compra adelantada desaparece de la proyección.
3. Si el importe es `0` se omite: un proyecto totalmente ejecutado no aporta filas.
4. **Mes de imputación** = `max(mes de ExecutionTime, mes actual)`, calculado con la zona horaria del cliente igual que los otros dos conjuntos. El `max` es lo que hace que un proyecto **retrasado** siga apareciendo: su fecha ya pasó, pero la plata todavía no se gastó y sigue siendo un egreso pendiente. Sin el `max` caería en un mes del pasado y el motor lo descartaría por estar antes del corte.
5. Si el mes resultante queda fuera del rango pedido, se omite.
6. Se emite un `TransactionTotalAmount{Type: EXPENSE, CategoryId: project.CategoryId, AccountId: project.AccountId, Amount: pendiente}`.

`AccountId` va poblado porque el frontend convierte cada importe a la moneda por defecto usando la moneda de esa cuenta ([projection.ts:54-70](../../src/stores/projection.ts#L54-L70)); sin cuenta, el importe se descarta y la tabla marca `hasUnconvertedAmounts`.

### Qué hay que invalidar

El store de proyecciones cachea por período y por día ([projection.ts:117-122](../../src/stores/projection.ts#L117-L122)), y ya se invalida cuando cambia una plantilla programada ([transactionTemplate.ts:154](../../src/stores/transactionTemplate.ts#L154)). Hay que invalidarlo también al crear, modificar, completar o borrar un proyecto, y al registrar o anular una compra.

### Lo que el usuario ve

En la tabla de proyecciones, el proyecto no aparece como una fila propia: **se suma a la subcategoría de gasto que tiene asignada**, en su mes. Es lo pedido, y es lo correcto: desde el punto de vista del flujo de caja, un proyecto de "Renovar la cocina" en la subcategoría *Hogar → Reformas* es exactamente eso, un egreso de Reformas.

---

## 7. Registrar compra

El flujo más delicado de la feature: puede producir **dos transacciones y dos filas de relación**, y todo tiene que ser atómico.

### Entradas

| Dato | Origen |
|---|---|
| Ítem | El botón desde el que se abre |
| Costo real | Editable, por defecto el presupuesto del ítem |
| Vendedor | Opcional, texto libre |
| Fecha de compra | Por defecto hoy |
| Cuenta de compra | Por defecto la cuenta del proyecto |

### Qué se crea

**Caso A — la cuenta de compra es la cuenta del proyecto.** Una transacción de gasto:

```
Type       = TRANSACTION_DB_TYPE_EXPENSE
CategoryId = project.CategoryId
AccountId  = project.AccountId
Amount     = costo real
Comment    = nombre del ítem
```

Más una fila `ProjectItemTransaction` de tipo `PURCHASE_EXPENSE` con el vendedor.

**Caso B — la cuenta de compra es otra.** Primero la transferencia, después el gasto:

```
1) Transferencia   AccountId = project.AccountId      Amount = costo en moneda del proyecto
                   RelatedAccountId = cuenta de compra RelatedAccountAmount = costo en moneda de compra
2) Gasto           AccountId = cuenta de compra        Amount = costo en moneda de compra
                   CategoryId = project.CategoryId
```

Más **dos** filas de relación: `FUNDING_TRANSFER` y `PURCHASE_EXPENSE`, ambas apuntando al mismo `ItemId`.

Orden importante: primero la transferencia, para que la plata esté en la cuenta antes del gasto en el orden del libro mayor. Las colisiones de `transaction_time` dentro del mismo segundo las resuelve `doCreateTransaction` corriendo el instante ([transactions.go:2822](../../pkg/services/transactions.go#L2822)); la transferencia ocupa dos milisegundos (ida y vuelta) y el gasto cae en el siguiente libre.

La transferencia necesita una **categoría de transferencia**, que el usuario elige en el diálogo (no se puede derivar de la categoría de gasto del proyecto).

### Atomicidad: hay que abrir `CreateTransaction`

`CreateTransaction` abre su propia transacción de base de datos ([transactions.go:632-634](../../pkg/services/transactions.go#L632-L634)), así que llamarla dos veces seguidas **no es atómico**: un fallo en la segunda deja la transferencia sin el gasto que la justifica.

Se extrae de `CreateTransaction` un método que reciba la sesión ya abierta:

```go
// pkg/services/transactions.go
// Prepara los uuids y normaliza el instante de la transacción, y la inserta dentro de la sesión
// recibida. CreateTransaction pasa a ser este método envuelto en DoTransaction.
func (s *TransactionService) createTransactionInSession(c core.Context, database *datastore.Database,
    sess *xorm.Session, transaction *models.Transaction, tagIds []int64, pictureIds []int64) error
```

`CreateTransaction` queda como una envoltura de una línea sobre él, con el mismo comportamiento observable, y `ProjectService` —que vive en el mismo paquete `services`— lo usa para meter las dos transacciones y las dos relaciones en un solo `DoTransaction`.

### Anular una compra

Borra las transacciones que creó (con `DeleteTransaction`, que revierte los saldos de las cuentas) y marca las filas de relación como borradas, en la misma transacción de base de datos. El ítem vuelve a quedar disponible para "Registrar compra".

### Una compra por ítem en v1

El esquema soporta N compras por ítem, pero v1 permite **una activa**. Segundo intento → `ErrProjectItemAlreadyPurchased`. Compras parciales y devoluciones quedan para v2; abrirlas ahora obliga a decidir qué significa "costo real" cuando hay tres compras y una devolución, y no hay UI pedida para eso.

---

## 8. Endpoints

Todos bajo `/api/v1/projects/`, registrados en [cmd/webserver.go](../../cmd/webserver.go) después del bloque de exploradores (líneas 477-483), con el mismo `bindApi`.

```
GET  /v1/projects/list.json                    ?hidden=false
GET  /v1/projects/get.json                     ?id=<projectId>
POST /v1/projects/add.json
POST /v1/projects/modify.json
POST /v1/projects/complete.json                { id, completed }
POST /v1/projects/hide.json
POST /v1/projects/move.json
POST /v1/projects/delete.json

POST /v1/projects/items/add.json
POST /v1/projects/items/modify.json
POST /v1/projects/items/move.json
POST /v1/projects/items/delete.json

POST /v1/projects/items/purchase.json          registra la compra (sección 7)
POST /v1/projects/items/purchase/revoke.json   anula la compra

POST /v1/projects/scheduled/add.json           plantilla de ahorro + relación (sección 5)
POST /v1/projects/scheduled/delete.json        desenlaza y borra la plantilla
```

`list.json` devuelve los proyectos **con estado, costo real y resumen de ítems ya resueltos** — la pantalla de lista los necesita todos y pedirlos uno por uno sería N+1. `get.json` devuelve además los ítems con su compra y las plantillas de ahorro enlazadas.

### DTO principal

```go
type ProjectInfoResponse struct {
    Id                  int64                          `json:"id,string"`
    Title               string                         `json:"title"`
    Description         string                         `json:"description,omitempty"`
    CategoryId          int64                          `json:"categoryId,string"`
    AccountId           int64                          `json:"accountId,string"`
    ExecutionTime       *int64                         `json:"executionTime,omitempty"`
    DeadlineTime        *int64                         `json:"deadlineTime,omitempty"`
    BudgetAmount        int64                          `json:"budgetAmount"`
    ActualAmount        int64                          `json:"actualAmount"`    // calculado
    Status              ProjectStatus                  `json:"status"`          // calculado
    DisplayOrder        int32                          `json:"displayOrder"`
    Hidden              bool                           `json:"hidden"`
    Items               []*ProjectItemInfoResponse     `json:"items,omitempty"`
    ScheduledTemplates  []*TransactionTemplateInfoResponse `json:"scheduledTemplates,omitempty"`
}
```

`ActualAmount` y `Status` son **solo de salida**: el backend los calcula y los ignora si llegan en una petición.

### Errores nuevos

`NormalSubcategoryProject = 21` en [pkg/errs/error.go](../../pkg/errs/error.go#L46) y un [pkg/errs/project.go](../../pkg/errs/project.go) nuevo con el patrón de [errs/explorer.go](../../pkg/errs/explorer.go):

```
ErrProjectIdInvalid                  ErrProjectItemIdInvalid
ErrProjectNotFound                   ErrProjectItemNotFound
ErrProjectCategoryInvalid            ErrProjectItemAlreadyPurchased
ErrProjectAccountInvalid             ErrProjectPurchaseNotFound
ErrProjectDeadlineBeforeExecutionTime ErrProjectPurchaseTransactionMissing
ErrProjectExecutionTimeRequired      ErrProjectSavingPlanHasNoOccurrences
ErrProjectAlreadyCompleted
```

### Borrado de datos del usuario

`ClearAllDataHandler` ([data_managements.go:170](../../pkg/api/data_managements.go#L170)) borra plantillas, transacciones, categorías, tags, tipos de cambio y exploraciones. Hay que agregar `DeleteAllProjects`, y **antes** de `DeleteAllTransactions`, para que las relaciones se borren mientras las transacciones que referencian todavía existen.

---

# Cambios en Frontend

## 9. Ventanas a implementar

### 9.1 Lista de proyectos — `/projects/list`

```
Proyectos                                              [+ Nuevo proyecto]
─────────────────────────────────────────────────────────────────────────
 Renovar la cocina           Hogar › Reformas      ▓▓▓▓▓▓░░░░  60%
 EN EJECUCIÓN · vence 30 sep        $6.000 / $10.000        Ahorros
─────────────────────────────────────────────────────────────────────────
 Viaje a Japón               Ocio › Viajes         ▓▓░░░░░░░░  20%
 PLANIFICADO · 15 mar 2027          $1.600 / $8.000         Caja de ahorro
─────────────────────────────────────────────────────────────────────────
 Cambiar la moto             Transporte › Vehículo ░░░░░░░░░░   0%
 BORRADOR · sin fecha                   $0 / $4.500          Ahorros
─────────────────────────────────────────────────────────────────────────
```

Barra de progreso = costo real / presupuesto, con color por estado. Ordenable por `DisplayOrder` con arrastrar y soltar, igual que la lista de plantillas ([templates/ListPage.vue](../../src/views/desktop/templates/ListPage.vue)). Filtro por estado y conmutador de ocultos.

### 9.2 Detalle del proyecto — `/projects/:id`

```
◄ Renovar la cocina                          [Editar] [Completar] [⋮]
  EN EJECUCIÓN · ejecución 1 sep · límite 30 sep · Hogar › Reformas
─────────────────────────────────────────────────────────────────────────
  Presupuesto  $10.000     Costo real  $6.000     Pendiente  $4.000
  Cuenta del proyecto: Ahorros                 ▓▓▓▓▓▓░░░░  60%
─────────────────────────────────────────────────────────────────────────
  AHORRO PROGRAMADO                                   [+ Programar]
    Mensual, día 5 · desde Cuenta corriente · hasta 1 sep
    $1.000 por ocurrencia · 4 restantes                    [Editar] [🗑]
─────────────────────────────────────────────────────────────────────────
  ÍTEMS                                                    [+ Ítem]
    Horno empotrable          $3.000    ✔ $2.900  Frávega  12 ago  [⋮]
    Mesada de granito         $4.500    ✔ $3.100  Mármoles 20 ago  [⋮]
    Grifería                  $1.200    [Registrar compra]          [⋮]
    Instalación                 $800    [Registrar compra]          [⋮]
  ─────────────────────────────────────────────────────────────────────
    Presupuestado en ítems    $9.500    Comprado  $6.000
```

La fila de cierre expone la diferencia entre el presupuesto del proyecto y la suma de los ítems, que el backend no valida a propósito (sección 3.2).

Un ítem ya comprado muestra costo real, vendedor y fecha, con enlace a la transacción; su menú ofrece **Anular compra**.

### 9.3 Diálogos

| Diálogo | Contenido |
|---|---|
| **Nuevo/Editar proyecto** | Título, descripción, categoría (`two-column-select` filtrado a gasto), cuenta, fechas, presupuesto. Al editar con ahorro enlazado, muestra el recálculo del importe antes de guardar |
| **Nuevo/Editar ítem** | Nombre, descripción, presupuesto, URL |
| **Registrar compra** | Costo real (precargado con el presupuesto del ítem), vendedor, fecha, cuenta de compra. Si la cuenta difiere de la del proyecto, despliega el bloque de transferencia: categoría de transferencia e importe convertido |
| **Programar ahorro** | Cuenta de origen, categoría de transferencia, frecuencia (`ScheduleFrequencySelect`), fecha de inicio, hora. Muestra en vivo el importe calculado y cuántas ocurrencias quedan hasta la fecha de ejecución |
| **Completar proyecto** | Confirmación, avisando que las plantillas de ahorro enlazadas se ocultarán |

**El diálogo de ahorro es propio, no reutiliza [EditDialog.vue](../../src/views/desktop/transactions/list/dialogs/EditDialog.vue).** Ese componente pesa más de 1.100 líneas y está cableado al store de plantillas, que guardaría contra `/transaction/templates/add.json` y rompería la atomicidad plantilla+relación de la sección 5. El diálogo propio expone cinco campos —el resto los fija el proyecto— y reutiliza los componentes de bajo nivel: `ScheduleFrequencySelect`, `AmountInput`, `DateTimeSelect`, `two-column-select`.

## 10. Archivos a modificar / crear

### Backend

| Archivo | Ubicación | Acción |
|---|---|---|
| `project.go` (modelos) | `pkg/models/` | **Crear**: `Project`, `ProjectStatus`, requests, `ProjectInfoResponse` |
| `project_item.go` | `pkg/models/` | **Crear** |
| `project_item_transaction.go` | `pkg/models/` | **Crear** |
| `project_scheduled_transaction.go` | `pkg/models/` | **Crear** |
| `projects.go` (servicio) | `pkg/services/` | **Crear**: CRUD, `GetSpentAmountsByProject`, propagación |
| `project_purchases.go` | `pkg/services/` | **Crear**: orquestación de compra y anulación (sección 7) |
| `project_projections.go` | `pkg/services/` | **Crear**: `GetProjectedProjectAmountsByMonth` |
| `transactions.go` | [pkg/services/transactions.go:567](../../pkg/services/transactions.go#L567) | Modificar: extraer `createTransactionInSession` |
| `transaction_projections.go` | [pkg/services/transaction_projections.go:41](../../pkg/services/transaction_projections.go#L41) | Modificar: tercer sumando en la fusión |
| `projects.go` (api) | `pkg/api/` | **Crear**: `ProjectsApi` |
| `data_managements.go` | [pkg/api/data_managements.go:170](../../pkg/api/data_managements.go#L170) | Modificar: `DeleteAllProjects` antes de las transacciones |
| `error.go` | [pkg/errs/error.go:46](../../pkg/errs/error.go#L46) | Modificar: `NormalSubcategoryProject = 21` |
| `project.go` (errores) | `pkg/errs/` | **Crear** |
| `uuid_type.go` | [pkg/uuid/uuid_type.go](../../pkg/uuid/uuid_type.go) | Modificar: dos tipos nuevos (sección 3.5) |
| `database.go` | [cmd/database.go:45](../../cmd/database.go#L45) | Modificar: cuatro `SyncStructs` |
| `webserver.go` | [cmd/webserver.go:483](../../cmd/webserver.go#L483) | Modificar: registrar las rutas |
| `projects_test.go` | `pkg/services/` | **Crear**: estado derivado, costo real, cálculo del ahorro |
| `project_projections_test.go` | `pkg/services/` | **Crear**: imputación por mes, corte, no doble conteo |

### Frontend

| Archivo | Ubicación | Acción |
|---|---|---|
| `project.ts` (core) | `src/core/` | **Crear**: `ProjectStatus` y su info de presentación |
| `project.ts` (modelos) | `src/models/` | **Crear**: `Project`, `ProjectItem`, requests y respuestas |
| `project.ts` (lib) | `src/lib/` | **Crear**: funciones puras (progreso, pendiente, resumen de ítems) |
| `services.ts` | [src/lib/services.ts:854](../../src/lib/services.ts#L854) | Modificar: los 15 endpoints, mismo esqueleto que los de exploradores |
| `project.ts` (store) | `src/stores/` | **Crear** |
| `projection.ts` (store) | [src/stores/projection.ts](../../src/stores/projection.ts) | Modificar: invalidar cuando cambia un proyecto |
| `ProjectProgressBar.vue` | `src/components/desktop/` | **Crear**: barra + porcentaje, reutilizada en lista y detalle |
| `ListPage.vue` | `src/views/desktop/projects/` | **Crear** |
| `DetailPage.vue` | `src/views/desktop/projects/` | **Crear** |
| `dialogs/*.vue` | `src/views/desktop/projects/dialogs/` | **Crear**: los cinco diálogos de 9.3 |
| `desktop.ts` (router) | [src/router/desktop.ts:160](../../src/router/desktop.ts#L160) | Modificar: `/projects/list` y `/projects/:id` |
| `MainLayout.vue` | [src/views/desktop/MainLayout.vue:49-54](../../src/views/desktop/MainLayout.vue#L49-L54) | Modificar: `nav-link` después de "Projections" |
| `desktop-main.ts` | [src/desktop-main.ts:594](../../src/desktop-main.ts#L594) | Modificar: registrar `ProjectProgressBar` |
| `*.json` (i18n) | [src/locales/](../../src/locales/) | Modificar: **20 archivos** de idioma |

### Store `src/stores/project.ts`

Patrón de [src/stores/explorer.ts](../../src/stores/explorer.ts):

- **Estado**: lista de proyectos, proyecto abierto con sus ítems y plantillas, flags de invalidación.
- **Acciones**: `loadProjects`, `loadProject`, `saveProject`, `completeProject`, `deleteProject`, `saveProjectItem`, `deleteProjectItem`, `registerPurchase`, `revokePurchase`, `scheduleSaving`, `deleteScheduledSaving`.
- **Invalidación cruzada**: toda acción que altere presupuesto, categoría, fecha de ejecución, estado o compras llama a `projectionsStore.updateProjectionInvalidState(true)`. Registrar una compra además invalida el store de transacciones y el de cuentas, porque creó transacciones y movió saldos.
- **`status` y `actualAmount` no se recalculan en el cliente**: se toman de la respuesta (sección 2). El único cálculo local es el porcentaje de progreso.

### i18n

Claves nuevas en los 20 archivos de `src/locales/`, entre otras:

`"Projects"`, `"New Project"`, `"Project Title"`, `"Project Account"`, `"Execution Date"`, `"Deadline"`, `"Budget"`, `"Actual Cost"`, `"Remaining"`, `"Draft"`, `"Planned"`, `"In Progress"`, `"Delayed"`, `"Completed"`, `"Complete Project"`, `"Project Items"`, `"New Item"`, `"Register Purchase"`, `"Revoke Purchase"`, `"Vendor"`, `"Purchase Date"`, `"Purchase Account"`, `"Schedule Saving"`, `"Saving Plan"`, `"Per occurrence"`, `"Occurrences left"`, `"Budgeted in items"`

`"Category"`, `"Account"`, `"Description"`, `"Amount"` y `"Transfer"` **ya existen** — reutilizarlas.

---

## 11. Limitaciones conocidas de v1

- **La propagación de categoría al ahorro programado no existe**, por incompatibilidad de tipos de categoría (sección 1). Es la única parte del pedido original que no se implementa tal cual.
- **Una compra activa por ítem.** Compras parciales, cuotas y devoluciones quedan para v2 (sección 7).
- **El recálculo del ahorro no se dispara solo con el paso del tiempo.** Si el usuario deja de aportar durante tres meses, el importe por ocurrencia no se ajusta hasta que edite el proyecto. La pantalla de detalle muestra ocurrencias restantes y pendiente, así que la desviación es visible, pero no se corrige sola.
- **Un proyecto con fecha de ejecución muy lejana** cae fuera del rango de la proyección (tope de 60 meses, [transaction_projections.go:17](../../pkg/services/transaction_projections.go#L17)) y simplemente no aparece.
- **Presupuesto y costo real en la moneda de la cuenta del proyecto.** Cambiar la cuenta del proyecto a una de otra moneda **no reexpresa** los importes ya cargados. La UI advierte; el backend no convierte.
- **Relaciones huérfanas.** Borrar desde la pantalla de transacciones una transacción que pertenece a una compra deja la fila de relación sin efecto. El costo real se corrige solo, pero el ítem queda marcado como comprado sin transacción detrás; se muestra como "compra incompleta" con opción de anular.

## 12. Fuera de alcance de v1

- **Mobile.** No se agregan páginas en `src/views/mobile/`, igual que Proyecciones e Insights Explorer.
- **Ahorro programado de tipo gasto.** Rompe la regla de oro (sección 1).
- **Exportación e importación de proyectos** en las herramientas de datos.
- **Adjuntar imágenes** a ítems o compras, aunque el modelo de transacción lo soportaría.
- **Compartir o colaborar** en un proyecto entre usuarios.
- **Fila propia del proyecto en la tabla de proyecciones.** El importe se suma a su subcategoría, sin desglose (sección 6).
