---
Title: Users
---

# Users

## User login

User can login using their username and password.

## Adding users

Only system `Administrator` can create new users. This includes:

- **Username**: Has to be unique.
- **Password**: You have to write the password twice. Has to be at least 12 characters, but shorter than 128. A password also can't be the top 1 million most used passwords.
- **First name**: No limitations.
- **Last name**: No limitations.
- **Email**: Has to be unique.
- **System role**: Can be either `user` or `admin`. Any admin can add users.

## Changing user data

Users can change their own user data, this includes:

- **Username**
- **Password**: You have to write the new password twice.
- **First name**
- **Last name**
- **Email**

Here apply the same constraints for fields as with user creation.

This action has to be authorized with the current password.

## Admin change user data

In addition to user changing their own data, `admins` can change all attributes that the user can for himself and additionally the user system role.
