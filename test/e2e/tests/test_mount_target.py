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

"""Integration tests for the EFS MountTarget API.
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
from .test_file_system import simple_file_system

RESOURCE_PLURAL = "mounttargets"

CREATE_WAIT_AFTER_SECONDS = 100
UPDATE_WAIT_AFTER_SECONDS = 15
DELETE_WAIT_AFTER_SECONDS = 30

@pytest.fixture(scope="module")
def simple_mount_target(efs_client, simple_file_system):
    (_, cr, file_system_id) = simple_file_system

    resource_name = random_suffix_name("efs-mt", 24)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["MOUNT_TARGET_NAME"] = resource_name
    replacements["FILE_SYSTEM_ID"] = file_system_id

    # Load efs CR
    resource_data = load_efs_resource(
        "mount_target",
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
    assert 'mountTargetID' in cr['status']
    mount_target_id = cr['status']['mountTargetID']

    yield (ref, cr, mount_target_id)

    _, deleted = k8s.delete_custom_resource(
        ref,
        period_length=DELETE_WAIT_AFTER_SECONDS,
    )
    assert deleted

    time.sleep(DELETE_WAIT_AFTER_SECONDS)

    validator = EFSValidator(efs_client)
    assert not validator.mount_target_exists(mount_target_id)

@service_marker
@pytest.mark.canary
class TestMountTarget:
    def test_create_delete(self, efs_client, simple_mount_target):
        (ref, _, mount_target_id) = simple_mount_target
        assert mount_target_id is not None

        validator = EFSValidator(efs_client)
        assert validator.mount_target_exists(mount_target_id)
        
        cr = k8s.get_resource(ref)
        assert 'spec' in cr
        assert 'ipAddress' in cr['spec']
