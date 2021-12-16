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
    { "token": "<token>" }
    ```

    Distrugge il token dell'utente `<username>` o eventualmente dell'utente associato a `<token>`.

- `GET /token?username=<username>`

    Ritorna il token associato all'utente

- `GET /users`

    Ritorna una lista degli utenti con tutte le informazioni pubbliche fornite da LDAP.

- `GET /users?username=<username>`

    Ritorna tutte le informazioni pubbliche dell'utente `<username>` fornite da LDAP.

- **TODO** `POST /update`

    ```js
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

- **TODO** `POST /change-password`

    ```js
    {
        // A valid token for the given user  
        "token": "<token>",
        // New password
        "password": "<new-password>",
    }
    ```

    Cambia la password dell'utente associato a `<token>` con `<new-password>`

- `GET /debug`

    Logga informazioni di debug sulle sessioni correnti.
