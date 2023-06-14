// Build tags are not added as SQL Managed Instance takes 6-8 hours for deployment, Exclude this from the worflow
// Please refer to examples/azure/terraform-azure-sqlmanagedinstance-example/README.md for more details

package test

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureSQLManagedInstanceExample(t *testing.T) {
	t.Parallel()

	uniquePostfix := strings.ToLower(random.UniqueId())
	expectedLocation := "West US2"
	expectedAdminLogin := "sqlmiadmin"
	expectedSQLMIState := "Ready"
	expectedSKUName := "GP_Gen5"
	expectedDatabaseName := "testdb"

	// Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-sqlmanagedinstance-example",
		Vars: map[string]interface{}{
			"postfix":       uniquePostfix,
			"location":      expectedLocation,
			"admin_login":   expectedAdminLogin,
			"sku_name":      expectedSKUName,
			"sqlmi_db_name": expectedDatabaseName,
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables
	expectedResourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	expectedManagedInstanceName := terraform.Output(t, terraformOptions, "managed_instance_name")

	// Get the SQL Managed Instance details and assert them against the terraform output
	actualSQLManagedInstance := azure.GetManagedInstance(t, expectedResourceGroupName, expectedManagedInstanceName, "")
	actualSQLManagedInstanceDatabase := azure.GetManagedInstanceDatabase(t, expectedResourceGroupName, expectedManagedInstanceName, expectedDatabaseName, "")
	expectedDatabaseStatus := "Online"

	assert.Equal(t, expectedManagedInstanceName, *actualSQLManagedInstance.Name)
	assert.Equal(t, expectedLocation, *actualSQLManagedInstance.Location)
	assert.Equal(t, expectedSKUName, *actualSQLManagedInstance.Sku.Name)
	assert.Equal(t, expectedSQLMIState, *actualSQLManagedInstance.ManagedInstanceProperties.State)

	assert.Equal(t, expectedDatabaseName, *actualSQLManagedInstanceDatabase.Name)
	assert.Equal(t, expectedDatabaseStatus, actualSQLManagedInstanceDatabase.ManagedDatabaseProperties.Status)

}
