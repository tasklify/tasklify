# tasklify

Dev deps:

```sh
sudo corepack enable # This will install and enable yarn

go install github.com/go-task/task/v3/cmd/task@latest
go install github.com/cosmtrek/air@latest
go install github.com/a-h/templ/cmd/templ@latest
```

Run dev on [localhost:8080](localhost:8080):

```sh
task dev
```

Run build:

```sh
task build
```

Template created from:

- https://github.com/tomdoesTech/gotth
- https://github.com/bnprtr/go-templ-htmx-template
- https://github.com/jritsema/go-htmx-tailwind-example
