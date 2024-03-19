# tasklify

Live: [tasklify.project-0.dev](https://tasklify.project-0.dev/)

Dev deps:

```sh
sudo corepack enable # This will install and enable yarn

go install github.com/go-task/task/v3/cmd/task@latest
go install github.com/cosmtrek/air@latest
go install github.com/a-h/templ/cmd/templ@latest

docker network create reverse-proxy
```

Run dev on [localhost:8080](localhost:8080):

```sh
task dev
```

Run build:

```sh
task build
```

## Libraries

Session management:

- <https://github.com/gorilla/sessions>
- <https://github.com/alexedwards/scs>
- <https://github.com/gin-contrib/sessions>

Template created from:

- <https://github.com/tomdoesTech/gotth>
- <https://github.com/bnprtr/go-templ-htmx-template>
- <https://github.com/jritsema/go-htmx-tailwind-example>

## Documentation

### Creating User Stories

User stories can be created by product owner and SCRUM master. The form is as follows:

- **Title**: Short, descriptive title. (Required and can't be a duplicate)
- **Description**: Detailed story explanation. (Required)
- **Acceptance Tests**: Click "Add Acceptance Test" to specify criteria for completion. (Optional)
- **Priority**: Select from "Must have," "Should have," "Could have," or "Won't have this time." (Required)
- **Business Value**: Numeric value indicating importance. (Required and can't be negative)
- **Project ID**: Automatically filled based on the project. (Hidden)

After filling out the necessary fields, submit the form to add your story to the project.
