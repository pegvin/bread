# Bread ![:bread:](./.github/bread.svg)

Install, update and remove AppImage from GitHub using your CLI.

## Features
- Install from the GitHub Releases
- Update with ease

## Installation

Download The Bread Binary into `/usr/local/bin`
```bash
sudo curl -L https://github.com/DEVLOPRR/bread/releases/download/v0.1.0/bread-0.1.0-x86_64.AppImage -O /usr/local/bin/bread
```

Give Executable Permissions To Downloaded Binary
```bash
sudo chmod +x /usr/local/bin/bread
```

One Liner:
```bash
sudo curl -L https://github.com/DEVLOPRR/bread/releases/download/v0.1.0/bread-0.1.0-x86_64.AppImage -o /usr/local/bin/bread && sudo chmod +x /usr/local/bin/bread
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

### Install a app from GitHub release
```bash
bread install user/repo
```

### Update a app
```bash
bread update app-id
```

### List all of the installed apps
```bash
bread list
```

### Remove a installed app
```bash
bread remove app-id
```

**NOTE** here app-id is the Application id which you have installed. to get your application id run `imagehub list`

---

### Full usage

```bash
Usage: bread <command>

Flags:
  --help     Show context-sensitive help.

Commands:
  install <target>
    Install an application.

  list
    List installed applications.

  remove <id>
    Remove an application.

  update [<targets> ...]
    Update an application.

Run "bread <command> --help" for more information on a command.
```

---

## Thanks
