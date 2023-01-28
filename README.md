
# Bread ![:bread:](./.github/bread.svg)

Install, update, remove & run AppImage from GitHub using your CLI. (Fork of [AppImage ClI Tool](https://github.com/AppImageCrafters/appimage-cli-tool))

## Features
- Install from the GitHub Releases
- Automatically Integrate App To Desktop When Installing/Updating
- Run Applications From Remote Without Installing Them
- Update with ease

## Getting Started

### Installation

<details>
  <summary>Arch Linux & it's Derivatives</summary>
  <br>
  <p>you can use this step if your distribution does provide <code>libappimage</code> v1.0.0 or greater, which is the case on Arch Linux & it's Derivatives, kaOS, KDE Neon, Parabola Linux</p>
  <p>install <code>libappimage</code> dependency</p>
  <pre><code>pacman -S libappimage</code></pre>
  <p>then install bread</p>
  <pre><code>sudo curl -L https://github.com/pegvin/bread/releases/download/v0.7.2/bread-0.7.2-x86_64 -o /usr/local/bin/bread && sudo chmod +x /usr/local/bin/bread</code></pre>
</details>

<details>
  <summary>Debian & it's Derivatives</summary>
  <br>
  <p>you can use this step if your distribution doesn't provide <code>libappimage</code> v1.0.0 or greater, which is the case on Debian & it's derivatives</p>
  <p>get the appimage containing <code>libappimage</code> v1.0.3</p>
  <pre><code>sudo curl -L https://github.com/pegvin/bread/releases/download/v0.7.2/bread-0.7.2-x86_64.AppImage -o /usr/local/bin/bread && sudo chmod +x /usr/local/bin/bread</code></pre>
</details>

***Any version of libappimage will work with bread but it is recommended to use v1.0.0 or greater, You can also Refer to this [list](https://repology.org/project/libappimage/versions) to check what version of libappimage your Distribution provides.***

---

## Removal

Just Remove the binary
```bash
sudo rm -v /usr/local/bin/bread
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
  <p>To install an application from a different Tag name you can specify the tag name too</p>
  <pre><code>bread install user/repo tagname</code></pre>
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

  <p>the <code>-n</code> or <code>--no-pre-release</code> flag can be used to disable updates for pre-releases.</p>
  <pre><code>bread update --no-pre-release</code></pre>
</details>

<details>
  <summary>Search for an application</summary>
  <br>
  <p>You can search for a application from the <a href="https://appimage.github.io">AppImage</a> API</p>
  <pre><code>bread search "Your search text"</code></pre>
</details>

<details>
  <summary>List all the installed application</summary>
  <br>
  <p>You can list all the installed applications using list command</p>
  <pre><code>bread list</code></pre>
  <p>If you also want to see the SHA1 Hashes of the applications listed, you can pass the <code>-s</code> or <code>--show-sha1</code> flag</p>
  <pre><code>bread list --show-sha1</code></pre>
  <p>If you want to see the GitHub release tag name <code>-t</code> or <code>--show-tag</code> flag</p>
  <pre><code>bread list --show-tag</code></pre>
</details>

---

### Bugs
- Icons not showing in menus until there's a system reboot
- Update Command Crashing

### Limits
- Bread uses GitHub API to get information about a repository and it's release, but without authentication GitHub API limits the request per hour.

---

## Tested On:
- Ubuntu 20.04 - by me
- Debian 11 - by me
- Manjaro Linux - by me
- Arch Linux - by [my brother](https://github.com/idno34)

---

## File/Folder Layout
Bread installs all the applications inside the `Applications` directory in your Linux Home Directory `~`, inside this directory there can be also a directory named `run-cache` which contains all the appimages you've run from the remote via the `bread run` command.

In the `Applications` there is also a file named `.registry.json` which contains information related to the installed applications!
In the `Applications` directory there is also a file named `.AppImageFeed.json` which is AppImage Catalog From [AppImage API](https://appimage.github.io/feed.json)

---
## Related:
- [Zap - :zap: Delightful AppImage package manager ](https://github.com/srevinsaju/zap)
- [A AppImage Manager Written in Shell](https://github.com/ivan-hc/AM-Application-Manager)
- [The Original Tool Which Bread is Based On](https://github.com/AppImageCrafters/appimage-cli-tool)

---

## Building From Source

Make Sure You Have Go version 1.18.x & [AppImage Builder](https://appimage-builder.readthedocs.io/en/latest/) Installed.

Get The Repository Via Git:

```bash
git clone https://github.com/pegvin/bread
```

Go Inside The Source Code Directory & Get All The Dependencies:

```bash
cd bread
go mod tidy
```

Make The Build Script Executable And Run It

```bash
chmod +x ./make
./make --prod
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
- [ ] Switch To Some Other Language Since Go Module System is Shit
- [ ] Improve UI
- [x] Make AppImages Runnable From Remote Without Installing (Done in v0.3.6)
- [x] Work On Reducing Binary Sizes (Reduced From 11.1MB to 3.1MB)
- [ ] Add 32 Bit Builds (Currently not possible since [DL](https://github.com/rainycape/dl) dependency is not available for 32 bit machines)
- [ ] Add Auto Updater Which Can Update The Bread Itself
- [x] Add `--version` To Get The Version (Done in v0.2.2)
- [x] Mirrors:
  - :heavy_multiplication_x: I Would Like To Introduce Concept Of Mirror Lists Which Contain The List Of AppImages With The Download URL, tho currently i am not working on it but in future i might.
  - [x] I am dropping this idea, tho i've added a search command which can search for appimages from a central server API

---

# Thanks
