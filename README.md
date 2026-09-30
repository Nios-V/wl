# wl

Una bitácora de trabajo para la terminal.

Llega la daily, te preguntan qué hiciste ayer y la mente queda en blanco. `wl` existe para eso: anotas algo en dos segundos mientras trabajas y al día siguiente lo tienes ahí.

```
$ wl "revisé documentación del proyecto"
Guardada (id 14)

$ wl yesterday
Lunes, 28/09
 09:40  revisé documentación del proyecto
 11:15  reunión con producto para definir reglas de negocio
 16:02  deploy de la API a QA
```

Es un solo binario, sin servidor ni cuentas. Los datos quedan en un SQLite local.

## Instalación

Solo necesitas [Go](https://go.dev/dl/):

```
go install github.com/Nios-V/wl/cmd/wl@latest
```

El binario queda en `~/go/bin` (en Windows, `%USERPROFILE%\go\bin`). Si `wl` no se encuentra, esa carpeta no está en tu `PATH`.

## Uso

| Comando | |
|---|---|
| `wl "texto"` | Anota algo |
| `wl add "texto"` | Lo mismo, explícito |
| `wl today` | Lo de hoy |
| `wl yesterday` | Lo del día hábil anterior |
| `wl undo` | Borra la última entrada (pregunta antes) |

### Tags

Cualquier `#palabra` dentro del texto se vuelve un tag. Se guarda en minúscula y se muestra aparte:

```
wl "se comienza desarrollo de proyecto en #Go"
```

### Un par de cosas que conviene saber

**Usa comillas.** Tanto en PowerShell como en bash, un `#` al inicio de una palabra empieza un comentario, así que el shell se come el tag antes de que llegue a `wl`.

**"Ayer" es el día hábil anterior.** El lunes, `wl yesterday` muestra el viernes, y lo mismo pasa el sábado y el domingo. Por ahora no considera feriados.

**Si el texto empieza con un subcomando**, `wl` lo toma como comando. Para anotar literalmente "hoy toca revisar logs", usa `wl add "hoy toca revisar logs"`.

## Dónde quedan los datos

En `~/.wl/wl.db`. Es un SQLite normal, así que puedes abrirlo con cualquier cliente si quieres hacer consultas propias. Las horas se guardan en UTC y se muestran en tu hora local.

Cada máquina tiene su propia base.

## Desarrollo

```
cmd/wl/          punto de entrada
internal/cli/    comandos (cobra)
internal/entry/  modelo, tags y fechas
internal/store/  SQLite y migraciones
internal/render/ formato de salida
```

## Lo que viene

- Mostrar los commits del día junto a las entradas, leídos directo de los repos locales
- `wl week`, `wl tag <nombre>`, `wl search <texto>`
