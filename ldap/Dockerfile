FROM osixia/openldap

LABEL maintainer="administrador@umag.cl"

ENV LDAP_ORGANISATION="Universidad" \
LDAP_DOMAIN="umag.cl"

COPY ./schema/sci.schema /container/service/slapd/assets/config/bootstrap/schema/sci.schema
COPY ./ldifs/usuarios.ldif /container/service/slapd/assets/config/bootstrap/ldif/usuarios.ldif
COPY ./ldifs/grupos.ldif /container/service/slapd/assets/config/bootstrap/ldif/grupos.ldif