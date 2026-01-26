A build tool for DayZ mods.  This is a work in progress.

## Releasing

You will need an appropriate GitHub personal access token, set to the environment variable `GITHUB_TOKEN`.

```
git tag -a v<semver> -m "<release note>"
git push origin v<semver>
goreleaser release
```
