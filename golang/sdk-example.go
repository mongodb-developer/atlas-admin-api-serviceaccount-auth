package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/atlas-sdk/v20241113001/admin"
)

// Example for Service Account Management API
// Example uses Service Account to create Service Account.
// Please ensure that Service Account has organizational admin permission.

// Required env variables to run example:
// export MONGODB_ATLAS_CLIENT_ID="your_client_id"
// export MONGODB_ATLAS_CLIENT_SECRET="your_client_secret"
// export MONGODB_ATLAS_ORG_ID="your_org_id"
func main() {
	host := os.Getenv("MONGODB_ATLAS_BASE_URL")
	if host == "" {
		host = "https://cloud.mongodb.com"
	}

	// Fetch clientID and clientSecret from environment variables
	clientID := os.Getenv("MONGODB_ATLAS_CLIENT_ID")
	clientSecret := os.Getenv("MONGODB_ATLAS_CLIENT_SECRET")
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")

	if clientID == "" || clientSecret == "" {
		log.Fatal("Missing required environment variables")
	}

	ctx := context.Background()
	sdk, err := admin.NewClient(
		admin.UseBaseURL(host),
		admin.UseOAuthAuth(ctx, clientID, clientSecret),
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	if orgID == "" {
		orgs, _, err := sdk.OrganizationsApi.ListOrganizations(ctx).Execute()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		if orgs.GetTotalCount() == 0 {
			log.Fatalf("no orgs")
		}
		orgID = orgs.GetResults()[0].GetId()
	}

	// 1. Create Service Account
	request := sdk.ServiceAccountsApi.CreateServiceAccount(
		ctx,
		orgID,
		admin.NewOrgServiceAccountRequest(
			"SA created by sdk-example",
			"example",
			[]string{"ORG_OWNER"},
			365*24,
		),
	)
	sa, _, err := request.Execute()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// 2. Rotate secret
	newSecret, _, err := sdk.ServiceAccountsApi.CreateServiceAccountSecret(
		ctx,
		orgID,
		sa.GetClientId(),
		&admin.ServiceAccountSecretRequest{
			SecretExpiresAfterHours: 365 * 24,
		},
	).Execute()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// 3. Delete rotated secret
	_, err = sdk.ServiceAccountsApi.DeleteServiceAccountSecret(
		ctx,
		sa.GetClientId(),
		sa.GetSecrets()[0].GetId(),
		orgID,
	).Execute()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// 4. Create new SDK client using New Service Account
	newSDK, err := admin.NewClient(
		admin.UseBaseURL(host),
		// 4.1 Using ClientId and Secret returned by API
		//  might have up to 2 secrets
		admin.UseOAuthAuth(ctx, sa.GetClientId(), newSecret.GetSecret()),
	)

	// 5. Make request using new Service Account
	projects, _, err := newSDK.ProjectsApi.ListProjectsWithParams(ctx,
		&admin.ListProjectsApiParams{}).Execute()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Projects size: %d\n", projects.GetTotalCount())

	// 6. Remove created Service Account. We would not be able to use it afterward without access to Secret value.
	sdk.ServiceAccountsApi.DeleteServiceAccount(ctx, sa.GetClientId(), orgID)

	// 7. Examples of new User Management operations
	
	// 7.1 List organization users
	orgUsers, _, err := sdk.MongoDBCloudUsersApi.ListOrganizationUsers(ctx, orgID).Execute()
	if err != nil {
		log.Fatalf("Error listing org users: %v", err)
	}
	fmt.Printf("Organization users count: %d\n", orgUsers.GetTotalCount())

	// 7.2 Create a new organization user
	newUserRoles := &admin.OrgUserRolesRequest{
		OrgRoles: []string{"ORG_MEMBER"},
	}
	newUser, _, err := sdk.MongoDBCloudUsersApi.CreateOrganizationUser(ctx, orgID, 
		&admin.CreateOrganizationRequest{
			Username: "new.user@example.com",
			Roles:    newUserRoles,
		}).Execute()
	if err != nil {
		log.Fatalf("Error creating org user: %v", err)
	}

	// 7.3 Add project role to user
	_, _, err = sdk.MongoDBCloudUsersApi.AddProjectRole(ctx, orgID, *newUser.Id, 
		&admin.AddOrRemoveGroupRole{
			Roles: []string{"GROUP_READ_ONLY"},
		}).Execute()
	if err != nil {
		log.Fatalf("Error adding project role: %v", err)
	}

	// 7.4 Get project users
	projectUsers, _, err := sdk.MongoDBCloudUsersApi.ListProjectUsers(ctx, orgID).Execute()
	if err != nil {
		log.Fatalf("Error listing project users: %v", err)
	}
	fmt.Printf("Project users count: %d\n", projectUsers.GetTotalCount())

	// 7.5 Clean up - remove roles and user
	_, _, err = sdk.MongoDBCloudUsersApi.RemoveProjectRole(ctx, orgID, *newUser.Id, 
		&admin.AddOrRemoveGroupRole{
			Roles: []string{"GROUP_READ_ONLY"},
		}).Execute()
	if err != nil {
		log.Fatalf("Error removing project role: %v", err)
	}

	_, err = sdk.MongoDBCloudUsersApi.DeleteOrganizationUser(ctx, orgID, *newUser.Id).Execute()
	if err != nil {
		log.Fatalf("Error deleting org user: %v", err)
	}
}
