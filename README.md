# Bread ![:bread:](./.github/bread.svg)

Install, update, remove & run AppImage from GitHub using your CLI. (Fork of [AppImage ClI Tool](https://github.com/AppImageCrafters/appimage-cli-tool))

## Features
- Install from the GitHub Releases
- Run Applications From Remote Without Installing Them
- Update with ease

## Installation

With Curl:
```bash
sudo curl -L https://github.com/DEVLOPRR/bread/releases/download/v0.3.5/bread-0.3.5-x86_64.AppImage -o /usr/local/bin/bread && sudo chmod +x /usr/local/bin/bread
```

With Wget:
```bash
sudo wget -O /usr/local/bin/bread https://github.com/DEVLOPRR/bread/releases/download/v0.3.5/bread-0.3.5-x86_64.AppImage && sudo chmod +x /usr/local/bin/bread
```

---

## Uninstallation

Just Remove the binary
```bash
rm -v /usr/local/bin/bread
```

**NOTE** this won't delete the app you've installed.

---

## Usage

**NOTE** if the user and repo in `user/repo` are same then you can just specify the single name itself, which means `libresprite/libresprite` is equal to `libresprite`

---

### Installing

Installing a App from GitHub Release
```bash
bread install user/repo
```

---

### Running

You can run a Application from remote without installing
```bash
bread run user/repo
```

You can also pass arguments to the application like this:
```bash
bread run user/repo -- --myArg1 --myArg2
```

Since all these applications are cached so that it isn't downloaded again, you can clean that cache too.
```bash
bread clean
```

---

### Updating

Just Checking For Update
```bash
bread update user/repo --check
```

Just Checking For All Updates
```bash
bread update user/repo --check --all
```

Updating A Single App
```bash
bread update user/repo
```

Updating All Of The Apps
```bash
bread update --all
```

---

### Removing

Completely Removing a installed app
```bash
bread remove user/repo
```

Only De-Integrating The App But Not Removing It 
```bash
bread remove user/repo --keep-file
```

---

### List all of the installed apps
```bash
bread list
```

---

## Full usage

```bash
Usage: bread <command>

Install, update and remove AppImage from GitHub using your CLI.

Flags:
  -h, --help       Show context-sensitive help.
      --version    Print version information and quit

Commands:
  install    Install an application.
  run        Run an application from Remote.
  list       List installed applications.
  remove     Remove an application.
  update     Update an application.
  clean      Clean all the cache.

Run "bread <command> --help" for more information on a command.
```

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

## Todo
- [ ] Improve UI
- [x] Make AppImages Runnable From Remote Without Installing (Done in v0.3.5)
- [ ] Work On Reducing Binary Sizes (Reduced A bit)
- [ ] Add 32 Bit Builds
- [ ] Add Auto Updater Which Can Update The Bread Itself
- [x] Add `--version` To Get The Version (Done in v0.2.2)
- [ ] Mirrors:
  - I Would Like To Introduce Concept Of Mirror Lists Which Contain The List Of AppImages With The Download URL, tho currently i am not working on it but in future i might.

---

# Thanks
