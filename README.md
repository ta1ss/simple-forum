# Forum
>tl;dr --> [link to website](https://176.112.158.14:8443/)


As a final project of ==Golang branch==, the purpose of this task was to build a forum using **Go** as a main language. 

### Implemented features:
+ Users can create posts, comments and like/dislike other peoples content
+ Ability to log-in/register with a local, ~~Google~~(not working without domain) or GitHub account
+ Adding images to your post
+ Basic security with encrypting passwords, UUID log-in sessions, HTTPS protocol and rate limiting.
+ Basic moderation system with two advanced account types and a system for flagging and reviewing content
+ Notification system for the users
+ Ability to edit posts and comments after submitting


### Account types:
> **Guest** - Non-logged user that can just browse the site.

> **User** - Logged-in user that can read, write and edit their content.

>**Moderator** - User with special privileges and access that can moderate and edit other users content.

>**Admin** - As above but can also manage the account access rights, for example give moderator rights.

--- 
### How To Run
The project is hosted on [this IP](https://176.112.158.14:8443/).

NB! It's hosted through a self-generated SSL Certificate for the learning purpose so the connection will be prompted as not secure. 

## Implementation
- Backend: `Golang`
- Frontend: `HTML` & `CSS`
- Database: `Sqlite3`

## Authors
*Viktor Veertee* & *Enri Suimets*

All rights reserved.
