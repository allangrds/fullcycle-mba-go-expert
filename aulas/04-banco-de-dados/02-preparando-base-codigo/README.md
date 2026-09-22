# 🔧 Passo a passo para subir o ambiente

1. Suba os containers com Docker:

```bash
  make up
```

2. Abra uma nova aba/terminal e acesse o container do MySQL:

```bash
  make mysql
```

3. Conecte-se ao banco de dados MySQL com o usuário root:

```bash
  mysql -u root -p goexpert
```

4. Crie a tabela `products`:

```bash
  CREATE TABLE products (
    id VARCHAR(255),
    name VARCHAR(80),
    price DECIMAL(10,2),
    PRIMARY KEY (id)
  );
```
