	if existing := accessPointAlreadyCreated(desired, err); existing != nil {
		rm.metrics.RecordAPICall("CREATE", "CreateAccessPoint", err)
		rlog.Info("access point already created for this resource, using it", "id", *existing.ko.Status.AccessPointID)
		return existing, nil
	}