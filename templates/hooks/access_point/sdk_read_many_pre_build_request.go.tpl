	if r.ko.Status.AccessPointID == nil {
		if r, err = rm.withAccessPointIDFromToken(ctx, r); err != nil {
			return nil, err
		}
	}