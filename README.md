**Description**

This is a simple standup bot for slack.

**Setup**

TODO add app manifest file and show setup instructions

**Test your bot**
In a channel type the below commands - this will create a mock standup with you as the sole participant.
`/standup solicit`
`/standup share`

**Roadmap**

- Provide terraform module for quick setup in fargate
- Update readme with a "how to" to set up slack bot or publish one
- Add native script authorization via roles
- Give go-slackbot a brain via dynamodb or s3 json file
- Catch fatal script errors and prevent exit
- Dockerize
