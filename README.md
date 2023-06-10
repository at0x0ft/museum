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
3. (optional) Make

## Quick Tutorial (Running Flow Example)

0. Install all of the Requirements.
1. Prepare `museum-collections` files for museum templates. For example, check out this repository: [museum-collections](https://github.com/at0x0ft/museum-collections) .
2. Prepare `.env` file to specify the go build config parameters. e.g.
   ```sh
   GOOS='linux'
   GOARCH='amd64'
   ```
3. Run commands below.
   ```sh
   make build
   ```
4. Create `.devcontainer` directory & `skeleton.yml` in the directory. `skeleton.yml` example is below.
   ```yml
   ---
   arguments:
     vscode_devcontainer:
       project_name: test project.
       attach_service: base_shell
     docker_compose:
       project_prefix: museum_dev
       files:
         - ../src/docker-compose.yml
         - ./docker-compose.yml
       vscode_extension_volumes:
         normal: vscode-extensions
         insider: vscode-insider-extensions

   collections:
     path: ../../museum-collections
     list:
       - path: ./base_shell
   ```
5. Run `mix` command giving input and output base `.devcontainer` directory as argument. This command will generate merged multiple `.devcontainer/seed.yml`(s) from `.devcontainer/skeleton.yml`!
   ```sh
   ./bin/museum mix test_project/.devcontainer
   ```
6. Check generated `.devconainer/seed.yml` out! This is a blue print of `devcontainer.json` & `docker-compose.yml` for VSCode Remote Development for Docker. Fix it at your preference.
7. Run `deploy` command giving input and output base `.devcontainer` directory as argument. This command will generate canonical `.devcontainer/devcontainer.json` & `.devcontainer/docker-composer.yml` from `.devcontainer/seed.yml`!
   ```sh
   ./bin/museum deploy test_project/.devcontainer
   ```
8. Check generated `.devconainer/devcontainer.json` & `.devcontainer/docker-compose.yml` out!
