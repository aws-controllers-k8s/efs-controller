# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
# 	 http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

"""Integration tests for the EFS FileSystem API.
"""

import pytest
import time
import logging
import json

from acktest.resources import random_suffix_name
from acktest.k8s import resource as k8s

from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_efs_resource
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.bootstrap_resources import get_bootstrap_resources
from e2e.tests.helper import EFSValidator

RESOURCE_PLURAL = "filesystems"

CREATE_WAIT_AFTER_SECONDS = 15
UPDATE_WAIT_AFTER_SECONDS = 30
DELETE_WAIT_AFTER_SECONDS = 15

def efs_policy_from_id(id: str):
    return "{\"Version\":\"2008-10-17\",\"Id\":\"1\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"AWS\":\"arn:aws:iam::637423241197:root\"},\"Action\":\"elasticfilesystem:ClientMount\",\"Resource\":\"arn:aws:elasticfilesystem:us-west-2:637423241197:file-system/FSID\"}]}".replace("FSID", id)

@pytest.fixture(scope="function")
def simple_file_system(efs_client):
    resource_name = random_suffix_name("efs", 24)

    resources = get_bootstrap_resources()
    logging.debug(resources)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["FILE_SYSTEM_NAME"] = resource_name

    # Load efs CR
    resource_data = load_efs_resource(
        "file_system",
        additional_replacements=replacements,
    )
    logging.debug(resource_data)

    # Create k8s resource
    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, RESOURCE_PLURAL,
        resource_name, namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)

    time.sleep(CREATE_WAIT_AFTER_SECONDS)
    cr = k8s.wait_resource_consumed_by_controller(ref)

    assert cr is not None
    assert k8s.get_resource_exists(ref)

    assert cr is not None
    assert 'status' in cr
    assert k8s.wait_on_condition(ref, "ACK.ResourceSynced", "True", wait_periods=15)
    assert 'fileSystemID' in cr['status']
    file_system_id = cr['status']['fileSystemID']

    yield (ref, cr, file_system_id)

    _, deleted = k8s.delete_custom_resource(
        ref,
        period_length=DELETE_WAIT_AFTER_SECONDS,
    )
    assert deleted

    time.sleep(DELETE_WAIT_AFTER_SECONDS)

    validator = EFSValidator(efs_client)
    assert not validator.file_system_exists(file_system_id)


@pytest.fixture(scope="function")
def backup_file_system(efs_client):
    resource_name = random_suffix_name("efs", 24)

    resources = get_bootstrap_resources()
    logging.debug(resources)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["FILE_SYSTEM_NAME"] = resource_name

    # Load efs CR
    resource_data = load_efs_resource(
        "file_system_backup",
        additional_replacements=replacements,
    )
    logging.debug(resource_data)

    # Create k8s resource
    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, RESOURCE_PLURAL,
        resource_name, namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)

    time.sleep(CREATE_WAIT_AFTER_SECONDS)
    cr = k8s.wait_resource_consumed_by_controller(ref)

    assert cr is not None
    assert k8s.get_resource_exists(ref)

    assert cr is not None
    assert 'status' in cr
    assert k8s.wait_on_condition(ref, "ACK.ResourceSynced", "True", wait_periods=15)
    assert 'fileSystemID' in cr['status']
    file_system_id = cr['status']['fileSystemID']

    yield (ref, cr, file_system_id)

    _, deleted = k8s.delete_custom_resource(
        ref,
        period_length=DELETE_WAIT_AFTER_SECONDS,
    )
    assert deleted

    time.sleep(DELETE_WAIT_AFTER_SECONDS)

    validator = EFSValidator(efs_client)
    assert not validator.file_system_exists(file_system_id)

@service_marker
@pytest.mark.canary
class TestFileSystem:

    def test_create_delete(self, efs_client, simple_file_system):
        (_, _, file_system_id) = simple_file_system
        assert file_system_id is not None

        validator = EFSValidator(efs_client)
        assert validator.file_system_exists(file_system_id)
    
    def test_create_update_delete(self, efs_client, simple_file_system):
        (ref, _, file_system_id) = simple_file_system
        assert file_system_id is not None

        validator = EFSValidator(efs_client)
        assert validator.file_system_exists(file_system_id)

        updates = {
            "spec": {
                "throughputMode": "elastic",
                "fileSystemProtection": {
                    "replicationOverwriteProtection": "DISABLED",
                }
            },
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(UPDATE_WAIT_AFTER_SECONDS)

        fs = validator.get_file_system(file_system_id)
        assert fs is not None
        assert fs['FileSystems'][0]['ThroughputMode'] == "elastic"
        assert fs['FileSystems'][0]['FileSystemProtection']['ReplicationOverwriteProtection'] == "DISABLED"


    def test_update_policies(self, efs_client, simple_file_system):
        (ref, _, file_system_id) = simple_file_system
        assert file_system_id is not None

        validator = EFSValidator(efs_client)
        assert validator.file_system_exists(file_system_id)

        policy = efs_policy_from_id(file_system_id)
        updates = {
            "spec": {
                "policy": policy,
            },
        }

        k8s.patch_custom_resource(ref, updates)
        time.sleep(CREATE_WAIT_AFTER_SECONDS)

        observedPolicy = validator.get_file_system_policy(file_system_id)
        assert json.loads(policy) == json.loads(observedPolicy)

    def test_update_backup_policy(self, efs_client, backup_file_system):
        (ref, _, file_system_id) = backup_file_system
        assert file_system_id is not None

        validator = EFSValidator(efs_client)
        assert validator.file_system_exists(file_system_id)
        assert k8s.wait_on_condition(ref, "ACK.ResourceSynced", "True", wait_periods=15)
        tags = validator.list_tags_for_resource(file_system_id)
        logging.error("tags: "+str(tags))

        assert tags is not None
        assert len(tags) > 0

        foundAWSTag = False
        for t in tags:
            if t["Key"].startswith("aws:"):
                foundAWSTag = True
        assert foundAWSTag

        updates = {
            "spec": {
                "backupPolicy": {
                    "status": "DISABLED",
                },
            },
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(UPDATE_WAIT_AFTER_SECONDS)

        bp = validator.get_backup_policy(file_system_id)
        assert bp is not None
        assert bp['Status'] == "DISABLED"

    def test_update_lifecycle_policies(self, efs_client, simple_file_system):
        (ref, _, file_system_id) = simple_file_system
        assert file_system_id is not None

        validator = EFSValidator(efs_client)
        assert validator.file_system_exists(file_system_id)

        updates = {
            "spec": {
                "lifecyclePolicies": [
                    {"transitionToIA": "AFTER_30_DAYS"}
                ],
            },
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(UPDATE_WAIT_AFTER_SECONDS)

        lfps = validator.get_file_system_lifecycle_policy(file_system_id)
        assert lfps is not None
        assert lfps[0]['TransitionToIA'] == "AFTER_30_DAYS"