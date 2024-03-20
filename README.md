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

After filling out the necessary fields, submit the form by clicking on the create button to add your story to the project.

### Creating Tasks

Tasks can be added to existing user stories within an active Sprint by members of the development team or the SCRUM master. The form is as follows:

- **Title**: Short, descriptive title. (Required)
- **Description**: Detailed description of what the task involves. (Required)
- **Time Estimate**: Time needed to complete the task in hours. Use of decimals for partial hours (e.g., 1.5 for one and a half hours). (Required)
- **User**: Optionally, specify a team member who is suggested to take the task. Note that tasks are not assigned until the team member accepts the task. (Optional)
- **Story ID**: Automatically filled based on the project. (Hidden)
- **Sprint ID**: Automatically filled based on the project. (Hidden)
- **Project ID**: Automatically filled based on the project. (Hidden)

After filling out the necessary information, submit the form by clicking on the create button to add the task to the current user story within the specified Sprint.