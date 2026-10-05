# wl

Una bitácora de trabajo para la terminal.

Llega la daily, te preguntan qué hiciste ayer y la mente queda en blanco. `wl` existe para eso: anotas algo en dos segundos mientras trabajas y al día siguiente lo tienes ahí.

```
$ wl "revisé documentación del proyecto"
Guardada (id 14)

$ wl yesterday
Lunes, 28/09
 09:40 revisé documentación del proyecto
 11:15 reunión con producto para definir reglas de negocio
 16:02 deploy de la API a QA

 Resumen: quedó lista la API en QA, falta validar reglas con producto

Pendientes
 [3] 25/09  pedir accesos a QA
 [5] 28/09  revisar PR de autenticación
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
| `wl todo "texto"` | Agrega un pendiente |
| `wl todo` | Lista los pendientes abiertos |
| `wl todo done <id>` | Marca un pendiente como hecho |
| `wl todo drop <id>` | Elimina un pendiente |
| `wl summary "texto"` | Resumen del día |
| `wl sprint` | Todo lo anotado en el sprint actual |
| `wl sprint summary "texto"` | Resumen del sprint |

### Tags

Cualquier `#palabra` dentro del texto se vuelve un tag. Se guarda en minúscula y se muestra aparte:

```
wl "se comienza desarrollo de proyecto en #Go"
```

### Pendientes

Lo que tienes que hacer pero todavía no haces. No pertenecen a un día: aparecen en amarillo al final de `today` y `yesterday`, con la fecha en que los anotaste, hasta que los cierres.

```
wl todo "pedir accesos a QA"
wl todo                  # [3] 25/09  pedir accesos a QA
wl todo done 3           # hecho, deja de aparecer
wl todo drop 3           # ya no aplica, se borra
```

### Resumen del día

Una frase para cerrar el día. Sale en cian después de los logs, así que al día siguiente `wl yesterday` te deja la daily casi armada.

```
wl summary "quedó lista la API en QA"
wl summary -y "..."      # para el día hábil anterior, si te acuerdas a la mañana siguiente
```

Si lo escribes de nuevo, reemplaza al anterior.

### Sprint

`wl sprint` muestra todas las entradas del sprint agrupadas por día, sin pendientes ni resúmenes diarios, y al final el resumen del sprint si existe.

```
wl sprint                         # sprint actual
wl sprint -p                      # el anterior, para la retro
wl sprint summary "..."           # resumen del sprint actual
wl sprint summary -p "..."        # resumen del anterior
```

## Configurar el sprint

`wl` no sabe cuándo empiezan los sprints de tu equipo, así que hay que decírselo **una sola vez**. Le das la fecha de inicio de cualquier sprint y desde ahí cuenta bloques de 14 días hacia adelante y hacia atrás. No hay que abrir ni cerrar sprints.

1. Averigua el primer día de tu sprint actual (o de cualquier otro, da lo mismo).

2. Crea el archivo `~/.wl/config.yaml` (en Windows, `C:\Users\<tu-usuario>\.wl\config.yaml`):

   ```
   code ~/.wl/config.yaml
   ```

3. Pega esto con tu fecha, en formato `AAAA-MM-DD`:

   ```yaml
   sprint:
     start: 2026-09-28
     days: 14
   ```

   `days` es lo que dura cada sprint. Si lo omites, son 14.

4. Comprueba que las fechas calcen:

   ```
   $ wl sprint
   Sprint 28/09 → 11/10
   ```

   Si salen corridas, la fecha de `start` no es un inicio de sprint.

Solo hay que volver a tocarlo si el equipo cambia el calendario (un sprint que se alarga, por ejemplo). En ese caso pones en `start` el inicio del sprint nuevo.

Sin este archivo todo lo demás funciona igual. Solo `wl sprint` te avisa que falta.

## Un par de cosas que conviene saber

**Usa comillas.** Tanto en PowerShell como en bash, un `#` al inicio de una palabra empieza un comentario, así que el shell se come el tag antes de que llegue a `wl`.

**"Ayer" es el día hábil anterior.** El lunes, `wl yesterday` muestra el viernes, y lo mismo pasa el sábado y el domingo. Por ahora no considera feriados.

**Si el texto empieza con un subcomando**, `wl` lo toma como comando. Para anotar literalmente "today toca revisar logs", usa `wl add "today toca revisar logs"`. Con los pendientes pasa algo parecido: `wl todo done con la API` falla porque `done` espera un id, pero `wl todo "done con la API"` lo guarda como texto.

**Sin colores:** si tu terminal no los muestra, prueba con Windows Terminal o la de VS Code. También puedes apagarlos a propósito con la variable `NO_COLOR=1`.

## Dónde quedan los datos

En `~/.wl/`:

- `wl.db`: un SQLite normal, así que puedes abrirlo con cualquier cliente si quieres hacer consultas propias. Las horas se guardan en UTC y se muestran en tu hora local.
- `config.yaml`: opcional, solo para el sprint.

Cada máquina tiene su propia base. Al actualizar `wl`, la base se migra sola la primera vez que lo usas.

## Desarrollo

```
go run ./cmd/wl today
```

```
cmd/wl/           punto de entrada
internal/cli/     comandos (cobra)
internal/config/  lectura de config.yaml
internal/entry/   modelo, tags y fechas
internal/store/   SQLite y migraciones
internal/render/  formato y colores de salida
```

## Lo que viene

- Mostrar los commits del día junto a las entradas, leídos directo de los repos locales
- `wl week`, `wl tag <nombre>`, `wl search <texto>`
