# marketplace-auth-service

- [Project bootstrap](#project-bootstrap)

- [Commit style](#commit-style)

- [Other scripts](#other-automation-scripts)

### Project bootstrap

For download dependencies use command:

```sh
go get -u -v -f all
```

or

```sh
go mod download
```

### Commit style

As we stated in a badge above, our project is commitizen-friendly. All of our commits follow the [`commitizen format`](https://gist.github.com/stephenparish/9941e89d80e2bc58a153#format-of-the-commit-message):

```
<type>(<optional_scope>): <subject>
<BLANK LINE>
<optional_body>
```

Example:

```
fix(transfers): fixed rerouting bug for cash by code transfer operation

Please see files cash.ts and cash.service.ts and take a look at new private methods. Make sure you understand new rerouting algorithm.
```

```
feat: implemented new authorization logic for trusted users
```

#### Other scripts
##### For generate RSA private and public key you can use the command

```sh
openssl genrsa -out your_name_key.pem 2048
```

for more info please check Makefile

