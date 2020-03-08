# auth-poisson

Servizio di autenticazione attraverso LDAP

## Interfaccia

- `POST /login`

    ```
    { 
        "username": "<username>", 
        "password": "<password>" 
    }
    ```

    Controlla le cretenziali dell'utente `<username>` siano `<password>` e ritorna un token identificativo per richieste successive (se un token è già associato all'utente ne crea uno nuovo).

- `POST /logout`

    ```
    { "username": "<username>" } | { "token": "<token>" }
    ```

    Distrugge il token dell'utente `<username>` o eventualmente dell'utente associato a `<token>`.

- `GET /info?username=<username>`

    Ritorna tutte le informazioni pubbliche dell'utente `<username>` fornite da LDAP.

- `POST /update`

    ```
    {
        "token": "<token>",
        "property": "<email | ...>",
        "value": "<value>"
    }
    ```

    Cambia la proprietà `<property>` dell'utente associato a `<token>` con il nuovo valore fornito.
