	if err := rm.setResourceAdditionalFields(ctx, ko); err != nil {
		return &resource{ko}, err
	}
