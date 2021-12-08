# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.16.0] - 2021-10-08

- Sending stats which are used for launching the SDK like storage service, logger, and integrations, etc. in tracking calls(track-user and batch-event). This is solely for debugging purpose. We are only sending whether a particular key(feature) is used not the actual value of the key

## [1.14.0] - 2021-06-17

### Added

- Webhooks support. Introduced new API `GetAndUpdateSettingsFile` to fetch and update settings-file in case of webhook-trigger
- Event Batching
  - Added support for batching of events sent to VWO server
  - Added `FlushEvents` API to manually flush the batch events queue when batch_events config is passed. Note: batch_events config i.e. events_per_request and request_time_interval won't be considered while manually flushing
  - If `RequestTimeInterval` is passed, it will only set the timer when the first event will arrive
  - If `EventsPerRequest` is provided, after flushing of events, new interval will be registered when the first event will arrive

  ```go
  callBack := func(err error, batch []map[string]interface{}) {
    log.Println(fmt.Sprintf("Batch events pushed %v %v", batch, err))
  }

  instance, err = vwo.Launch(
    settingsFile,
    api.WithBatchEventQueue(
      api.BatchConfig{RequestTimeInterval: 10, EventsPerRequest: 3},
      callBack
    )
  )

  // Manually flush events
  instance.FlushEvents()
  ```

- Integrations
  - Exposed lifecycle hook events. This feature allows sending VWO data to third party integrations.

  ```go
  Integrations := func(integrationsMap map[string]interface{}) {
    fmt.Println("Integrations Map : ", integrationsMap)
  }

  instance, err = vwo.Launch(settingsFile, api.WithIntegrationsCallBack(Integrations))
  ```

### Changed

- Send environment token in every network call initiated from SDK to the VWO server. This will help in viewing campaign reports on the basis of environment.
- Removed sending user-id, that is provided in the various APIs, in the tracking calls to VWO server as it might contain sensitive PII data.
- SDK Key will not be logged in any log message, for example, tracking call logs.

## [1.8.0] - 2020-08-05

### Changed
- Update track API to handle duplicate and unique conversions and corresponding changes in `Launch` API
- Update track API to track a goal globally across campaigns with the same `goalIdentififer` and corresponding changes in `Launch` API
- Update user storage to store `goalIdentififer`
- Handled new attributes, `GoalTypeToTrack` and `ShouldUserReturningUser` in `Options` and `Launch` API

```go
// It will track goal having `goalIdentifier` of campaign having `CampaignKey` for the user having `userId` as id.
campaignKeys = "CampaignKey"
isSuccessful = vwoInstance.Track(campaignKeys, goalIdentifier, userId, options);

// it will track goal having `goalIdentifier` of campaigns having `CampaignKey1` and `CampaignKey2` for the user having `userId` as id.
campaignKeys = []string{"CampaignKey1", "CampaignKey2"}
isSuccessful = vwoInstance.Track(campaignKeys, goalIdentifier, userId, options);

// it will track goal having `goalIdentifier` of all the campaigns
campaignKeys = nil
isSuccessful = vwoInstance.Track(campaignKeys, goalIdentifier, userId, options);

//Read more about configuration and usage - https://developers.vwo.com/reference#server-side-sdk-track
```

## [1.0.0] - 2020-06-03

### Added

- First release with server-side A/B capabilities.
- Feature rollout and Feature-test support
- Campaigns support pre-segmentation as well as post-segmentation
- Forced Variation aka whitelisting is also available in A/B and Feature Test
- Tests for ensuring code coverage
