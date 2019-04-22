# Recording Storage

## Introduction

This service works with recordings doing different jobs on them:
* transfer recordings files from AWS S3 to other storages (Dracoon etc.)

### Transferring from Amazon to Dracoon

The following steps are performed during transferring:
* Parse queue message
  * Get source file (bucket name and key)
  * Get destination folders (Dracoon parent room and path's templates)
    - Templates consist of the following variables: `$instance` and `$meetingType`
* Check if source file exists
* Check if destination rooms (defined by parent room and paths generated from templates).
  - If no, the missed room is created
* Start transferring
  - for Dracoon it means that upload channel is opened
* Downloading file from S3 by chunks
* Encrypt every chunk with hybrid crypto algorithm (using symmetric cryptography for file encryption and asymmetric cryptography for file's key encryption)
* Upload encrypted chunk to Dracoon room
* Complete transferring
  - for S3 it means that downloaded file removed from there
  - for Dracoon it means that upload channel is closed and encryption key is uploaded
* Share access to the uploaded file between all room's users who can read in that room
* Also access sharing for newly created users is done hourly

## Queues of messages

The services uses AWS SQS to get messages that describe jobs for service's workers

### Message structure

```json
{
    "storage": {
        "type":"dracoon",
        "extra": {
            "dracoon": {
                "targets": [
                    {
                        "path":"some room\/sub-room",
                        "expiration":"4 years 3 months 2 weeks 1 day"
                    }
                ]
            },
            "aws.s3": {
                "bucket": "24s-public",
                "keys": [
                    "sandbox\/test-video.mp4"
                ]
            }
        }
    },
    "meeting": {
        "id": 22980,
        "type": {
            "name": "live now",
            "duration": 60
        },
        "user": {
            "name": "Sessions support",
            "email": "support@24sessions.com"
        },
        "guest": {
            "name": "Guest",
            "email": "guest@24s.co"
        },
        "date": "2017-09-11T20:24:43+00:00",
        "description": "natus aperiam quia blanditiis dignissimos nam porro."
    },
    "instance": {
        "domain": "24slocal.new.com"
    }
}
```

Explanations:
* `storage.extra.dracoon.targets` - list of targets on Dracoon side where recording should be uploaded. Consist of the following parameters:
  * `path` - path to the room (relative, in the parent room)
  * `expiration` - string that describe when uploaded file is expired beginning from now. Each path can have different expiration or no expiration at all.
* `storage.extra.aws.s3` - describe where recording is now (where it should be downloaded from). Consists of the following parameters:
  * bucket - name of AWS S3 bucket
  * keys - list (array) of AWS S3 objects (sort of paths to files). There can be more than one key: for video, phone etc.
* `instance.domain` - domain of the instance where recording was made. Necessary to get instance configuration.

**Expiration format**

It can be any combinations of phrases `<number> <period>` where period is one of **year**(s), **month**(s), **week**(s) or **day**(s). Example: "4 years 3 months 2 weeks 1 day".

## New parameters in instance configuration

* `recordings.storage` - type of storage. Currently only **dracoon**
* `recordings.dracoon.*` - subset of Dracoon configuration for the instance
  * `recordings.dracoon.parent-room-id` - room ID where all sub-rooms with recordings of this instances are stored
  * `recordings.dracoon.login` - part of Dracoon credentials
  * `recordings.dracoon.password` - part of Dracoon credentials
  * `recordings.dracoon.encryption-password` - part of Dracoon credentials
  * `recordings.dracoon.targets` - describes where to upload each recording
* `wrapup.*` - for ING where wrapups are not used we use the wrapup input field to allow them to enter their inner customer ID
  * `wrapup.label` - set custom label for wrapup input field (that appears after meeting is completed)
  * `wrapup.send-notification` - send email about entered wrapup or not. For ING case - no.

**Dracoon target's format**

It's stored in instance configuration in JSON format

```json
    [{
        "meetingTypeId": 14,
        "copyTo": [{
          "path": "some room/sub-room", "expiration": "4 years 3 months 2 weeks 1 day"
        }]
    }]
```

So for every meeting type we set here its ID and list of Dracoon rooms (paths to them). For each room the expiration can be set.
Only recordings of meetings of the given meeting type will be transferred to Dracoon.

**Tip**. As meeting type ID you can use special value `*any*` that means that the given path and expiration will be applied for recording of meeting of any meeting type.

## Dependencies

See file `Gopkg.toml`

**To add new dependency**

* Add new dependency via `dep ensure -add <path-to-repository`
* For our libs use constraint with the necessary branch/version
* If new dependency has configuration then:
  * Add configuration file to `conf` directory and to `AWS S3` bucket
  * Add downloading of configuration to `entrypoint.sh` file

## Installation

* Clone repository by running the following commands:
  * `git clone git@gitlab.com:24sessions/service-recordingStorage.git recording-storage`
  * `cd service-recordingStorage`
* Install dependencies.
  * install `dep` if not exist:
    `curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh`
  * run `dep ensure`
  * run `git install -v`
* Add configuration files.
    - `cp conf/app.conf.default app.conf`
    - `cp conf/pubsub/sqs/topics.yml.dev conf/pubsub/sqs/topics.yml`
* Test project if necessary.
  * add test configuration:
    - `mkdir tests/conf`
    - `cp conf/app.conf.test tests/conf/app.conf`
    - `mkdir -p tests/conf/pubsub/sqs`
    - `cp conf/pubsub/sqs/topics.yml.test tests/conf/pubsub/sqs/topics.yml`
  * run `go test ./tests/...`

## Run service

**On test environment**

* Add entry in Git CI configuration in `AWS S3`, in `docker-compose.yml` file

**On production**

* Add entry in Git CI configuration in `AWS S3`, in `docker-compose.yml` file
* Set right value for `aws_queue_calendar_url` in Website parameters.yml

## API

### Healthcheck

Summary: Checks application status: AWS SQS, AWS S3, Dracoon, etc. Must be used by monitor.

Endpoint: /healthcheck

Method: **GET**

Auth header: **No Auth**

Example: GET: /healthcheck

Response, status **200**

```json
{
    "status": "success"
}
```

Response, status **500**

```json
{
    "error": "worker_error_file_transferer",
    "status": "error"
}
```

## Configurations

### app.conf

* `logger.*` Subset of parameters for logger
* `configurator.url` URL to our internal configurator (that manage instances configuration)
* `workers.FileTransferer.*` Subset of parameters for FileTransferer worker
  - `workers.FileTransferer.count` Count of workers to transfer files. For prod - minimum 1
  - `workers.FileTransferer.reportLevel`
  - `workers.FileTransferer.queue` Key name from **conf/pubsub/sqs/topics.yml** (in that file URLs of all AWS SQS queues are stored)
  - `workers.FileTransferer.source.type` Source storage. Currently only **s3**
  - `workers.FileTransferer.destination.type` Destination storage. Currently only **dracoon**
* `workers.DracoonGlobalSharer.*` Subset of parameters for DracoonGlobalSharer worker
  - `workers.DracoonGlobalSharer.count` Count of workers to share access for Dracoon files, maximum 1
  - `workers.DracoonGlobalSharer.jobDelayMinutes` delay in minutes between subsequent runnings of the worker, minimum 1 minute.
* `s3.*` Subset of parameters for AWS S3
  - `s3.chunkSize` Size of downloading file's chunks in bytes. Default value is 20 Mb.
  - `s3.removeAfterTransferring` Whether to remove file when transferring is completed
* `dracoon.*` Subset of parameters for Dracoon
  - `dracoon.baseUrl` Base URL of Dracoon API. Currently https://dracoon.team/api/v4
  - `dracoon.login`, `dracoon.password` Credentials for our account. On production they are stored in instance configuration. So this ones are only for test purpose
  - `dracoon.encryptionPassword` Password to decrypt private key of our Dracoon account. On production it's stored in instance configuration. So this one is only for test purpose
  - `dracoon.parentRoomId` Parent room's ID. On production it's stored in instance configuration. So this one is only for test purpose
  - `dracoon.test.*` The same as `dracoon.*` but only for test purpose

## Worker domain

* FileTransferer: transfer instance recordings from one storage to another (e.g., AWS S3, Dracoon etc.)
* DracoonGlobalSharer: share access for all Dracoon users to all files to which users have access rights and which are not accessible now.

## Description Codes

Recording Storage : 15000 - 15999

* 15000: Info:      Start application
* 15001: Info:      Start CLI application
* 15002: Alert:     Application (or worker) can not load config for workers
* 15003: Critical:  Amazon session error
* 15004: Critical:  Amazon configuration error
* 15005: Error:     SQS error
* 15006: Info:      transferer status
* 15007: Error:     transferer error
* 15008: Info:      Worker profiling
* 15009: Error:     AmazonToDracoon error
* 15010: Error:     REST client error
* 15011: Error:     Storage error
* 15012: Info:      Storage status
* 15013: Error:     InstanceJobRunner error
* 15014: Error:     HealthCheck error
* 15015: Error:     Worker error
* 15016: Error:     DracoonGlobalSharer error
* 15017: Info:      DracoonGlobalSharer status
* 15018: Error:     Configuration error
