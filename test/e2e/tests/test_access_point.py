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

from acktest.resources import random_suffix_name
from acktest import tags
from acktest.k8s import resource as k8s

from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_efs_resource
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.bootstrap_resources import get_bootstrap_resources
from e2e.tests.helper import EFSValidator
from .test_file_system import simple_file_system

RESOURCE_PLURAL = "accesspoints"

CREATE_WAIT_AFTER_SECONDS = 15
UPDATE_WAIT_AFTER_SECONDS = 30
DELETE_WAIT_AFTER_SECONDS = 15

@pytest.fixture(scope="module")
def simple_access_point(efs_client, simple_file_system):
    (_, cr, file_system_id) = simple_file_system

    resource_name = random_suffix_name("efs-ap", 24)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["ACCESS_POINT_NAME"] = resource_name
    replacements["FILE_SYSTEM_ID"] = file_system_id

    resource_data = load_efs_resource(
        "access_point",
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
    assert 'accessPointID' in cr['status']
    access_point_id = cr['status']['accessPointID']

    yield (ref, cr, access_point_id)

    _, deleted = k8s.delete_custom_resource(
        ref,
        period_length=DELETE_WAIT_AFTER_SECONDS,
    )
    assert deleted

    time.sleep(DELETE_WAIT_AFTER_SECONDS)

    validator = EFSValidator(efs_client)
    assert not validator.access_point_exists(access_point_id)

@service_marker
@pytest.mark.canary
class TestAccessPoint:
    def test_create_delete(self, efs_client, simple_access_point):
        (ref, _, access_point_id) = simple_access_point
        assert access_point_id is not None

        validator = EFSValidator(efs_client)
        assert validator.access_point_exists(access_point_id)

        assert k8s.wait_on_condition(ref, "ACK.ResourceSynced", "True", wait_periods=5)

        cr = k8s.get_resource(ref)
        assert cr['status']['lifeCycleState'] == "available"

    def test_update_tags(self, efs_client, simple_access_point):
        (ref, _, access_point_id) = simple_access_point
        assert access_point_id is not None

        validator = EFSValidator(efs_client)
        assert validator.access_point_exists(access_point_id)
        assert k8s.wait_on_condition(ref, "ACK.ResourceSynced", "True", wait_periods=5)

        desired_tags = [{
            "key": "Name",
            "value": "foobar"
        }]
        updates = {
            "spec": {
                "tags": desired_tags,
            },
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(UPDATE_WAIT_AFTER_SECONDS)

        assert k8s.wait_on_condition(ref, "ACK.ResourceSynced", "True", wait_periods=5)

        ap = validator.get_access_point(access_point_id)
        assert ap is not None
        latest_tags = ap["Tags"]

        desired_tags = [{"Key": d["key"], "Value": d["value"]} for d in desired_tags]
        tags.assert_equal_without_ack_tags(
            expected=desired_tags,
            actual=latest_tags
        )