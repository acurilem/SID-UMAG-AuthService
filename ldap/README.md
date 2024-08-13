# OpenLdap con usuarios y grupos

## Levantar servidor openldap
Para levantar el servidor openldap con usuarios y grupos, se deben ejecutar los siguientes comandos en el directorio del repositorio.
```
docker build -t wshihadeh/openldap:withdata .  
docker compose up -d
```

## Ingresar por PhpldapAdmin
```
Ingresar al sitio http://localhost:8090  
Usuario: cn=admin,dc=umag,dc=cl  
Password: secret  
```

## Usuarios de prueba
```
usuario		password  
jeblanco	secret  
misaaved	secret  
jherreye	secret  
corojas		secret  
```
