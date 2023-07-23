<h1 align="center">
  Stiletto ⚔️
</h1>
<p align="center">Portable, and <b>containerized</b> pipelines that runs everywhere! <b> works on top of Dagger ❤️️</b>.<br/><br/>

---
![GitHub release (latest by date)](https://img.shields.io/github/v/release/Excoriate/Stilettov2) [![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md) [![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)[![Run golangci-lint](https://github.com/Excoriate/stilettov2/actions/workflows/go-ci-lint.yaml/badge.svg)](https://github.com/Excoriate/stilettov2/actions/workflows/go-ci-lint.yaml)[![Go Unit Test](https://github.com/Excoriate/stilettov2/actions/workflows/go-ci-test.yml/badge.svg)](https://github.com/Excoriate/stilettov2/actions/workflows/go-ci-test.yml)


---
Stiletto (means "dagger" in Italian) is a portable, and containerized pipelines that runs everywhere, since it works on top of the wonderful
[Dagger](https://dagger.io). Its main motivation is to provide a simple way to run pipelines in a portable way, and also to provide a way to run
pipelines in a containerized way, so you can run your pipelines in any environment, and also in any CI/CD platform.
Stiletto follows the same principles as [GitHub Actions](https://github.com/features/actions) in the way that it defines a pipeline, its jobs and actions; nevertheless, **it's not a 'portable' version of GH actions**.

## 🔧 How to Install Stiletto

Stiletto provides binary distributions for every release which are generated using GoReleaser. To install it, you can use the pre-built binaries which are available for Linux, Windows, and macOS:

1. Navigate to the [Releases](https://github.com/Excoriate/stilettov2/releases) page.
2. Download the archive for your respective OS and architecture.
3. Extract the archive.
4. Move the `stiletto` binary into your `$PATH`.

Or, based on your OS. For MacOS, you can use [Homebrew](https://brew.sh/):

```bash
brew tap Excoriate/homebrew-tap https://github.com/Excoriate/homebrew-tap.git
brew install stiletto
```
>**NOTE**: There are compiled binaries available for most of the common platforms, including Windows. Check the
[Releases](https://github.com/Excoriate/stilettov2/releases) page.



## ▶️ How to Use Stiletto
### Core concepts
* 🤖 **Runner**: It's how the tasks and jobs are executed. Currently, the only runner supported is [Dagger](https://dagger.io).
* ⚡️ **Task**: It's the smallest unit of work that can be executed. It's composed by a set of `commands`. If you're familiar with GitHub actions, it's equivalent to the [steps](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#jobsjob_idsteps).
* 📦 **Job**: It's a set of tasks that are executed in a given order. If you're familiar with GitHub actions, it's equivalent to the [jobs](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#jobs).
* 📜 **manifest**: It's the file that defines the pipeline. It's a YAML file that contains the definition of the jobs, tasks or workflows. The specs are defined in the [manifests](./docs/manifests) folder.

### How to define a manifest
A manifest can be defined manually following the examples available in the [examples](./examples) folder
>**NOTE**: In the future, a CLI `command` will be added in order to automatically generate all the supported manifests.

Here's an example of a manifest that defines an IAC (infrastructure-as-code) task for Terragrunt (which works on top of terraform):
```yaml
---
apiVersion: v1
kind: Task
metadata:
    name: iac-terragrunt
spec:
    containerImage: alpine/terragrunt
    mountDir: .
    workdir: examples/terragrunt
    commandsSpec:
        - binary:
          commands:
              - ls -ltrah /mnt
        - binary: terragrunt
          commands:
              - init
              - plan
              - apply -auto-approve
              - destroy -auto-approve

```
Some of the tasks options and capabilities while being defined are:
* It can scan **environment variables** using the following options:
  * Scan `AWS` env vars out of the box.
  * Scan `terraform` (`TF_VARS_`) env vars out of the box.
  * Scan all the host environment variables if available.
  * Scan selectively environment variables, or set them explicitly.
* It can mount **directories** and work on top of them defining **workdir** as an independent option.
* It can define **commands** as _plain strings_, _Stiletto_ will take care of ensuring that the commands are executed in the right order.

### CLI
Stiletto provides a CLI that can be used to run the pipelines. Just run `stiletto help` to see the available commands. However, here there are some examples of how to use it:
- Running a task from a `taskfile`:
```bash
stiletto job dagger --task-files=mytasks/my-task.yaml
```
- Running a task from a `taskfile` and overriding the `workdir`:
```bash
stiletto job --mountdir=/tmp --workdir=/tmp --task-files=mytasks/my-task.yaml
```

## Roadmap 🗓️

- [ ] New `manifest` command to generate manifests.
- [ ] New types for manifest (e.g. `workflow`, `job`).
- [ ] Enable workflows (`workflow.yml`) for more complex pipelines.
- [ ] Cover necessary/critical parts of Stiletto with proper unit tests.
- [ ] Add an official DockerFile that can be available in [DockerHub](https://hub.docker.com/).
- [ ] Add an API to trigger Stiletto for automatic pipelines through _internal developer portals_ or other interfaces.

>**Note**: This is still work in progress, however, I'll be happy to receive any feedback or contribution. Ensure you've read the [contributing guide](./CONTRIBUTING.md) before doing so.


## Contributing

Please read our [contributing guide](./CONTRIBUTING.md).

## Community

Find me in:

- 📧 [Email](mailto:alex@ideaup.cl)
- 🧳 [Linkedin](https://www.linkedin.com/in/alextorresruiz/)


<a href="https://github.com/Excoriate/stilettov2/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=Excoriate/stiletto" />
</a>
