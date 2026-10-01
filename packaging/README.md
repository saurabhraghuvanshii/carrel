# Packaging

Templates for a Homebrew tap and a Scoop bucket. Nothing here is published yet.

`carrel.rb` and `carrel.json` hold `{{version}}` and `{{sha_<os>_<arch>}}` placeholders. After a release, fill them from that release's `checksums.txt`:

```
gh release download v0.1.0 --pattern checksums.txt --dir /tmp/carrel-0.1.0
sh packaging/fill.sh 0.1.0 /tmp/carrel-0.1.0/checksums.txt /tmp/carrel-0.1.0/out
```

## Homebrew tap

1. Create a public repository named `homebrew-tap` under the same GitHub owner.
2. Copy the filled `carrel.rb` to `Formula/carrel.rb` in that repository and push.
3. Users then run `brew install saurabhraghuvanshii/tap/carrel`.
4. Check it locally first with `brew install --build-from-source ./Formula/carrel.rb` and `brew test carrel`.

## Scoop bucket

1. Create a public repository named `scoop-bucket` under the same GitHub owner.
2. Copy the filled `carrel.json` to `bucket/carrel.json` and push.
3. Users then run `scoop bucket add carrel https://github.com/saurabhraghuvanshii/scoop-bucket` and `scoop install carrel`.
4. `checkver` and `autoupdate` let Scoop's update tooling pick up later releases on its own.

Repeat step 2 of each for every release, or automate it later with a job in `release.yml` that pushes to both repositories with a token.
