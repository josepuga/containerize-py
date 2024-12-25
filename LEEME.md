# CONTAINERIZE PYTHON PROJECT

Es una simple pero efectiva utilidad que genera los ficheros de contenedor necesarios para crear una imagen del proyecto Python indicado.

## Características:
- Por defecto crea imágenes desde `python-alpine`, pero puede usarse cualquier `python-slim-bullseye`, `python-slim-bookworm`, etc.
- Selecciona la versión Python que hubiera en el posible Virtual Enviroment o en su defecto, la actual del sistema.

## Instalación:
- En el directorio `bin/` hay ejecutables para Linux y Windows. No requieren librerías adicionales.
- Si quieres clonar el repositorio y compilarlo por ti mismo, aconsejo usar el script `build.sh`.

## Consideraciones:
- No está ampliamente probado, por lo que no se garantiza su correcto funcionamiento para proyectos grandes que requieran muchas dependencias.
- Los ficheros generados son `Containerfile` y `.containerignore`. Ambos son reconocidos por docker (con prioridad respecto a `Dockerfile` y `.dockerignore`) ya que son el nuevo estándar para no hacerlo "marca-dependiente".
- Es necesario que exista un `main.py` y el `requirements.txt`.


## Ejemplo de funcionamiento:
Aunque tiene varias opciones, la forma más sencilla es crear los ficheros en el mismo proyecto y usar la imagen alpine por defecto.
```bash
$ cd /path/to/my_project
$ containerize-py .
```
Esto empezará el proceso que generará los ficheros.
```txt
Checking directories...
Checking for main.py...
Checking for requirements.txt...
Getting Python project version...
Writing Ignore File...
Writing Container File...
Done!.
Check that no files containing sensitive information (passwords, API keys, etc.) are present in your project. If there are any, make sure to include them in .containerignore
```

Se habrán generado 2 ficheros. `Containerfile`:
```Dockerfile
# Generated with containerize-py v0.1.0
FROM python:3.13.0-alpine
WORKDIR /app
COPY . .
RUN pip install --no-cache-dir -r requirements.txt
CMD ["python3", "main.py"]
```

Y otro para excluir ficheros innecesarios. `.containerignore`:
```txt
# Generated with containerize-py v0.1.0
# Virtual enviroments
.venv/
venv/
.env/
env/

# Cache
__pycache__/
*.pyc
*.pyo

# 
Containerfile
.containerignore
```