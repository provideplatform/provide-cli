package common

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/provideservices/provide-go/api/ident"
)

const requireApplicationSelectLabel = "Select an application:"
const requireOrganizationSelectLabel = "Select an organization:"
const requireWorkgroupSelectLabel = "Select a workgroup:"

// RequireApplication is equivalent to a required --application flag
func RequireApplication() error {
	opts := make([]string, 0)
	apps, _ := ident.ListApplications(RequireUserAuthToken(), map[string]interface{}{})
	for _, app := range apps {
		opts = append(opts, *app.Name)
	}

	prompt := promptui.Select{
		Label: requireApplicationSelectLabel,
		Items: opts,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return err
	}

	fmt.Printf("selected application %s at index: %v", *apps[i].Name, i)
	ApplicationID = apps[i].ID.String()
	return nil
}

// RequireWorkgroup is equivalent to a required --workgroup flag
// (yes, this is identical to RequireApplication() with exception to the Printf content...)
func RequireWorkgroup() error {
	opts := make([]string, 0)
	apps, _ := ident.ListApplications(RequireUserAuthToken(), map[string]interface{}{})
	for _, app := range apps {
		opts = append(opts, *app.Name)
	}

	prompt := promptui.Select{
		Label: requireWorkgroupSelectLabel,
		Items: opts,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return err
	}

	fmt.Printf("selected workgroup %s at index: %v", *apps[i].Name, i)
	ApplicationID = apps[i].ID.String()
	return nil
}

// RequireOrganization is equivalent to a required --organization flag
func RequireOrganization() error {
	opts := make([]string, 0)
	orgs, _ := ident.ListOrganizations(RequireUserAuthToken(), map[string]interface{}{})
	for _, org := range orgs {
		opts = append(opts, *org.Name)
	}

	prompt := promptui.Select{
		Label: requireOrganizationSelectLabel,
		Items: opts,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return err
	}

	fmt.Printf("selected organization %s at index: %v", *orgs[i].Name, i)
	OrganizationID = orgs[i].ID.String()
	return nil
}
