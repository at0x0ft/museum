# museum

VSCode docker remote development environment template collections.

## Note

**This project has been under development yet**. The specifications would be **changed disruptively** in the future.

## Purpose (Focus? Goal?)

1. Splitting environments between development base shell container and runtime containers.
2. Establishing high configurable container development environment.
3. Managing extensions with code-base and in the container sandbox.
4. Absorbing the differences between between host machines (e.g. Linux, MacOS, Windows...) and ensuring portability & reproducibility.

## Requirements

1. VSCode
2. Docker (with Docker Buildkit)
3. Docker Compose
4. [yq](https://github.com/mikefarah/yq)

## Quick Tutorial (Running Flow Example)

0. Install all of the Requirements.
1. Create `.devcontainer` directory & `skeleton.yml` in the directory. `skeleton.yml` example is below.
   ```yml
   ---
   version: "0"

   base_shell:
     path: ./services/base_shell
   ```
2. Run `restore` command giving input and output base `.devcontainer` directory as argument. This command will generate `.devcontainer/seed.yml` from `.devcontainer/skeleton.yml`!
   ```sh
   docker-compose run --rm museum restore test_project/.devcontainer
   ```
3. Check generated `.devconainer/seed.yml` out! This is a blue print of `devcontainer.json` & `docker-compose.yml` for VSCode Remote Development for Docker. Fix it at your preference.
4. Run `deploy` command giving input and output base `.devcontainer` directory as argument. This command will generate canonical `.devcontainer/devcontainer.json` & `.devcontainer/docker-composer.yml` from `.devcontainer/seed.yml`!
   ```sh
   docker-compose run --rm museum deploy test_project/.devcontainer
   ```
5. Check generated `.devconainer/devcontainer.json` & `.devcontainer/docker-compose.yml` out!
