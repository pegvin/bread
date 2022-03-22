# Bread ![:bread:](./.github/bread.svg)

Install, update, remove & run AppImage from GitHub using your CLI. (Fork of [AppImage ClI Tool](https://github.com/AppImageCrafters/appimage-cli-tool))

## Features
- Install from the GitHub Releases
- Run Applications From Remote Without Installing Them
- Update with ease

## Installation

If you already have installed or installing [LibAppImage](https://repology.org/project/libappimage/versions) you can just install bread binary with Curl:

```bash
sudo curl -L https://github.com/DEVLOPRR/bread/releases/download/v0.4.2/bread-0.4.2-x86_64 -o /usr/local/bin/bread && sudo chmod +x /usr/local/bin/bread
```

<br>

If LibAppImage is not available for your distribution or you don't want to install it specially for this software, you can get the bread's AppImage which contains LibAppImage v1.0.2, using Curl:
```bash
sudo curl -L https://github.com/DEVLOPRR/bread/releases/download/v0.4.2/bread-0.4.2-x86_64.AppImage -o /usr/local/bin/bread && sudo chmod +x /usr/local/bin/bread
```

***It is recommended to install the Binary instead of AppImage, if your distribution provides LibAppImage Version 1.0.2 Or Higher.***

---

## Uninstallation

Just Remove the binary
```bash
rm -v /usr/local/bin/bread
```

**NOTE** this won't delete the app you've installed.

---

## Usage

<details>
  <summary>NOTE</summary>
  <br>
  <p>Often there are many times when the GitHub user and repo both are same, for example <a href="https://github.com/LibreSprite/LibreSprite" target="_blank">libresprite</a>, so in this case you can just specify single name like this <code>bread install libresprite</code>, this works with all the commands</p>
</details>

<details>
  <summary>Install a application</summary>
  <br>
  <p>To install an Application from GitHub you can use the install command where user is the github repo owner and repo is the repository name</p>
  <pre><code>bread install user/repo</code></pre>
</details>

<details>
  <summary>Run a application from remote</summary>
  <br>
  <p>If you want to run a application from remote without installing it you can use the run command</p>
  <pre><code>bread run user/repo</code></pre>
  <p>You can pass CLI arguments to the application too like this</p>
  <pre><code>bread run user/repo -- --arg1 --arg2</code></pre>
  <p>You can clear the download cache using clean command <code>bread clean</code>, Since all the applications you run from remote are cached so that it isn't downloaded everytime</p>
</details>

<details>
  <summary>Remove a application</summary>
  <br>
  <p>you can remove a installed application using the remove command</p>
  <pre><code>bread remove user/repo</code></pre>
</details>

<details>
  <summary>Update a applicationn</summary>
  <br>
  <p>You can update a application using the update command</p>
  <pre><code>bread update user/repo</code></pre>

  <p>if you just want to check if update is available you can use the <code>--check</code> flag</p>
  <pre><code>bread update user/repo --check</code></pre>

  <p>if you want to update all the applications you can use the <code>--all</code> flag</p>
  <pre><code>bread update --all</code></pre>

  <p>the <code>--check</code> & <code>--all</code> flag can be used together</p>
  <pre><code>bread update --all --check</code></pre>
</details>

<details>
  <summary>Search for an application</summary>
  <br>
  <p>You can search for a application from the [AppImage](https://appimage.github.io) API</p>
  <pre><code>bread search "Your search text"</code></pre>
</details>

<details>
  <summary>List all the installed application</summary>
  <br>
  <p>You can list all the installed applications using list command</p>
  <pre><code>bread list</code></pre>
</details>

---

## Tested On:
- Ubuntu 20.04 - by me
- Arch Linux - by [my brother](https://github.com/idno34)

---

## File/Folder Layout
Bread installs all the applications inside the `Applications` directory in your Linux Home Directory `~`, inside this directory there can be also a directory named `run-cache` which contains all the appimages you've run from the remote via the `bread run` command.

In the `Applications` there is also a file named `.registry.json` which contains information related to the installed applications!

---

## Building From Source

Make Sure You Have Go version 1.18.x & [AppImage Builder](https://appimage-builder.readthedocs.io/en/latest/) Installed.

Get The Repository Via Git:

```bash
git clone https://github.com/DEVLOPRR/bread
```

Go Inside The Source Code Directory & Get All The Dependencies:

```bash
cd bread
go mod tidy
```

Make The Build Script Executable And Run It

```bash
chmod +x ./make
./make
```

And To Build The AppImage Run

```bash
./make appimage
```

---
## Build Script
The `make` bash script can build your go code, make appimage out of it, and clean the left over stuff including the genrated builds.

#### Building in Development Mode
This will build the go code into a binary inside the `build` folder
```
./make
```

#### Building in Production Mode
Building for production requires passing `--prod` flag which will enable some compiler options resulting in a small build size.
```
./make --prod
```

#### Building the AppImage
Bread requires libappimage0 for integrating your apps to desktop, which is done via libappimage, to make End user's life easier we package the libappimage with bread and that's why we build the binaries into AppImages so that user doesn't need to install anything.

To make a appimage out the pre built binaries
```
./make appimage
```

#### Get Dependency
To install the dependencies require to build go binary
```
./make get-deps
```

---

## Todo
- [ ] Improve UI
- [x] Make AppImages Runnable From Remote Without Installing (Done in v0.3.6)
- [x] Work On Reducing Binary Sizes (Reduced From 11.1MB to 3.1MB)
- [ ] Add 32 Bit Builds
- [ ] Add Auto Updater Which Can Update The Bread Itself
- [x] Add `--version` To Get The Version (Done in v0.2.2)
- [ ] Mirrors:
  - I Would Like To Introduce Concept Of Mirror Lists Which Contain The List Of AppImages With The Download URL, tho currently i am not working on it but in future i might.

---

# Thanks
