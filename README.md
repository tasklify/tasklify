# tasklify

Live: [tasklify.project-0.dev](https://tasklify.project-0.dev/)

Dev deps:

```sh
sudo corepack enable # This will install and enable yarn

go install github.com/go-task/task/v3/cmd/task@latest
go install github.com/cosmtrek/air@latest
go install github.com/a-h/templ/cmd/templ@latest

docker network create gateway_prod
```

Run dev on [localhost:8080](localhost:8080):

```sh
task dev
```

Run build:

```sh
task build
```

Run tests (`task dev` has to be running):

```sh
task test
```

## Other

Template created from:

- <https://github.com/tomdoesTech/gotth>
- <https://github.com/bnprtr/go-templ-htmx-template>
- <https://github.com/jritsema/go-htmx-tailwind-example>

Common password list from [here](https://github.com/danielmiessler/SecLists/blob/master/Passwords/Common-Credentials/10-million-password-list-top-1000000.txt).
