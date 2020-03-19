package vsphere

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccResourceVSphereContentLibrary_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccResourceVSphereContentLibraryPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVSphereContentLibraryConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"vsphere_content_library.library", "id", regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"),
					),
					resource.TestMatchResourceAttr(
						"vsphere_content_library.library", "description", regexp.MustCompile("Library Description"),
					),
					testAccResourceVSphereContentLibraryDescription(regexp.MustCompile("Library Description")),
					testAccResourceVSphereContentLibraryName(regexp.MustCompile("ContentLibrary_test")),
				),
			},
			{
				ResourceName:      "vsphere_content_library.library",
				ImportState:       true,
				ImportStateVerify: true,
				Config:            testAccResourceVSphereContentLibraryConfig(),
			},
		},
	})
}

func testAccResourceVSphereContentLibraryPreCheck(t *testing.T) {
	if os.Getenv("VSPHERE_DATACENTER") == "" {
		t.Skip("set VSPHERE_DATACENTER to run vsphere_content_library acceptance tests")
	}
	if os.Getenv("VSPHERE_DATASTORE") == "" {
		t.Skip("set VSPHERE_DATASTORE to run vsphere_content_library acceptance tests")
	}
}

func testAccResourceVSphereContentLibraryDescription(expected *regexp.Regexp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		library, err := testGetContentLibrary(s, "library")
		if err != nil {
			return err
		}
		if !expected.MatchString(library.Description) {
			return fmt.Errorf("Content Library description does not match. expected: %s, got %s", expected.String(), library.Description)
		}
		return nil
	}
}

func testAccResourceVSphereContentLibraryName(expected *regexp.Regexp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		library, err := testGetContentLibrary(s, "library")
		if err != nil {
			return err
		}
		if !expected.MatchString(library.Name) {
			return fmt.Errorf("Content Library name does not match. expected: %s, got %s", expected.String(), library.Name)
		}
		return nil
	}
}

func testAccResourceVSphereContentLibraryConfig() string {
	return fmt.Sprintf(`
variable "datacenter" {
  type    = "string"
  default = "%s"
}

variable "datastore" {
  type    = "string"
  default = "%s"
}

data "vsphere_datacenter" "dc" {
  name = var.datacenter
}

data "vsphere_datastore" "ds" {
  datacenter_id = data.vsphere_datacenter.dc.id
  name = var.datastore
}

resource "vsphere_content_library" "library" {
  name            = "ContentLibrary_test"
  storage_backing = [ data.vsphere_datastore.ds.id ]
  description     = "Library Description"
}
`,
		os.Getenv("VSPHERE_DATACENTER"),
		os.Getenv("VSPHERE_DATASTORE"),
	)
}
