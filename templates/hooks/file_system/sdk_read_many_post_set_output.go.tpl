	if err := rm.setResourceAdditionalFields(ctx, ko); err != nil {
		return nil, err
	}
	if !filesystemActive(&resource{ko}) {
		return &resource{ko}, requeueWaitState(r)
	}
