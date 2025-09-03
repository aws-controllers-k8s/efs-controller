	// Check replication status first and requeue if deleting
	if !filesystemActive(r) {
		return nil, requeueWaitState(r)
	}
	
	if !replicationConfigurationActive(r) {
		return nil, requeueWaitReplicationConfiguration
	}

	// Handle replication configuration deletion using sync logic
	// Set desired state to empty for deletion and call sync
	if replicationConfigurationExists(r) {
		// Create a resource with empty replication config to trigger deletion
		deleteResource := &resource{r.ko.DeepCopy()}
		deleteResource.ko.Spec.ReplicationConfiguration = nil

		err = rm.syncReplicationConfiguration(ctx, deleteResource, r)
		if err != nil {
			return nil, err
		}
		// Requeue to wait for deletion to complete
		return nil, requeueWaitReplicationConfiguration
	}