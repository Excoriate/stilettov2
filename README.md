<h1 align="center">
  Stiletto ⚔️
</h1>
<p align="center">Portable, and <b>containerized</b> pipelines that runs everywhere! <b> works on top of Dagger ❤️️</b>.<br/><br/>

---
![GitHub release (latest by date)](https://img.shields.io/github/v/release/Excoriate/Stiletto) [![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md) [![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)


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
Just executes:
```bash
stiletto

```




## Roadmap 🗓️

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
