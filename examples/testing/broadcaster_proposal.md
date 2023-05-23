Broadcaster

### Add user
1. Subscribe the user's channel (and start processing incoming events)
2. Load user's subscriptions and subscribe corresponding topics' channels

### Add a user's subscription
1. Receive an event from the user's channel (user A subscribes a topic B)
2. Subscribe the topic channel
3. (There may be some gap until all broadcaster synchronize their state)

