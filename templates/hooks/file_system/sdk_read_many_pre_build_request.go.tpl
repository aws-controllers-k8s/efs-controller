	if r.ko.Spec.Backup != nil {
		msg := fmt.Errorf("field Backup is a NO-OP field. Set BackupPolicy.Status instead")
		return nil, ackerr.NewTerminalError(msg)
	}
	if r.ko.Status.FileSystemID == nil {
		if r, err = rm.withFileSystemIDFromToken(ctx, r); err != nil {
			return nil, err
		}
	}