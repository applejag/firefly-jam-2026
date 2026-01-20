<!--
SPDX-FileCopyrightText: 2025 Kalle Fagerberg

SPDX-License-Identifier: CC0-1.0
-->

# Epic Wizard Firefly Gladiators

[![latest version](https://img.shields.io/github/v/release/applejag/epic-wizard-firefly-gladiators)](https://github.com/applejag/epic-wizard-firefly-gladiators/releases)
[![download count](https://img.shields.io/github/downloads/applejag/epic-wizard-firefly-gladiators/total)](https://github.com/applejag/epic-wizard-firefly-gladiators/releases)
[![NO AI](https://raw.githubusercontent.com/nuxy/no-ai-badge/master/badge.svg)](https://github.com/nuxy/no-ai-badge)
[![REUSE status](https://api.reuse.software/badge/github.com/applejag/epic-wizard-firefly-gladiators)](https://api.reuse.software/info/github.com/applejag/epic-wizard-firefly-gladiators)

Game written in [Go](https://go.dev/) (with [TinyGo](https://tinygo.org/))
for the [Firefly Zero Jam #1](https://itch.io/jam/firefly-jam).

![main menu screenshot](./docs/startmenu.gif)

Jam theme: firefly

---

I originally tried creating this game in [MoonBit](https://www.moonbitlang.com/)
but I had so many compiler bugs that I ended up giving up and started porting
the game over to Go. See the [`moon-wasm-bug*` tags](https://github.com/applejag/epic-wizard-firefly-gladiators/tags)
on this repo.

## Play the game

Game can be found in the Firefly Zero catalog: <https://catalog.fireflyzero.com/applejag.ewfg>

On OS X (Mac) or Linux, run the following:

```bash
PLAY="$(curl https://fireflyzero.com/play.sh)"
bash -c "$PLAY" -- applejag.ewfg 
```

On Windows, install [firefly_cli](https://docs.fireflyzero.com/user/installation/)
and then run:

```bash
firefly_cli import applejag.ewfg
firefly_cli emulator --id applejag.ewfg 
```

## Play a specific version

Go to [releases](https://github.com/applejag/epic-wizard-firefly-gladiators/releases)
page and find the specific version, download the `applejag.ewfg.zip` file
from the "Assets" list on that version, and then load in the game with the
firefly CLI:

```bash
firefly_cli import ~/Downloads/applejag.ewfg.zip
firefly_cli emulator --id applejag.ewfg 

# or if you have "firefly_cli" installed as "ff"
ff import ~/Downloads/applejag.ewfg.zip
firefly_cli emulator --id applejag.ewfg 
```

## Credits

Made by: team applejuice

- Code: [@applejag](https://github.com/applejag)

- Art: [@JusJuice](https://github.com/JusJuice)

- Firefly names: Brooks Bedore, Dane, Spark, Kalob,
  and the cast of Titus Production's FROZEN
