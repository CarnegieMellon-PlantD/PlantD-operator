# Contributing

We, the maintainers, love pull requests from everyone, but often find
we must say "no" despite how reasonable the proposal may seem.

For this reason, we ask that you open an issue to discuss proposed
changes prior to submitting a pull request for the implementation.
This helps us to provide direction as to implementation details, which
branch to base your changes on, and so on.

1. Open an issue to describe your proposed improvement or feature
1. Create a fork of this repository
1. Create a new feature branch in the forked repo with an appropriate name
1. If applicable, add a [CHANGELOG-7.md entry](#changelog) describing your change.
1. Create a Pull Request into the main repo as appropriate based on the issue discussion

PlantD is and always will be open source, and we continue to highly
value community contribution. 

## Changelog


All new changes go underneath the _Unreleased_ heading at the top of the Changelog.
Beyond that, here are some additional guidelines that should make it more clear where your
change goes in the Changelog.

### Added

Any _new_ functionality goes here. This may be a new field on a data type or a new data
type altogether; a new API endpoint; or possibly a whole new feature. In general, these
are sentences that start with the word "added."

Examples:

- `begin` field to silences that initiates silencing at a given timestamp
- /healthz endpoint that reports health of the prometheus-agent process

### Changed

Changes to any existing component or functionality of the system that does not cause
breaking changes to users or developers go here. _Changed_ is distinguishable from
_Fixed_ in that it is an intentional change to existing functionality.

Examples:

- Refactored the API to use reusable controller logic

### Fixed

Fixed bugs go here.

Examples:

- Don't delete auth tokens at startup

### Deprecated

Deprecated should include any soon-to-be removed functionality. An entry here that
is user facing will likely yield entries in _Removed_ or _Breaking_ eventually.

Examples:

- The /health API endpoint is being replaced by /healthz on the backend
- The /stash API endpoint is being removed in a future release

### Removed

Removed is for the removal of functionality that does not directly impact users,
these entries most likely only impact developers of PlantD.  If user facing
functionality is removed, an entry should be added to the _Breaking Changes_
section instead.

Examples:

- Removed references to `encoding/json` in favor of `json-iter`.
- Removed unused `Store` interface for `BlobStore`.

### Security

Any fixes to address security exploits should be added to this section. If
available, include an associated CVE entry.

Examples:

- Upgraded build to use Go 1.9.1 to address [CVE-2017-15041](https://www.cvedetails.com/cve/CVE-2017-15041/)
- Fixed issue where users could view entities without permission

### Breaking Changes

Whenever you have to make a change that will cause users to be unable to
upgrade versions of PlantD without intervention by an operator, your change
goes here. Try to avoid these. If they're required, we should have documented
justification in a GitHub issue and preferably a proposal. We should also bump
minor versions at this time.

Examples:

- Refactored how Checks are stored in Etcd, `plantd-operator` docker image is required to upgrade

## Git Workflow

Our git workflow is largely inspired by [GitHub Flow](https://guides.github.com/introduction/flow/) and [Oneflow](https://www.endoflineblog.com/oneflow-a-git-branching-model-and-workflow) but adapted to our reality and our needs.

Here are the highlights:
- There's only one eternal branch named `main`. All other branches are temporary.
- Feature branches are where the day-to-day development work happens. They are based from main and pushed continuously back into it whenever possible so the pull requests are small and simple, while keeping main stable.
- Release branches are branched off from main at the point all the necessary features are present. From then on, new work aimed for the next release is pushed to main as always, while any necessary changes for the release (updating the changelog, last minute bugfixes, updating dependencies etc.) are pushed to the release branch. Once the release is ready, we tag the top of the release branch. Finally, we merge the release branch into main.
- Hotfixes are very similar to releases, except we branch off from a release tag. A hotfix is basically an immediate fix for something that's really getting in the way of our users.
