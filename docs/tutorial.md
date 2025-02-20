# Tutorial: Using GitHub Copilot to Update SDK Examples

This tutorial walks through the process of using GitHub Copilot to update SDK examples based on changelog.json changes.

## Prerequisites

- VS Code with GitHub Copilot extension installed
- Go development environment set up

## Steps

### 1. Set Up Working Set

1. Open VS Code
2. Open the following files in your editor:
   ```
   - golang/sdk-example.go
   - context/changelog.json
   - context/openapi.yaml
   ```

### 2. Prepare Copilot

1. Open the Copilot chat panel:
   - Press `Cmd/Ctrl + Shift + P`
   - Type "Copilot Chat" and select "GitHub Copilot: Open Chat"

2. Set the context by copying and pasting the following prompt:
   ```
   Context: openapi.yaml file represents MongoDB Atlas OpenAPI. 
   sdk-example.go is the file I want to edit. 
   changelog.json contains new changes to the attached OpenAPI.yaml

   SDK uses OpenAPI tag and operationId fields to build SDK methods.

   Action: Act as software agent that for provided SDK example extending it by 
   providing sdk method calls from changelog.json path
   ```

### 3. Generate Code Updates

1. Wait for Copilot to analyze the files and suggest changes
2. Review the suggested changes in the chat panel
3. Accept or modify the suggestions as needed

### 4. Test the Changes

1. Set up environment variables:
   ```bash
   export MONGODB_ATLAS_CLIENT_ID="your_client_id"
   export MONGODB_ATLAS_CLIENT_SECRET="your_client_secret"
   export MONGODB_ATLAS_ORG_ID="your_org_id"
   ```

2. Run the example:
   ```bash
   cd golang
   go run sdk-example.go
   ```

### 5. Automate Updates (Optional)

1. Set up the GitHub Action:
   - Copy the provided workflow file to `.github/workflows/update-sdk-example.yml`
   - Ensure your repository has the necessary permissions set

2. Trigger the automation:
   - Push changes to changelog.json, or
   - Manually trigger the workflow from GitHub Actions tab

## Troubleshooting

### Common Issues

1. **Copilot Not Responding**
   - Try reloading VS Code
   - Check your Copilot authentication

2. **SDK Method Not Found**
   - Check SDK version releases

3. **Authentication Errors**
   - Verify environment variables are set correctly
   - Check service account permissions

## Best Practices

1. Always review Copilot suggestions before accepting
2. Test changes locally before pushing
3. Keep changelog.json up to date
4. Document new SDK methods in code comments
5. Follow the established error handling patterns

## Next Steps

- Review the generated documentation in `docs/user-management.md`
- Set up automated testing for new SDK methods
- Consider contributing improvements back to the SDK
