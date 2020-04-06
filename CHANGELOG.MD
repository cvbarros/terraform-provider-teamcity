# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0]

This release is the first major release and includes an upgrade to TeamCity 2019.2.2 as the supported version.
Whenever a major version is bumped, backwards compatibility with previous TeamCity versions aren't guaranteed.

### Features:
- **New resource**: Support for TeamCity User Group [#69]

### Fixes:
- Resolved [#43] Migrated _step_ attribute for Build Config to `schema.TypeList` [#70]

### Maintenance:
- Upgrade client library to `v1.0.0` and TeamCity in tests to 2019.2.2 [#71]
- Migrate to `terraform-plugin-sdk` instead of Terraform Core [#72]

[//]: # (Release links)
[1.0.0]: https://github.com/cvbarros/terraform-provider-teamcity/releases/tag/v1.0.0

[//]: # (Issue/PR links)
[#43]: https://github.com/cvbarros/terraform-provider-teamcity/issues/43
[#69]: https://github.com/cvbarros/terraform-provider-teamcity/pull/69
[#70]: https://github.com/cvbarros/terraform-provider-teamcity/pull/70
[#71]: https://github.com/cvbarros/terraform-provider-teamcity/pull/71
[#72]: https://github.com/cvbarros/terraform-provider-teamcity/pull/72
