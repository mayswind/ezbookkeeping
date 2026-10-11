#!/usr/bin/env python3
"""Genera datos de prueba para ezBookkeeping directamente en la base de datos.

Emite SQL por stdout. Aplicar con:

    python3 scripts/seed-dev-data.py | docker exec -i ezbookkeeping-db \
        psql -U sebastianmaldonado -d financial_app_db -v ON_ERROR_STOP=1 -1

Los datos respetan las invariantes que mantiene pkg/services/transactions.go:

  * Los ids siguen el formato del generador interno (pkg/uuid/internal_generator.go):
    unixTime(32b) << 31 | uuidType(4b) << 27 | serverId(8b) << 19 | seq(19b).
    Usamos un serverId propio para no colisionar nunca con los ids del servidor.
  * transaction_time = unixTime * 1000 + indice, y es unico por usuario
    (UQE_transaction_uid_time).
  * Una transferencia son dos filas: TRANSFER_OUT en T y TRANSFER_IN en T+1,
    apuntandose mutuamente con related_id.
  * account.balance es la suma de los deltas de sus transacciones:
    MODIFY_BALANCE -> related_account_amount, INCOME/TRANSFER_IN -> +amount,
    EXPENSE/TRANSFER_OUT -> -amount.
  * Ninguna transaccion de una cuenta es anterior a su MODIFY_BALANCE (saldo inicial).

Los importes se expresan en centavos porque COP se define con fraction: 2
(src/consts/currency.ts), igual que el resto de la app.
"""

import datetime as dt
import random
import sys

# ----------------------------------------------------------------------------
# Parametros
# ----------------------------------------------------------------------------

UID = 3837566490429620224          # usuario "sebas"
CURRENCY = 'COP'
TZ_OFFSET_MIN = -300               # America/Bogota, en minutos
SEED_START = dt.date(2026, 1, 1)
SEED_END = dt.date(2026, 10, 10)

# serverId reservado para los datos sembrados: el servidor corre con server_id=0
# (conf/ezbookkeeping.ini), asi que estos ids no pueden colisionar con los suyos.
SEED_SERVER_ID = 170

UUID_TYPE_ACCOUNT = 2
UUID_TYPE_TRANSACTION = 3
UUID_TYPE_TEMPLATE = 7

TYPE_MODIFY_BALANCE = 1
TYPE_INCOME = 2
TYPE_EXPENSE = 3
TYPE_TRANSFER_OUT = 4
TYPE_TRANSFER_IN = 5

# Categorias existentes del usuario (transaction_category)
CAT = {
    'salario':        3837566490832273425,
    'bonificacion':   3837566490832273426,
    'extras':         3837566490832273427,
    'freelance':      3837566490832273428,
    'inversiones':    3837566490832273429,
    'intereses':      3837566490832273431,
    'comida':         3837566490832273436,
    'bebida':         3837566490832273437,
    'ropa':           3837566490832273439,
    'hogar':          3837566490832273443,
    'electronica':    3837566490832273444,
    'servicios':      3837566490832273447,
    'arriendo':       3837566490832273448,
    'transporte':     3837566490832273449,
    'taxi':           3837566490832273450,
    'carro':          3837566490832273451,
    'telefono':       3837566490832273454,
    'internet':       3837566490832273455,
    'fitness':        3837566490832273457,
    'cine':           3837566490832273459,
    'suscripciones':  3837566490832273461,
    'medicamentos':   3837566490832273470,
    'seguros':        3837566490832273474,
    'transferencia':  3837566490832273478,
    'pago_tc':        3837566490832273479,
    'deposito':       3837566490832273480,
}

PESO = 100  # centavos por peso


def pesos(value):
    return int(round(value * PESO))


# ----------------------------------------------------------------------------
# Generacion de ids
# ----------------------------------------------------------------------------

_now = int(dt.datetime.now().timestamp())
_seq = {}


def new_id(uuid_type):
    """Reproduce assembleUuid() de pkg/uuid/internal_generator.go."""
    seq = _seq.get(uuid_type, 0)
    _seq[uuid_type] = seq + 1

    if seq > (1 << 19) - 1:
        raise SystemExit('demasiados ids para el tipo %d' % uuid_type)

    return (_now << 31) | (uuid_type << 27) | (SEED_SERVER_ID << 19) | seq


# ----------------------------------------------------------------------------
# Tiempo de transaccion
# ----------------------------------------------------------------------------

TZ = dt.timezone(dt.timedelta(minutes=TZ_OFFSET_MIN))
_used_seconds = set()


def local(year, month, day, hour=12, minute=0):
    return dt.datetime(year, month, day, hour, minute, tzinfo=TZ)


def transaction_time(moment, slots=1):
    """Devuelve el transaction_time de un instante, reservando `slots` valores.

    transaction_time es unixTime*1000 + indice y debe ser unico por usuario; una
    transferencia necesita dos valores consecutivos dentro del mismo segundo.
    """
    unix = int(moment.timestamp())

    while unix in _used_seconds:
        unix += 1

    _used_seconds.add(unix)

    if slots > 1000:
        raise SystemExit('no caben %d transacciones en un segundo' % slots)

    return unix * 1000


# ----------------------------------------------------------------------------
# Acumuladores
# ----------------------------------------------------------------------------

accounts = []       # dicts listos para INSERT
balances = {}       # account_id -> saldo acumulado en centavos
transactions = []   # dicts listos para INSERT
templates = []


def add_account(name, category, icon, color, opening, opened_at, display_order,
                statement_date=None, comment=''):
    account_id = new_id(UUID_TYPE_ACCOUNT)
    extend = '{"lastReconciledTime":null,"creditCardStatementDate":%s}' % (
        'null' if statement_date is None else statement_date
    )

    accounts.append({
        'account_id': account_id,
        'category': category,
        'icon': icon,
        'color': color,
        'name': name,
        'display_order': display_order,
        'comment': comment,
        'extend': extend,
    })
    balances[account_id] = 0

    # El saldo inicial es una transaccion MODIFY_BALANCE: amount es el saldo
    # fijado y related_account_amount el delta aplicado a la cuenta (que partia
    # de cero), tal como lo calcula doCreateTransaction.
    emit(TYPE_MODIFY_BALANCE, 0, account_id, opened_at, opening,
         related_account_id=account_id, related_account_amount=opening,
         comment='Saldo inicial')

    return account_id


def emit(txn_type, category_id, account_id, moment, amount,
         related_account_id=0, related_account_amount=0, related_id=0,
         comment='', transaction_id=None, time_value=None):
    txn = {
        'transaction_id': transaction_id or new_id(UUID_TYPE_TRANSACTION),
        'type': txn_type,
        'category_id': category_id,
        'account_id': account_id,
        'transaction_time': time_value if time_value is not None else transaction_time(moment),
        'amount': amount,
        'related_id': related_id,
        'related_account_id': related_account_id,
        'related_account_amount': related_account_amount,
        'comment': comment,
    }
    transactions.append(txn)

    if txn_type == TYPE_MODIFY_BALANCE:
        balances[account_id] += related_account_amount
    elif txn_type in (TYPE_INCOME, TYPE_TRANSFER_IN):
        balances[account_id] += amount
    else:
        balances[account_id] -= amount

    return txn


def income(category, account_id, moment, amount, comment=''):
    emit(TYPE_INCOME, CAT[category], account_id, moment, amount, comment=comment)


def expense(category, account_id, moment, amount, comment=''):
    emit(TYPE_EXPENSE, CAT[category], account_id, moment, amount, comment=comment)


def transfer(category, from_account, to_account, moment, amount, comment=''):
    """Inserta el par TRANSFER_OUT / TRANSFER_IN que espera el servidor."""
    out_id = new_id(UUID_TYPE_TRANSACTION)
    in_id = new_id(UUID_TYPE_TRANSACTION)
    out_time = transaction_time(moment, slots=2)

    emit(TYPE_TRANSFER_OUT, CAT[category], from_account, moment, amount,
         related_account_id=to_account, related_account_amount=amount,
         related_id=in_id, comment=comment, transaction_id=out_id,
         time_value=out_time)

    # GetRelatedTransferTransaction coloca la fila entrante en T+1 y cruza los
    # importes: el amount de cada fila es el de su propia cuenta.
    emit(TYPE_TRANSFER_IN, CAT[category], to_account, moment, amount,
         related_account_id=from_account, related_account_amount=amount,
         related_id=out_id, comment=comment, transaction_id=in_id,
         time_value=out_time + 1)


def add_template(name, txn_type, category, account_id, amount, days,
                 hour, display_order, comment=''):
    """Plantilla programada mensual, para la pagina de proyecciones.

    scheduled_at son minutos desde la medianoche UTC, por eso se compensa el
    offset del usuario (ver GetScheduledOccurrences).
    """
    scheduled_at = hour * 60 - TZ_OFFSET_MIN

    templates.append({
        'template_id': new_id(UUID_TYPE_TEMPLATE),
        'name': name,
        'type': txn_type,
        'category_id': CAT[category],
        'account_id': account_id,
        'scheduled_frequency': ','.join(str(d) for d in days),
        'scheduled_start_time': int(local(SEED_START.year, SEED_START.month, SEED_START.day, 0, 0).timestamp()),
        'scheduled_end_time': int(local(2028, 12, 31, 23, 59).timestamp()),
        'scheduled_at': scheduled_at,
        'amount': amount,
        'comment': comment,
        'display_order': display_order,
    })


# ----------------------------------------------------------------------------
# Cuentas
# ----------------------------------------------------------------------------

random.seed(20261010)  # reproducible

opened_at = local(SEED_START.year, SEED_START.month, SEED_START.day, 8, 0)

NOMINA = add_account('Cuenta Nomina', 2, 110, '1976d2', pesos(4_500_000), opened_at, 2,
                     comment='Cuenta corriente donde se recibe el salario')
EFECTIVO = add_account('Efectivo', 1, 1, '4caf50', pesos(350_000), opened_at, 3)
AHORRO = add_account('Ahorro Vivienda', 8, 30, 'ff9800', pesos(12_000_000), opened_at, 4,
                     comment='Meta: cuota inicial')
# Las cuentas de pasivo guardan el saldo en negativo: la UI muestra -balance
# como saldo pendiente (getAccountBalance en src/stores/account.ts).
VISA = add_account('Tarjeta Visa', 3, 100, 'e91e63', pesos(-1_200_000), opened_at, 5,
                   statement_date=15)
FONDO = add_account('Fondo de Inversion', 7, 800, '9c27b0', pesos(6_000_000), opened_at, 6)


# ----------------------------------------------------------------------------
# Movimientos mes a mes
# ----------------------------------------------------------------------------

def months(start, end):
    year, month = start.year, start.month

    while (year, month) <= (end.year, end.month):
        yield year, month
        year, month = (year + 1, 1) if month == 12 else (year, month + 1)


def jitter(base, spread):
    return pesos(base + random.randint(-spread, spread))


for year, month in months(SEED_START, SEED_END):
    last_day = SEED_END.day if (year, month) == (SEED_END.year, SEED_END.month) else 31

    def on(day, hour=12, minute=0):
        """Instante del dia indicado, o None si cae fuera del periodo sembrado."""
        try:
            moment = local(year, month, day, hour, minute)
        except ValueError:
            return None

        return moment if day <= last_day else None

    # --- Ingresos ---
    when = on(1, 8, 0)
    if when:
        income('salario', NOMINA, when, pesos(4_800_000), 'Salario mensual')

    if month in (6, 12):
        when = on(15, 8, 30)
        if when:
            income('bonificacion', NOMINA, when, pesos(2_400_000), 'Prima semestral')

    if month in (3, 6, 9):
        when = on(28, 9, 0)
        if when:
            income('inversiones', FONDO, when, jitter(180_000, 40_000),
                   'Rendimientos del trimestre')

    if month % 2 == 0:
        when = on(18, 19, 0)
        if when:
            income('freelance', NOMINA, when, jitter(900_000, 250_000),
                   'Proyecto freelance')

    # --- Gastos fijos ---
    for day, hour, category, account, base, spread, note in (
        (3, 9, 'arriendo', NOMINA, 1_450_000, 0, 'Arriendo apartamento'),
        (5, 10, 'servicios', NOMINA, 280_000, 45_000, 'Energia, agua y gas'),
        (6, 10, 'internet', NOMINA, 95_000, 0, 'Fibra hogar'),
        (6, 11, 'telefono', NOMINA, 65_000, 8_000, 'Plan movil'),
        (8, 20, 'suscripciones', VISA, 48_000, 0, 'Streaming y musica'),
        (9, 11, 'seguros', NOMINA, 135_000, 0, 'Seguro de vida'),
        (22, 18, 'fitness', NOMINA, 120_000, 0, 'Mensualidad gimnasio'),
    ):
        when = on(day, hour)
        if when:
            expense(category, account, when, jitter(base, spread), note)

    # --- Gastos variables ---
    for i, day in enumerate((7, 14, 21, 28)):
        when = on(day, 17, 30)
        if when:
            # El mercado se paga alternando cuenta y tarjeta.
            account = VISA if i % 2 else NOMINA
            expense('comida', account, when, jitter(250_000, 70_000), 'Mercado semanal')

    for day in (4, 11, 18, 25):
        when = on(day, 7, 15)
        if when:
            expense('transporte', EFECTIVO, when, jitter(42_000, 12_000),
                    'Transporte publico')

    when = on(12, 20, 30)
    if when:
        expense('bebida', VISA, when, jitter(95_000, 35_000), 'Salida con amigos')

    when = on(19, 13, 0)
    if when:
        expense('taxi', EFECTIVO, when, jitter(28_000, 10_000), 'Taxi')

    if month % 3 == 0:
        when = on(23, 16, 0)
        if when:
            expense('carro', NOMINA, when, jitter(320_000, 90_000),
                    'Mantenimiento del carro')

    if month % 2:
        when = on(16, 15, 0)
        if when:
            expense('medicamentos', EFECTIVO, when, jitter(75_000, 25_000), 'Farmacia')

    if month in (2, 5, 8):
        when = on(24, 19, 30)
        if when:
            expense('ropa', VISA, when, jitter(260_000, 80_000), 'Ropa de temporada')

    if month in (4, 7):
        when = on(26, 20, 0)
        if when:
            expense('electronica', VISA, when, jitter(850_000, 150_000),
                    'Equipo de computo')

    when = on(27, 21, 0)
    if when:
        expense('cine', EFECTIVO, when, jitter(45_000, 15_000), 'Cine')

    if month in (1, 5, 9):
        when = on(13, 11, 0)
        if when:
            expense('hogar', NOMINA, when, jitter(180_000, 60_000), 'Articulos del hogar')

    # --- Transferencias ---
    when = on(2, 9, 30)
    if when:
        transfer('transferencia', NOMINA, AHORRO, when, pesos(800_000),
                 'Ahorro programado')

    when = on(15, 12, 0)
    if when:
        transfer('deposito', NOMINA, EFECTIVO, when, pesos(300_000),
                 'Retiro para gastos del mes')

    when = on(20, 10, 0)
    if when:
        transfer('pago_tc', NOMINA, VISA, when, pesos(900_000),
                 'Pago tarjeta de credito')

    when = on(25, 10, 30)
    if when:
        transfer('transferencia', NOMINA, FONDO, when, pesos(500_000),
                 'Aporte al fondo')


# ----------------------------------------------------------------------------
# Plantillas programadas (alimentan la pagina de proyecciones)
# ----------------------------------------------------------------------------

add_template('Salario mensual', TYPE_INCOME, 'salario', NOMINA,
             pesos(4_800_000), [1], 8, 3, 'Nomina')
add_template('Arriendo', TYPE_EXPENSE, 'arriendo', NOMINA,
             pesos(1_450_000), [3], 9, 4, 'Arriendo apartamento')
add_template('Servicios publicos', TYPE_EXPENSE, 'servicios', NOMINA,
             pesos(280_000), [5], 10, 5, 'Energia, agua y gas')
add_template('Internet y telefonia', TYPE_EXPENSE, 'internet', NOMINA,
             pesos(160_000), [6], 10, 6, 'Fibra y plan movil')
add_template('Suscripciones', TYPE_EXPENSE, 'suscripciones', VISA,
             pesos(48_000), [8], 20, 7, 'Streaming y musica')
add_template('Mercado quincenal', TYPE_EXPENSE, 'comida', NOMINA,
             pesos(250_000), [7, 21], 17, 8, 'Mercado')
add_template('Gimnasio', TYPE_EXPENSE, 'fitness', NOMINA,
             pesos(120_000), [22], 18, 9, 'Mensualidad gimnasio')


# ----------------------------------------------------------------------------
# Salida SQL
# ----------------------------------------------------------------------------

def sql_str(value):
    return "'" + value.replace("'", "''") + "'"


out = sys.stdout.write

out('-- Datos de prueba generados por scripts/seed-dev-data.py\n')
out('-- usuario %d, periodo %s .. %s\n\n' % (UID, SEED_START, SEED_END))

out('-- Borra una ejecucion anterior de este script (identificada por el serverId %d)\n'
    % SEED_SERVER_ID)
out("DELETE FROM transaction WHERE uid = %d AND ((transaction_id >> 19) & 255) = %d;\n"
    % (UID, SEED_SERVER_ID))
out("DELETE FROM transaction_template WHERE uid = %d AND ((template_id >> 19) & 255) = %d;\n"
    % (UID, SEED_SERVER_ID))
out("DELETE FROM account WHERE uid = %d AND ((account_id >> 19) & 255) = %d;\n\n"
    % (UID, SEED_SERVER_ID))

out('-- Cuentas (%d)\n' % len(accounts))
for account in accounts:
    out(
        'INSERT INTO account (account_id, uid, deleted, category, type, parent_account_id, '
        'name, display_order, icon, icon_type, color, currency, balance, comment, extend, '
        'hidden, created_unix_time, updated_unix_time, deleted_unix_time) VALUES '
        '(%d, %d, false, %d, 1, 0, %s, %d, %d, 0, %s, %s, %d, %s, '
        'decode(%s, %s), false, %d, %d, 0);\n' % (
            account['account_id'], UID, account['category'], sql_str(account['name']),
            account['display_order'], account['icon'], sql_str(account['color']),
            sql_str(CURRENCY), balances[account['account_id']], sql_str(account['comment']),
            sql_str(account['extend'].encode().hex()), sql_str('hex'), _now, _now,
        )
    )

out('\n-- Transacciones (%d)\n' % len(transactions))
for txn in transactions:
    out(
        'INSERT INTO transaction (transaction_id, uid, deleted, type, category_id, account_id, '
        'transaction_time, timezone_utc_offset, amount, related_id, related_account_id, '
        'related_account_amount, hide_amount, comment, geo_longitude, geo_latitude, created_ip, '
        'scheduled_created, created_unix_time, updated_unix_time, deleted_unix_time) VALUES '
        '(%d, %d, false, %d, %d, %d, %d, %d, %d, %d, %d, %d, false, %s, 0, 0, %s, false, '
        '%d, %d, 0);\n' % (
            txn['transaction_id'], UID, txn['type'], txn['category_id'], txn['account_id'],
            txn['transaction_time'], TZ_OFFSET_MIN, txn['amount'], txn['related_id'],
            txn['related_account_id'], txn['related_account_amount'],
            sql_str(txn['comment']), sql_str('127.0.0.1'), _now, _now,
        )
    )

out('\n-- Plantillas programadas (%d)\n' % len(templates))
for tpl in templates:
    out(
        'INSERT INTO transaction_template (template_id, uid, deleted, template_type, name, type, '
        'category_id, account_id, scheduled_frequency_type, scheduled_frequency, '
        'scheduled_start_time, scheduled_end_time, scheduled_at, scheduled_timezone_utc_offset, '
        'tag_ids, amount, related_account_id, related_account_amount, hide_amount, comment, '
        'display_order, hidden, created_unix_time, updated_unix_time, deleted_unix_time) VALUES '
        '(%d, %d, false, 2, %s, %d, %d, %d, 2, %s, %d, %d, %d, %d, %s, %d, 0, 0, false, %s, '
        '%d, false, %d, %d, 0);\n' % (
            tpl['template_id'], UID, sql_str(tpl['name']), tpl['type'], tpl['category_id'],
            tpl['account_id'], sql_str(tpl['scheduled_frequency']), tpl['scheduled_start_time'],
            tpl['scheduled_end_time'], tpl['scheduled_at'], TZ_OFFSET_MIN, sql_str(''),
            tpl['amount'], sql_str(tpl['comment']), tpl['display_order'], _now, _now,
        )
    )

out('\n')
for account in accounts:
    out('-- %-22s saldo final %15d\n' % (account['name'], balances[account['account_id']]))
