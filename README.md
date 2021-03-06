# wercker_build_trigger
Trigger [Wercker](http://www.wercker.com/) build

[![Go project version](https://badge.fury.io/go/github.com%2Fsue445%2Fwercker_build_trigger.svg)](https://badge.fury.io/go/github.com%2Fsue445%2Fwercker_build_trigger)
[![wercker status](https://app.wercker.com/status/e4c5f1e0f5898b33ffdc26ca29ef4e2c/s/master "wercker status")](https://app.wercker.com/project/byKey/e4c5f1e0f5898b33ffdc26ca29ef4e2c)
[![Coverage Status](https://coveralls.io/repos/github/sue445/wercker_build_trigger/badge.svg?branch=HEAD)](https://coveralls.io/github/sue445/wercker_build_trigger?branch=HEAD)

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
```

* `application_path` : application path **(required)**
  * If wercker application url is https://app.wercker.com/sue445/wercker_build_trigger, `application_path` is `sue445/wercker_build_trigger`
* `pipeline_name` : pipeline name
  *  default is `build`
* `branch` : Branch you want to build
  *  default is `master`

## ProTip
### Weekly build
e.g.) Run a build every Sunday at 3:00

crontab

```
0 3 * * 0 /path/to/wercker_build_trigger --config /path/to/wercker_build_trigger.yml --token xxxxxxx
```

## Changelog
https://github.com/sue445/wercker_build_trigger/releases
