**Description**

This is a simple standup bot for slack.

**Credit**

This is a rewrite of https://github.com/justmiles/standup-bot to allow for in-app and mutli channel standup meeting configuration.

**Setup**

1. Create standup bot with `app-manifest.yaml` file provided in the root of this directory
1. Get bot and bot user tokens and set environment variables defined in `lib/config/main.go`
1. Set up s3 bucket and provide AWS Creds if you would like standup settings to persists between restarts

**Test your bot**
In a channel type the below commands - this will create a mock standup with you as the sole participant.
`/standup solicit`
`/standup share`

**Users in multiple standups Scenarios**

If User recieves a standup solicitation from multiple channel configurations at the same time, then:
- If User replies as normal, then both standups will receive the same notes.
- If user replies in thread, then only the thread response will be used for the standup notes.

**Roadmap**

- Provide terraform module for quick setup in fargate
- Update readme with a "how to"
