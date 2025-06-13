# notify-slack

A GitHub Action to send Slack notifications for your GitHub workflow status. This action can be easily integrated into any GitHub workflow to provide Slack notifications for your CI/CD pipeline events.

## Usage

To use this action in your GitHub workflow, add the following step to your `.github/workflows/your-workflow.yml` file:

```yaml
- name: Send Slack Notification
  uses: pal-paul/notify-slack@v1  # Use the latest version tag
  with:
    # Required inputs
    status: ${{ job.status }}  # success/failure
    notify_when: 'always'      # always/success/failure
    slack_token: ${{ secrets.SLACK_TOKEN }}
    slack_channel: 'your-channel-name'
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
        uses: pal-paul/notify-slack@v1
        if: always()  # This ensures the notification is sent regardless of previous steps
        with:
          status: ${{ job.status }}
          notify_when: always
          slack_token: ${{ secrets.SLACK_TOKEN }}
          slack_channel: 'deployments'
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
        uses: pal-paul/notify-slack@v1
        if: failure()  # This ensures notification is sent only on failure
        with:
          status: ${{ job.status }}
          notify_when: failure
          slack_token: ${{ secrets.SLACK_TOKEN }}
          slack_channel: 'alerts'
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

The action will automatically use the following GitHub environment variables:

- `GITHUB_WORKFLOW`: Name of the workflow
- `GITHUB_SHA`: The commit SHA
- `GITHUB_REF`: The branch or tag ref
- `GITHUB_REPOSITORY`: The repository name
- `GITHUB_SERVER_URL`: The GitHub server URL
- `GITHUB_RUN_ID`: The unique run identifier
- `GITHUB_JOB`: The job name

## Outputs

The action will send notifications to your configured Slack channel with the following details:

- Workflow name
- Job status (with appropriate emoji)
- Commit information
- Links to:
  - GitHub workflow run
  - Commit details
  - Workflow file
