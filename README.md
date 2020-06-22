# auth-poisson

Servizio di autenticazione attraverso LDAP

## Interfaccia

- `POST /login`

    ```json
    { 
        "username": "<username>", 
        "password": "<password>" 
    }
    ```

    Controlla le cretenziali dell'utente `<username>` siano `<password>` e ritorna un token identificativo per richieste successive (se un token è già associato all'utente ne crea uno nuovo).

- `POST /logout`

    ```json
    {
        "username": "<username>"
    }
    ```

    or

    ```json
    {
        "token": "<token>"
    }
    ```

- `POST /logout`

    ```json
    { "username": "<username>" } | { "token": "<token>" }
    ```

    Distrugge il token dell'utente `<username>` o eventualmente dell'utente associato a `<token>`.

- `GET /token?username=<username>`

    Ritorna il token associato all'utente

- `GET /q?username=<username>`

    Ritorna tutte le informazioni pubbliche dell'utente `<username>` fornite da LDAP.

- `POST /update`

    ```json
    {
        // A valid token for the given user  
        "token": "<token>",
        // Name of the property to change  
        "property": "<email | ...>",
        // New value to set
        "value": "<value>"
    }
    ```

    Cambia la proprietà `<property>` dell'utente associato a `<token>` con il nuovo valore fornito.

- `GET /debug`

    Logga informazioni di debug sulle sessioni correnti.