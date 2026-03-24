# Auto Review (`cmd/autoreview`)

AI-powered code review tool. It analyzes the latest commit or your current
working directory changes using Google's Gemini API and provides a summary and
an LGTM status.

## Usage

Build and run the tool:

```bash
bazel run //cmd/autoreview
```

Display a help message:

```bash
bazel run //cmd/autoreview -- --help
```

Pass command line options:

```bash
bazel run //cmd/autoreview -- [options]
```

## Authentication Setup

The tool uses Google Cloud application-default credentials.
Run the following commands to authenticate via Google Cloud:

```bash
gcloud auth application-default set-quota-project skia-infra-corp
gcloud auth application-default login
```

See more go/gemini-api/authentication.
