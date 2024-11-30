# Take home assignment Submission
This is a submission for the take home assignment, client file contains a simple VueJS app that allows you to communicate with the app.<br />
All endpoints were created simple as the app itself should not be over complicated, the app saves running jobs in a Jobhub and backs up all historical job info in the DB.<br />
If `TH_HOOK_URL` is set, whenever the status of a job changes, it will fire a GET request to the url set in the env variable, the jobId and status values will be added as query strings to the url.

## Requirements
Since the app uses sqlite by default, it does require gcc in order to run the app, in addition it requires you to set environment variable`CGO_ENABLED=1` in order to run the app.

## Environment Variables
| Variable | Env  | Default  |
|---|---|---|
| Port  | TH_PORT  | 8080  |
| Token  | TH_TOKEN  | 123  |
| HookUrl  | TH_HOOK_URL  |  |

## API Endpoints
In order to use the api endpoints you need to set `Authorization` header to the value set in the "token" env value (123 by default).<br />

- `/jobs` - POST creates a job, request body is as follows
```
{
    "command": "python3",
    "args": ["long_task.py"]
}
```

- `/jobs` - GET Gets all jobs.
- `/jobs/{id}` - GET job by  ID.
- `/jobs/{id}/status` - GET returns the job status.
- `/jobs/{id}/logs` - GET returns the job logs, query string value `lines` can be set to return the latest number of lines in the log file.
- `/jobs/{id}/kill` - GET kills the job.