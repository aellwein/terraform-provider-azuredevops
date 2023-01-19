//go:build (all || permissions || resource_build_Folder_permissions) && (!exclude_permissions || !exclude_resource_build_Folder_permissions)
// +build all permissions resource_build_Folder_permissions
// +build !exclude_permissions !exclude_resource_build_Folder_permissions

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/aellwein/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/aellwein/terraform-provider-azuredevops/azuredevops/internal/utils/datahelper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func hclBuildFolderPermissions(projectName string, path string, permissions map[string]string) string {
	rootPermissions := datahelper.JoinMap(permissions, "=", "\n")
	description := "Integration Test Folder"

	return fmt.Sprintf(`
%s

data "azuredevops_group" "tf-project-readers" {
	project_id = azuredevops_project.project.id
	name       = "Readers"
}

resource "azuredevops_build_folder_permissions" "permissions" {
	project_id  = azuredevops_project.project.id
	principal   = data.azuredevops_group.tf-project-readers.id
	path        = azuredevops_build_folder.test_folder.path

	permissions = {
		%s
	}
}
`,
		testutils.HclBuildFolder(projectName, path, description),
		rootPermissions,
	)
}

func TestAccBuildFolderPermissions_SetPermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config := hclBuildFolderPermissions(projectName, `\test-folder`, map[string]string{
		"ViewBuilds":                 "Allow",
		"EditBuildQuality":           "Allow",
		"RetainIndefinitely":         "Allow",
		"DeleteBuilds":               "Deny",
		"ManageBuildQualities":       "Allow",
		"DestroyBuilds":              "Allow",
		"UpdateBuildInformation":     "Allow",
		"QueueBuilds":                "Allow",
		"ManageBuildQueue":           "Allow",
		"StopBuilds":                 "Allow",
		"ViewBuildDefinition":        "Allow",
		"EditBuildDefinition":        "Allow",
		"DeleteBuildDefinition":      "Deny",
		"AdministerBuildPermissions": "NotSet",
	})
	tfNodeRoot := "azuredevops_build_folder_permissions.permissions"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "project_id"),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "principal"),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "path"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.%", "14"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.ViewBuilds", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.DeleteBuilds", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.DeleteBuildDefinition", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.AdministerBuildPermissions", "notset"),
				),
			},
		},
	})
}

func TestAccBuildFolderPermissions_UpdatePermissions(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config1 := hclBuildFolderPermissions(projectName, `\dir1`, map[string]string{
		"ViewBuilds":         "Deny",
		"EditBuildQuality":   "NotSet",
		"RetainIndefinitely": "Deny",
		"DeleteBuilds":       "Deny",
	})
	config2 := hclBuildFolderPermissions(projectName, `\dir1`, map[string]string{
		"ViewBuilds":         "Deny",
		"EditBuildQuality":   "Allow",
		"RetainIndefinitely": "Deny",
		"DeleteBuilds":       "Deny",
	})
	tfNodeRoot := "azuredevops_build_folder_permissions.permissions"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "project_id"),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "principal"),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "path"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.%", "4"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.ViewBuilds", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.EditBuildQuality", "notset"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.RetainIndefinitely", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.DeleteBuilds", "deny"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "project_id"),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "principal"),
					resource.TestCheckResourceAttrSet(tfNodeRoot, "path"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.%", "4"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.ViewBuilds", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.EditBuildQuality", "allow"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.RetainIndefinitely", "deny"),
					resource.TestCheckResourceAttr(tfNodeRoot, "permissions.DeleteBuilds", "deny"),
				),
			},
		},
	})
}
