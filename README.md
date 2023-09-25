## Dependências

* Docker
* Make

## Rodar 

### Subir projeto 

```bash
make up
```

### Endpoints 

#### Saldo do usuário.

```bash
  curl -v --location --request GET 'http://localhost:3005/users/f2f0e0d1-e37e-45c3-ad06-e6c2a66544fc' 
```

### Fazer transfer

```bash
curl -v --location --request POST 'http://localhost:3005/transfers' \
--header 'Content-Type: application/json' \
--data-raw '{
    "amount" : 10,
    "debtor_id" : "f2f0e0d1-e37e-45c3-ad06-e6c2a66544fc",
    "beneficiary_id" : "089557bc-ddf2-4ec5-8077-d8bf09fe3ddc" 
}'
```