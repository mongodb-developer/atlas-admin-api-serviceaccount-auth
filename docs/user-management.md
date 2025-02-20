# MongoDB Atlas User Management Operations

This document describes the User Management operations demonstrated in the SDK example.

## Operations Overview

### 1. List Organization Users
```go
orgUsers, _, err := sdk.MongoDBCloudUsersApi.ListOrganizationUsers(ctx, orgID).Execute()
```
Lists all users in an organization. Returns a paginated result containing user details and their roles.

### 2. Create Organization User
```go
newUserRoles := &admin.OrgUserRolesRequest{
    OrgRoles: []string{"ORG_MEMBER"},
}
newUser, _, err := sdk.MongoDBCloudUsersApi.CreateOrganizationUser(ctx, orgID, 
    &admin.CreateOrganizationRequest{
        Username: "new.user@example.com",
        Roles:    newUserRoles,
    }).Execute()
```
Creates a new user in the organization with specified roles. The user is identified by their email address.

### 3. Add Project Role
```go
_, _, err = sdk.MongoDBCloudUsersApi.AddProjectRole(ctx, orgID, userId, 
    &admin.AddOrRemoveGroupRole{
        Roles: []string{"GROUP_READ_ONLY"},
    }).Execute()
```
Adds project-level roles to an existing user. Available roles include:
- GROUP_READ_ONLY
- GROUP_DATA_ACCESS_ADMIN
- GROUP_OWNER
- And others as defined in the API

### 4. List Project Users
```go
projectUsers, _, err := sdk.MongoDBCloudUsersApi.ListProjectUsers(ctx, orgID).Execute()
```
Returns a list of all users who have access to the specified project, including their roles and permissions.

### 5. Remove Project Role
```go
_, _, err = sdk.MongoDBCloudUsersApi.RemoveProjectRole(ctx, orgID, userId, 
    &admin.AddOrRemoveGroupRole{
        Roles: []string{"GROUP_READ_ONLY"},
    }).Execute()
```
Removes specified project roles from a user. The user retains any other roles not specified in the request.

### 6. Delete Organization User
```go
_, err = sdk.MongoDBCloudUsersApi.DeleteOrganizationUser(ctx, orgID, userId).Execute()
```
Completely removes a user from the organization. This action cannot be undone.

## Error Handling

All operations return error objects that should be checked. Common errors include:
- Invalid role names
- User not found
- Insufficient permissions
- Rate limiting

## Best Practices

1. Always validate role names before making requests
2. Use appropriate error handling for each operation
3. Clean up unused user accounts and roles
4. Follow principle of least privilege when assigning roles
5. Keep track of user IDs for subsequent operations
