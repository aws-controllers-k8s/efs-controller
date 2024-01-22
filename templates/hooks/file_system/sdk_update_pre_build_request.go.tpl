
	if delta.DifferentAt("Spec.Tags") {
		err := syncTags(
			ctx, rm.sdkapi, rm.metrics, 
			string(*desired.ko.Status.ACKResourceMetadata.ARN), 
			desired.ko.Spec.Tags, latest.ko.Spec.Tags,
		)
		if err != nil {
			return nil, err
		}
	}
	if delta.DifferentAt("Spec.Policy") {
		err := rm.syncPolicy(ctx, desired)
		if err != nil {
			return nil, err
		}
	}
	if delta.DifferentAt("Spec.BackupPolicy") {
		err := rm.syncBackupPolicy(ctx, desired)
		if err != nil {
			return nil, err
		}
	}
	if delta.DifferentAt("Spec.LifecyclePolicies") {
		err := rm.syncLifecyclePolicies(ctx, desired)
		if err != nil {
			return nil, err
		}
	}
	if delta.DifferentAt("Spec.FileSystemProtection") {
		err := rm.syncFilesystemProtection(ctx, desired)
		if err != nil {
			return nil, err
		}
	}
	// To trigger to normal update we need to make sure that at least
	// one of the following fields are different.
    if !delta.DifferentAt("Spec.ProvisionedThroughputInMiBps") && !delta.DifferentAt("Spec.ThroughputMode") {
        return desired, nil
    }