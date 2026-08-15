# NAS Deploy Instructions

Despliegue de `ezbookkeeping` en un NAS (UGOS / UGREEN, arquitectura **aarch64/arm64**) que no soporta `docker compose build` porque le falta el plugin `buildx`. La imagen se compila en la máquina de desarrollo (x86_64) y se transfiere ya construida.

## 1. Compilar la imagen para la arquitectura del NAS

El NAS es ARM64, distinto al x86_64 de esta máquina, así que hay que cross-compilar con `buildx` + QEMU y pasar `BUILD_PIPELINE=1` para que los tests que llaman APIs externas de bancos no rompan el build:

```bash
docker buildx build --platform linux/arm64 --build-arg BUILD_PIPELINE=1 -t ezbookkeeping:local --load .
```

> Si alguna vez el NAS cambia de arquitectura, verifica con `ssh <usuario>@<nas> "uname -m"` (`x86_64` → `--platform linux/amd64`, `aarch64` → `--platform linux/arm64`).

## 2. Exportar la imagen a un .tar

```bash
docker save ezbookkeeping:local -o ezbookkeeping-local.tar
```

## 3. Transferir el .tar al NAS

El cliente SSH moderno (OpenSSH ≥ 9.0) usa SFTP por defecto en `scp`, y el `sshd` del NAS no siempre lo soporta bien — usar el flag `-O` para forzar el protocolo SCP clásico:

```bash
scp -O ezbookkeeping-local.tar Maldonado@maldocloud.local:/volume1/docker/app-financiera
```

Si `scp -O` sigue fallando, alternativa con `rsync`:

```bash
rsync -avP -e ssh ezbookkeeping-local.tar Maldonado@maldocloud.local:/volume1/docker/app-financiera
```

## 4. Cargar la imagen en el Docker del NAS

```bash
ssh Maldonado@maldocloud.local
docker load -i /volume1/docker/app-financiera/ezbookkeeping-local.tar
docker images ezbookkeeping   # confirma que aparece ezbookkeeping:local
```

## 5. Copiar el compose sin sección `build:`

El repo trae `docker-compose.nas.yml`, idéntico a `docker-compose.yml` pero sin `build:` (para que el NAS nunca intente compilar ni hacer un pull ciego, solo use la imagen ya cargada). Cópialo junto con un `.env` que defina `POSTGRES_PASSWORD` a `/volume1/docker/app-financiera` en el NAS.

## 6. Levantar el stack en el NAS

```bash
docker compose -f docker-compose.nas.yml up -d
```

`postgres:17-alpine` se descarga solo (es imagen pública); solo `ezbookkeeping:local` necesita el paso manual de save/load.

## Actualizar tras cambios de código

Repetir el ciclo completo cada vez que cambie el código:

```bash
docker buildx build --platform linux/arm64 --build-arg BUILD_PIPELINE=1 -t ezbookkeeping:local --load .
docker save ezbookkeeping:local -o ezbookkeeping-local.tar
scp -O ezbookkeeping-local.tar Maldonado@maldocloud.local:/volume1/docker/app-financiera
ssh Maldonado@maldocloud.local "docker load -i /volume1/docker/app-financiera/ezbookkeeping-local.tar && docker compose -f /volume1/docker/app-financiera/docker-compose.nas.yml up -d"
```
