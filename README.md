# tasklify

Live: [tasklify.project-0.dev](https://tasklify.project-0.dev/)

Dev deps:

```sh
sudo corepack enable # This will install and enable yarn

go install github.com/go-task/task/v3/cmd/task@latest
go install github.com/cosmtrek/air@latest
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/onsi/ginkgo/v2/ginkgo@latest

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

Common password list from [here](https://github.com/danielmiessler/SecLists/blob/master/Passwords/Common-Credentials/10-million-password-list-top-1000000.txt).

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

### Sprint backlog

The Sprint Backlog page is a central hub for project participants, offering a clear breakdown of sprint-specific tasks and user stories:

- **Sprint Information**

  - At the top, users find a concise summary of the sprint, including its title, duration, and a convenient button to navigate back to the Product Backlog. This section sets the context for the sprint, outlining the timeframe and providing a quick link for easy access to broader project details.

- **Tasks and User Stories Breakdown**
  - The core of the page is dedicated to detailing the user stories within the sprint, along with their associated tasks. Each story is listed with options for deeper exploration or task creation, subject to the user's role and the story's completion status. Tasks are categorized by their current state: Unassigned, Assigned, Active, or Done, making it easy to gauge progress and workload at a glance. If a task has been assigned but not accepted by the user, it is displayed in the "Unassigned" section with a pending status.

### Unassign Task

A user assigned to a task has the option to unassign themselves from it by clicking on the **Unassign button** next to the "Assigned to" table rows in the user story tabels in the Sprint backlog page. After unassigning, a new **Assign button** appears, which allows a new user to assign the task to themselves.
