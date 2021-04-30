package notify

type SubscriptionAlterationResult struct {
	log ChangeLog
}

func (c SubscriptionAlterationResult) ChangeLog() ChangeLog {
	return c.log
}

func (c SubscriptionAlterationResult) OK() bool {
	if len(c.ChangeLog().NotFoundUsers()) > 0 || len(c.ChangeLog().NotFoundClients()) > 0 {
		return false
	}

	for _, log := range c.log.Topics() {
		if len(log.FailedClients()) > 0 || len(log.FailedUsers()) > 0 {
			return false
		}
	}

	return true
}