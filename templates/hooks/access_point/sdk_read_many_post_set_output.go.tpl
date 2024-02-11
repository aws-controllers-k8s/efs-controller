	if !accessPointActive(&resource{ko}) {
		return &resource{ko}, requeueWaitState(r)
	}
