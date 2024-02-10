# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
#	 http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

"""Helper functions for EFS e2e tests
"""

import logging

class EFSValidator:
    def __init__(self, efs_client):
        self.efs_client = efs_client

    def get_file_system(self, filesystem_id: str) -> dict:
        try:
            resp = self.efs_client.describe_file_systems(
                FileSystemId=filesystem_id
            )
            return resp

        except Exception as e:
            logging.debug(e)
            return None

    def file_system_exists(self, filesystem_id) -> bool:
        return self.get_file_system(filesystem_id) is not None
    
    def get_file_system_policy(self, filesystem_id: str) -> dict:
        try:
            resp = self.efs_client.describe_file_system_policy(
                FileSystemId=filesystem_id
            )
            return resp["Policy"]

        except Exception as e:
            logging.debug(e)
            return None
        
    def get_file_system_lifecycle_policy(self, filesystem_id: str) -> dict:
        try:
            resp = self.efs_client.describe_lifecycle_configuration(
                FileSystemId=filesystem_id
            )
            return resp["LifecyclePolicies"]

        except Exception as e:
            logging.debug(e)
            return None
        
    def get_backup_policy(self, filesystem_id: str) -> dict:
        try:
            resp = self.efs_client.describe_backup_policy(
                FileSystemId=filesystem_id
            )
            return resp["BackupPolicy"]

        except Exception as e:
            logging.debug(e)
            return None
        
    def get_mount_target(self, mount_target_id: str) -> list:
        try:
            resp = self.efs_client.describe_mount_targets(
                MountTargetId=mount_target_id
            )
            return resp["MountTargets"][0]

        except Exception as e:
            logging.debug(e)
            return None
        
    def mount_target_exists(self, mount_target_id) -> bool:
        return self.get_mount_target(mount_target_id) is not None