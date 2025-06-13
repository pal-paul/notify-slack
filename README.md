# notify-slack

A GitHub Action to send Slack notifications for your GitHub workflow status. This action can be easily integrated into any GitHub workflow to provide Slack notifications for your CI/CD pipeline events.

## Usage

To use this action in your GitHub workflow, add the following step to your `.github/workflows/your-workflow.yml` file:

```yaml
- name: Send Slack Notification
  uses: pal-paul/notify-slack@v1.3.1  # Use the latest version tag
  with:
    # Required inputs
    status: ${{ job.status }}  # success/failure
    notify_when: 'always'      # always/success/failure
    slack_token: ${{ secrets.SLACK_TOKEN }}
    slack_channel: 'your-channel-name'
  env:
    GITHUB_TOKEN: ${{ github.token }}  # Required for accessing GitHub URLs
```

## Example Workflows

### 1. Basic Deployment Notification

```yaml
name: Deploy and Notify
on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      # Your deployment steps here...
      
      - name: Notify Slack
        uses: pal-paul/notify-slack@v1.3.1
        if: always()  # This ensures the notification is sent regardless of previous steps
        with:
          status: ${{ job.status }}
          notify_when: always
          slack_token: ${{ secrets.SLACK_TOKEN }}
          slack_channel: 'deployments'
        env:
          GITHUB_TOKEN: ${{ github.token }}
```

### 2. Notify Only on Failure

```yaml
name: Test and Notify
on:
  push:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      # Your test steps here...
      
      - name: Notify Slack on Failure
        uses: pal-paul/notify-slack@v1.3.1
        if: failure()  # This ensures notification is sent only on failure
        with:
          status: ${{ job.status }}
          notify_when: failure
          slack_token: ${{ secrets.SLACK_TOKEN }}
          slack_channel: 'alerts'
        env:
          GITHUB_TOKEN: ${{ github.token }}
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| status | The status of the job (success/failure) | Yes | success |
| notify_when | When to send notification (always/success/failure) | Yes | always |
| slack_token | Slack token for authentication | Yes | - |
| slack_channel | Slack channel to send notifications to | Yes | - |

## Environment Variables

The action uses environment variables for sensitive information. Make sure to set these in your repository secrets:

- `SLACK_TOKEN`: Your Slack bot token or webhook URL

### Required Environment Variables

- `GITHUB_TOKEN`: Automatically provided by GitHub Actions. Must be passed to the action for proper URL authentication.

### Automatic Environment Variables

The action will automatically use the following GitHub environment variables:

- `GITHUB_WORKFLOW`: Name of the workflow
- `GITHUB_SHA`: The commit SHA
- `GITHUB_REF`: The branch or tag ref
- `GITHUB_REPOSITORY`: The repository name
- `GITHUB_SERVER_URL`: The GitHub server URL
- `GITHUB_RUN_ID`: The unique run identifier
- `GITHUB_JOB`: The job name

## Slack Notification Format

The action will send formatted notifications to your configured Slack channel with the following details:

### Header

- Workflow name with status (succeeded/failed)

### Content

- Status icon (✅ for success, ❌ for failure)
- Workflow name and status
- Run link: Direct link to the workflow run
- Commit: SHA with link to commit details

### Example

```
[Workflow Name] succeeded
✅ workflow-name succeeded
• Run: View workflow run
• Commit: 1234abc

```

## Troubleshooting

### Common Issues

1. "Contribute to..." message in Slack links
   - Make sure you're passing the `GITHUB_TOKEN` in the `env` section of the action
   - This token is required for proper URL authentication

2. Missing workflow information
   - Ensure all required inputs are provided
   - Check if the workflow has proper permissions to access GitHub API
