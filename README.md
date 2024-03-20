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

### Creating Tasks

Tasks can be added to existing user stories within an active Sprint by members of the development team or SCRUM master on the sprint backlog page. The form is as follows:

- **Title**: Short, descriptive title. (Required)
- **Description**: Detailed description of what the task involves. (Required)
- **Time Estimate**: Time needed to complete the task in hours. Use of decimals for partial hours (e.g., 1.5 for one and a half hours). (Required)
- **User**: Optionally, specify a team member who is suggested to take the task. Note that tasks are not assigned until the team member accepts the task. (Optional)
- **Story ID**: Automatically filled based on the project. (Hidden)
- **Sprint ID**: Automatically filled based on the project. (Hidden)
- **Project ID**: Automatically filled based on the project. (Hidden)

After filling out the necessary information, submit the form by clicking on the create button to add the task to the current user story within the specified Sprint.

### Creating Sprints

Sprints can be created in the current project by SCRUM master on the product backlog page. The form is as follows:

- **Start Date**: Date when the Sprint begins. Must be a future or current date. (Required)
- **End Date**: Date when the Sprint ends. Must be in the future and after the start date. (Required)
- **Velocity**: The expected amount of work (in story points) the team believes can be completed during the Sprint. (Required)
- **Project ID**: Automatically filled based on the project. (Hidden)

After filling out the necessary information, submit the form by clicking on the create button to add the sprint to the current project.

### Adding user stories to sprint

Stories can be added to a sprint by SCRUM master through the product backlog page.

Adding a new user story to the currently active sprint involves the following steps:

1. **Navigation to backlog**: All available stories not yet assigned to a sprint are listed in the backlog tab, located on the left side of the page.
2. **Story selection**: Stories can be selected through checkboxes. Note that only the stories with a defined time complexity can be selected.
3. **Add to sprint**: After selection, click the "Add to Sprint" button to assign the chosen stories to the currently active sprint. This action moves the stories from the backlog to the sprint view.
4. **Removing from sprint**: Stories that are not yet realized (completed) can be removed from the sprint by clicking the 'X' button within the sprint view, which will return them to the backlog.

### Product backlog

The Product Backlog page is accessible to all project members. It is structured into two main sections:

- **Backlog Tab**
  - Located on the left side of the page, this tab contains all unrealized user stories that have not been assigned to a sprint yet.
- **Sprint Views**
  - Positioned on the right, this section features all sprints along with their associated user stories.
  - For each sprint its title (sequence number) and date range are presented.
  - A badge that indicates the current status (Done, Active, or Upcoming) is displayed adjacent to the sprint title.
  - For the currently active sprint, a "Sprint backlog" button is available, providing a direct link to the sprint backlog page.

#### User stories

Within both the backlog tab and the sprint view, user stories are listed in a table format. The basic information presented is as follows:

- **Title**: Short, descriptive title.
- **Description**: Detailed story explanation.
- **Priority**: Importance of the story (one of the following options: Must have, Should have, Could have, Won't have this time).
- **Assignment Status**: Indicates whether a team member has been assigned to the story.
- **Realization Status**: Denotes whether the story has been realized.

To the right of each story is a "Details" button which provides access to additional information.
