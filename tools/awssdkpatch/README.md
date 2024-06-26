# awssdkpatch

`awssdkpatch` creates [gopatch](https://github.com/uber-go/gopatch) patches that aid in the migration to AWS SDK v2.
For more context on why this is required, visit the [AWS SDK migration meta-issue](https://github.com/hashicorp/terraform-provider-aws/issues/32976) on Github.

## Options

```console
awssdkpatch -h
```

```
Generate a patch file to migrate a service from AWS SDK for Go V1 to V2.

Usage: awssdkpatch [flags]

Flags:
  -out string
        output file (optional) (default "awssdk.patch")
  -service string
        service to migrate (required)
```

## Usage

For most cases, the `awssdkpatch` executable is called as follows:

```console
awssdkpatch -service xray
```

This generates a patch file, `awssdk.patch`, in the root of the project.
This patch can then be applied to the appropriate service directory with `gopatch`:

```console
gopatch -p awssdk.patch internal/service/xray/...
```

To preview the changes in stdout without modifying files, include the `-d/--diff` flag:

```console
gopatch -d -p awssdk.patch internal/service/xray/...
```
