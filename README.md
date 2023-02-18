# Forum

This project consists of creating a simple forum using mainly GO & HTML:
Features:
+ Users can create posts, comments and like/dislike other peoples content
+ Ability to log-in/register with a local, Google or GitHub account
+ Adding images to your post
+ Basic security with encrypting passwords, UUID log-in sessions, HTTPS protocol and rate limiting.
+ Basic moderation system with two advanced account types and a system for flagging and reviewing content
+ Notification system for the users
+ Ability to edit posts and comments after submitting


Account types:
+ Guest - Non-logged user that can just browse the site
+ User - Logged in user that can read, write and edit their content
+ Moderator - User with special privileges and access that can moderate and edit other peoples content
+ Admin - As above but can also manage the account access rights, for example give moderator rights.

## How To Run
For this project, you need to have Docker installed. If you don't have it, you can get it [here](https://www.docker.com/).
1. Clone the repo
2. Run the bash script that will take care of building and running the image, as well as stopping and deleting itself after completion. 
 > `bash letsgo.sh`

## Implementation
- Backend: `Golang`
- Frontend: `HTML` & `CSS` & `JS`
- Database: `Sqlite3`
- Container service: `Docker`

## Authors
*Viktor Veertee* & *Enri Suimets*

All rights reserved.
