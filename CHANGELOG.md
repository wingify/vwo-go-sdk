# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.8.0] - 2020-08-05

### Changes
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
