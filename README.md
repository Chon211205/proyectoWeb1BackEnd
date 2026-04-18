# Proyecto 1 Web

## Instrucciones para correrlo localmente

### Paso 1: git clone a los siguientes repositorios:

- FrontEnd: https://github.com/Chon211205/proyecto1WebFront.git
- Backend: https://github.com/Chon211205/proyectoWeb1BackEnd.git

### Paso 2: correr FrontEnd.

- Levantar el servidor en la terminal de vscode con el siguiente comando: python -m http.server 5500
- Dirigirse al siguiente puerto donde esta corriendo el FrontEnd: http://localhost:5500

### Paso 3: correr BackEnd.

- Levantar el servidor en otra terminal de vscode con el siguiente comando: go run .
- Dirigirse al frontEnd y esperar a que conecte la API.
- Puerto del Swagger-UI:http://localhost:8080/docs

## Screen Shoot de la app funcionando
### Aplicacion:
<img width="3727" height="1914" alt="image" src="https://github.com/user-attachments/assets/25e6212e-4393-42f0-88b4-d2d84268f220" />

### Swapper UI:
<img width="3732" height="1910" alt="image" src="https://github.com/user-attachments/assets/c401c888-7d2a-4a5f-99b8-808f50c36ed7" />



## Challenges implementados

- [Subjetivo] Calidad visual del cliente — ¿Se ve como una app real o como una tarea?	0 – 30
- [Subjetivo] Calidad del historial de Git — commits descriptivos, progresión lógica, no un solo commit con todo	0 – 20
- [Subjetivo] Organización del código — archivos separados con responsabilidades claras, en ambos repositorios	0 – 20
- pec de OpenAPI/Swagger escrita y precisa (el contrato de la API en YAML o JSON)	20
- Swagger UI corriendo y siendo servido desde el backend (no solo el archivo)	20
- Códigos HTTP correctos en toda la API (201 al crear, 204 al eliminar, 404 si no existe, 400 en input inválido, etc.)	20
- Validación server-side con respuestas de error en JSON descriptivas	20
- Paginación en GET /series con parámetros ?page= y ?limit=	30
- Búsqueda por nombre con ?q=	15
- Ordenamiento con ?sort= y ?order=asc|desc	15
- Exportar la lista de series a CSV — generado manualmente desde JavaScript, sin librerías. El archivo debe descargarse desde el navegador.	20
- Exportar la lista de series a Excel (.xlsx) — generado manualmente desde JavaScript, sin librerías de ningún tipo. El archivo debe ser un .xlsx real que abra correctamente en Excel o LibreOffice. Tip: investiguen el formato SpreadsheetML.	30
- Sistema de rating — tabla propia en la base de datos, endpoints REST propios (POST /series/:id/rating, GET /series/:id/rating, etc.), y visible en el cliente.	30
- Permite subir imágenes. (pongan un tome de como 1 mega a la imagen)	30

## Reflexion

Este proyecto para mi fue interesante debido a que puse mas en practica todo lo aprendido en los laboratorios y ejercicios anteriores en este proyecto como Go, JavaScript, HTML, CSS. Tambien fue divertido como funcionan otras tecnologias relacionadas a los lenguajes anteriores mencionados como por ejemplo el openapi.json y el swagger, es interesante como uno puede comprobar que el backend responde de manera correcta antes de plasmarlo en un frontend. Utilizaria de nuevo estas herramientas ya que agarre mayor practica y confianza.  Por otra parte, al tener repositorios separados me permitio tener mayor orden al momento de trabajar, se me hizo mas comodo al momento de modificar o implementar algo conforme hacia los puntos challenges.
