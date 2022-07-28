# VSCode docker development environments

VSCode docker remote development environment collections.

## Purpose (Focus? Goal?)

1. Splitting environments between development base shell container and runtime containers.
2. Establishing high configurable container development environment.
3. Managing extensions with code-base and in the container sandbox.
4. Absorbing the differences between between host machines (e.g. Linux, MacOS, Windows...) and ensuring portability & reproducibility.

## Requirements

1. Docker (with Docker Buildkit)
2. Docker Compose
3. [yq](https://github.com/mikefarah/yq)
