// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vsphere

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/kube-cloud/terraform-provider-vsphere/vsphere/internal/helper/testhelper"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/kube-cloud/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/kube-cloud/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/kube-cloud/terraform-provider-vsphere/vsphere/internal/helper/virtualmachine"
	"github.com/vmware/govmomi/vim25/types"
)

func TestAccResourceVSphereComputeClusterVMAntiAffinityRule_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			RunSweepers()
			testAccPreCheck(t)
			testAccResourceVSphereComputeClusterVMAntiAffinityRulePreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccResourceVSphereComputeClusterVMAntiAffinityRuleExists(false),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVSphereComputeClusterVMAntiAffinityRuleConfig(2, true),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceVSphereComputeClusterVMAntiAffinityRuleExists(true),
					testAccResourceVSphereComputeClusterVMAntiAffinityRuleMatchBase(
						true,
						false,
						"terraform-test-cluster-affinity-rule",
					),
					testAccResourceVSphereComputeClusterVMAntiAffinityRuleMatchMembership(),
				),
			},
			{
				ResourceName:      "vsphere_compute_cluster_vm_anti_affinity_rule.cluster_vm_anti_affinity_rule",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					cluster, err := testGetComputeClusterFromDataSource(s, "rootcompute_cluster1")
					if err != nil {
						return "", err
					}

					rs, ok := s.RootModule().Resources["vsphere_compute_cluster_vm_anti_affinity_rule.cluster_vm_anti_affinity_rule"]
					if !ok {
						return "", errors.New("no resource at address vsphere_compute_cluster_vm_anti_affinity_rule.cluster_vm_anti_affinity_rule")
					}
					name, ok := rs.Primary.Attributes["name"]
					if !ok {
						return "", errors.New("vsphere_compute_cluster_vm_anti_affinity_rule.cluster_vm_anti_affinity_rule has no name attribute")
					}

					m := make(map[string]string)
					m["compute_cluster_path"] = cluster.InventoryPath
					m["name"] = name
					b, err := json.Marshal(m)
					if err != nil {
						return "", err
					}

					return string(b), nil
				},
				Config: testAccResourceVSphereComputeClusterVMAntiAffinityRuleConfig(1, true),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceVSphereComputeClusterVMAntiAffinityRuleExists(true),
					testAccResourceVSphereComputeClusterVMAntiAffinityRuleMatchMembership(),
				),
			},
		},
	})
}

func TestAccResourceVSphereComputeClusterVMAntiAffinityRule_updateEnabled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			RunSweepers()
			testAccPreCheck(t)
			testAccResourceVSphereComputeClusterVMAntiAffinityRulePreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccResourceVSphereComputeClusterVMAntiAffinityRuleExists(false),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVSphereComputeClusterVMAntiAffinityRuleConfig(2, true),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceVSphereComputeClusterVMAntiAffinityRuleExists(true),
					testAccResourceVSphereComputeClusterVMAntiAffinityRuleMatchBase(
						true,
						false,
						"terraform-test-cluster-affinity-rule",
					),
					testAccResourceVSphereComputeClusterVMAntiAffinityRuleMatchMembership(),
				),
			},
			{
				Config: testAccResourceVSphereComputeClusterVMAntiAffinityRuleConfig(2, false),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceVSphereComputeClusterVMAntiAffinityRuleExists(true),
					testAccResourceVSphereComputeClusterVMAntiAffinityRuleMatchBase(
						false,
						false,
						"terraform-test-cluster-affinity-rule",
					),
					testAccResourceVSphereComputeClusterVMAntiAffinityRuleMatchMembership(),
				),
			},
		},
	})
}

func TestAccResourceVSphereComputeClusterVMAntiAffinityRule_updateCount(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			RunSweepers()
			testAccPreCheck(t)
			testAccResourceVSphereComputeClusterVMAntiAffinityRulePreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccResourceVSphereComputeClusterVMAntiAffinityRuleExists(false),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVSphereComputeClusterVMAntiAffinityRuleConfig(2, true),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceVSphereComputeClusterVMAntiAffinityRuleExists(true),
					testAccResourceVSphereComputeClusterVMAntiAffinityRuleMatchBase(
						true,
						false,
						"terraform-test-cluster-affinity-rule",
					),
					testAccResourceVSphereComputeClusterVMAntiAffinityRuleMatchMembership(),
				),
			},
			{
				Config: testAccResourceVSphereComputeClusterVMAntiAffinityRuleConfig(3, true),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceVSphereComputeClusterVMAntiAffinityRuleExists(true),
					testAccResourceVSphereComputeClusterVMAntiAffinityRuleMatchBase(
						true,
						false,
						"terraform-test-cluster-affinity-rule",
					),
					testAccResourceVSphereComputeClusterVMAntiAffinityRuleMatchMembership(),
				),
			},
		},
	})
}

func testAccResourceVSphereComputeClusterVMAntiAffinityRulePreCheck(t *testing.T) {
	if os.Getenv("TF_VAR_VSPHERE_DATACENTER") == "" {
		t.Skip("set TF_VAR_VSPHERE_DATACENTER to run vsphere_compute_cluster_vm_anti_affinity_rule acceptance tests")
	}
	if os.Getenv("TF_VAR_VSPHERE_NFS_DS_NAME") == "" {
		t.Skip("set TF_VAR_VSPHERE_NFS_DS_NAME to run vsphere_compute_cluster_vm_anti_affinity_rule acceptance tests")
	}
	if os.Getenv("TF_VAR_VSPHERE_CLUSTER") == "" {
		t.Skip("set TF_VAR_VSPHERE_CLUSTER to run vsphere_compute_cluster_vm_anti_affinity_rule acceptance tests")
	}
	if os.Getenv("TF_VAR_VSPHERE_PG_NAME") == "" {
		t.Skip("set TF_VAR_VSPHERE_PG_NAME to run vsphere_compute_cluster_vm_anti_affinity_rule acceptance tests")
	}
}

func testAccResourceVSphereComputeClusterVMAntiAffinityRuleExists(expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		info, err := testGetComputeClusterVMAntiAffinityRule(s, "cluster_vm_anti_affinity_rule")
		if err != nil {
			if expected == false {
				if viapi.IsManagedObjectNotFoundError(err) {
					// This is not necessarily a missing rule, but more than likely a
					// missing cluster, which happens during destroy as the dependent
					// resources will be missing as well, so want to treat this as a
					// deleted rule as well.
					return nil
				}
			}
			return err
		}

		switch {
		case info == nil && !expected:
			// Expected missing
			return nil
		case info == nil && expected:
			// Expected to exist
			return errors.New("cluster rule missing when expected to exist")
		case !expected:
			return errors.New("cluster rule still present when expected to be missing")
		}

		return nil
	}
}

func testAccResourceVSphereComputeClusterVMAntiAffinityRuleMatchBase(
	enabled bool,
	mandatory bool,
	name string,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		actual, err := testGetComputeClusterVMAntiAffinityRule(s, "cluster_vm_anti_affinity_rule")
		if err != nil {
			return err
		}

		if actual == nil {
			return errors.New("cluster rule missing")
		}

		expected := &types.ClusterAntiAffinityRuleSpec{
			ClusterRuleInfo: types.ClusterRuleInfo{
				Enabled:      structure.BoolPtr(enabled),
				Mandatory:    structure.BoolPtr(mandatory),
				Name:         name,
				UserCreated:  structure.BoolPtr(true),
				InCompliance: actual.InCompliance,
				Key:          actual.Key,
				RuleUuid:     actual.RuleUuid,
				Status:       actual.Status,
			},
			Vm: actual.Vm,
		}

		if !reflect.DeepEqual(expected, actual) {
			return spew.Errorf("expected %#v got %#v", expected, actual)
		}

		return nil
	}
}

func testAccResourceVSphereComputeClusterVMAntiAffinityRuleMatchMembership() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		actual, err := testGetComputeClusterVMAntiAffinityRule(s, "cluster_vm_anti_affinity_rule")
		if err != nil {
			return err
		}

		if actual == nil {
			return errors.New("cluster rule missing")
		}

		vms, err := testAccResourceVSphereComputeClusterVMAntiAffinityRuleMatchMembershipVMIDs(s)
		if err != nil {
			return err
		}

		expectedSort := structure.MoRefSorter(vms)
		sort.Sort(expectedSort)

		expected := &types.ClusterAntiAffinityRuleSpec{
			ClusterRuleInfo: actual.ClusterRuleInfo,
			Vm:              actual.Vm,
		}

		actualSort := structure.MoRefSorter(actual.Vm)
		sort.Sort(actualSort)
		actual.Vm = []types.ManagedObjectReference(actualSort)

		if !reflect.DeepEqual(expected, actual) {
			return spew.Errorf("expected %#v got %#v", expected, actual)
		}

		return nil
	}
}

func testAccResourceVSphereComputeClusterVMAntiAffinityRuleMatchMembershipVMIDs(s *terraform.State) ([]types.ManagedObjectReference, error) {
	var ids []string
	if rs, ok := s.RootModule().Resources["vsphere_virtual_machine.vm"]; ok {
		ids = []string{rs.Primary.ID}
	} else {
		ids = testAccResourceVSphereComputeClusterVMAntiAffinityRuleGetMultiple(s)
	}

	results, err := virtualmachine.MOIDsForUUIDs(testAccProvider.Meta().(*Client).vimClient, ids)
	if err != nil {
		return nil, err
	}
	return results.ManagedObjectReferences(), nil
}

func testAccResourceVSphereComputeClusterVMAntiAffinityRuleGetMultiple(s *terraform.State) []string {
	var i int
	var ids []string
	for {
		rs, ok := s.RootModule().Resources[fmt.Sprintf("vsphere_virtual_machine.vm.%d", i)]
		if !ok {
			break
		}
		ids = append(ids, rs.Primary.ID)
		i++
	}
	return ids
}

func testAccResourceVSphereComputeClusterVMAntiAffinityRuleConfig(count int, enabled bool) string {
	return fmt.Sprintf(`
%s

variable "vm_count" {
  default = "%d"
}

resource "vsphere_virtual_machine" "vm" {
  count            = "${var.vm_count}"
  name             = "terraform-test-${count.index}"
  resource_pool_id = "${data.vsphere_compute_cluster.rootcompute_cluster1.resource_pool_id}"
  datastore_id     = vsphere_nas_datastore.ds1.id

  num_cpus = 2
  memory   = 2048
  guest_id = "other3xLinuxGuest"

  wait_for_guest_net_timeout = -1

  network_interface {
    network_id = "${data.vsphere_network.network1.id}"
  }

  disk {
    label = "disk0"
    size  = 20
  }
}

resource "vsphere_compute_cluster_vm_anti_affinity_rule" "cluster_vm_anti_affinity_rule" {
  name                = "terraform-test-cluster-affinity-rule"
  compute_cluster_id  = "${data.vsphere_compute_cluster.rootcompute_cluster1.id}"
  virtual_machine_ids = "${vsphere_virtual_machine.vm.*.id}"
	enabled             = %t
}
`,
		testhelper.CombineConfigs(
			testhelper.ConfigDataRootDC1(),
			testhelper.ConfigDataRootHost1(),
			testhelper.ConfigDataRootHost2(),
			testhelper.ConfigResDS1(),
			testhelper.ConfigDataRootComputeCluster1(),
			testhelper.ConfigResResourcePool1(),
			testhelper.ConfigDataRootPortGroup1()),
		count,
		enabled,
	)
}
