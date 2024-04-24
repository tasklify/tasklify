---
Title: Task Logs
---

## Task Logs

When a task is assigned, the assignee can begin logging time by clicking the **Start** button in the sprint backlog. This action triggers an automatic timer. To pause the timer, the user clicks the **Stop** button and can later resume by clicking the **Resume** button. Additionally, the user has the option to manually edit the logged and remaining time for the task by clicking on the duration/remaining time and modifying it. If the timer is not manually stopped within the same day, it will automatically halt at midnight, and the user will be alerted with a red exclamation mark.

### Past Logs

- **Adding Past Logs**

  - A past log can be added by clicking the **Add past log** button in the past logs section. This action opens a form where users can enter the date, duration, and predicted remaining time for the task. Once the necessary fields are completed, the form can be submitted by clicking the **Add log** button. There are several restrictions on the date: it must be a date from the past, not today or in the future; it must fall within the sprint's duration; there must not be an existing log for this date; and it must be a valid date. Additionally, the duration entered must not exceed 24 hours.

- **Editing Past Logs**

  - A past log can be edited if it was made by you by clicking on the duration, entering the new time in the XXhXXm format, and clicking the **Change** button. Similarly, the remaining time can be adjusted by clicking on the indicated time, updating it in the XXhXXm format, and then clicking the **Change** button. The duration entered should not exceed 24 hours or be negative, and the remaining time must not be negative.

- **Deleting Past Logs**

  - A past log can be deleted if it was made by you by clicking the **Delete** button next to the log. This action will remove the log from the past logs section completely and permanently.

### Today's Log

- If a log has not been started today, and the predicted remaining time from the previous day is not 0, the user can start a log for today by clicking the **Start** button in the today's log section. This action will start an automatic timer, which can be paused and resumed as needed. The user can also manually edit the logged and remaining time for today's log by clicking on the duration/remaining time and modifying it. If the timer is not manually stopped within the same day, it will automatically halt at midnight, and the user will be alerted with a red exclamation mark.

<!-- - If the predicted remaining time from the previous log or today's log is 0, the user can end the task by clicking the **End task** button. This action will mark the task as done and the user can no longer log time for the task. -->

### End Task

If the predicted remaining time from the previous log or today's log is 0, the task is automatically marked as **Done**. This indicates that work on the task is finished and the user can no longer log time for it.
However, if the user wishes to resume work on the task at a later time, they can simply update the remaining time to a value greater than zero, and the task's status will revert to **In progress**.
