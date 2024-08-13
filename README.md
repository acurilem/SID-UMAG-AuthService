# UMAG-AuthService

Microservicio para loggearse dentro de los servicios de UMAG.

## Docker

Corra el comando:

```bash
docker build -t wshihadeh/openldap:withdata ../ldap

docker compose up -d
```

Al momento de ejecutar el comando se levantar tres componentes:

- LDAP Server: Servidor de LDAP donde se almacenan los usuarios con sus contraseñas.
- LDAP Server Admin: Interfaz web que permite conectar al componente de LDAP Server y servir una web para interactuar con LDAP Server.
- AuthService: Microservicio que se conecta a LDAP Server y base de datos SQL Server de Microsoft, para generar el servicio de autenticación de usuarios, produciendo finalmente un token JWT para autenticarse en el resto de microservicios del proyecto UMAG.

## Swagger

Para consultar la documentación de la API se tiene que acceder a estos puntos:

### Interfaz web:

```
[GET] /swagger/index.html
```

### Archivos JSON

```
[GET] /swagger/doc.json
```
