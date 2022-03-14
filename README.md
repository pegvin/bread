# Bread ![:bread:](./.github/bread.svg)

Install, update and remove AppImage from GitHub using your CLI.

## Features
- Install from the GitHub Releases
- Update with ease

## Installation

Download The Bread Binary into `/usr/local/bin`
```bash
sudo curl -L https://github.com/DEVLOPRR/bread/releases/download/v0.2.2/bread-0.2.2-x86_64.AppImage -O /usr/local/bin/bread
```

Give Executable Permissions To Downloaded Binary
```bash
sudo chmod +x /usr/local/bin/bread
```

One Liner:
```bash
sudo curl -L https://github.com/DEVLOPRR/bread/releases/download/v0.2.2/bread-0.2.2-x86_64.AppImage -o /usr/local/bin/bread && sudo chmod +x /usr/local/bin/bread
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

### Install a app from GitHub release
```bash
bread install user/repo
```

### Update a app
```bash
bread update user/repo
```

### List all of the installed apps
```bash
bread list
```

### Remove a installed app
```bash
bread remove user/repo
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
  list       List installed applications.
  remove     Remove an application.
  update     Update an application.

Run "bread <command> --help" for more information on a command.
```

---

## Building From Source

Make Sure You Have Go version 1.17.x & [AppImage Builder](https://appimage-builder.readthedocs.io/en/latest/) Installed.

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
- [ ] Work On Reducing Binary Sizes
- [x] Add `--version` To Get The Version (Done in v0.2.2)
- [ ] Mirrors:
  - I Would Like To Introduce Concept Of Mirror Lists Which Contain The List Of AppImages With The Download URL, tho currently i am not working on it but in future i might.

---

# Thanks
