	if existing := fileSystemAlreadyCreated(desired, err); existing != nil {
		rm.metrics.RecordAPICall("CREATE", "CreateFileSystem", err)
		rlog.Info("file system already created for this resource, using it", "id", *existing.ko.Status.FileSystemID)
		return existing, nil
	}