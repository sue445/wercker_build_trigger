# wercker_build_trigger
Trigger [Wercker](http://www.wercker.com/) build

[![wercker status](https://app.wercker.com/status/e4c5f1e0f5898b33ffdc26ca29ef4e2c/m/master "wercker status")](https://app.wercker.com/project/byKey/e4c5f1e0f5898b33ffdc26ca29ef4e2c)

## Getting

Download latest binary from here

https://github.com/sue445/wercker_build_trigger/releases

## Usage
```bash
# Print version
$ wercker_build_trigger -version

# Trigger a build
$ wercker_build_trigger -token=xxxxx -config=/path/to/config.yml
```

## Options
```bash
$ wercker_build_trigger
  -config string
    	Path to config file
  -token string
    	API token
  -version
    	Whether showing version
```

* `-token` : personal token
  * Go to https://app.wercker.com/profile/tokens and generate your token
* `-config` : path to config file
  * described later

## Config file format
Example

```yaml
pipelines:
  - application_path: "wercker/docs"
    pipeline_name: "build"
    branch: "master"
  - application_path: "sue445/wercker_build_trigger"
    pipeline_name: "build"
    branch: "master"
```

* `application_path` : application path **(required)**
  * If wercker application url is https://app.wercker.com/sue445/wercker_build_trigger, `application_path` is `sue445/wercker_build_trigger`
* `pipeline_name` : pipeline name **(required)**
  * e.g.) build
* `branch` : Branch you want to build
  *  default is `master`
