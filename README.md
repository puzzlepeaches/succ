
<h1 align="center">succ</h1>

<h3 align="center">Suck up some domains from MS</h3>

## Why? 

`succ` is a simple command line tool that queries Microsoft for a list of domains associated with an Office 365 tenant. Specifically, it queries an Autodiscover endpoint using a specially crafted XML payload. The response contains a list of domains that are associated with the tenant.

This tool is a simple continuation of tools already on the market such as:

* [bbot](https://blog.blacklanternsecurity.com/p/bbot)
* [letItGo](https://github.com/SecurityRiskAdvisors/letItGo)
* [AADInternals](https://github.com/Gerenios/AADInternals)
* etc.

The aim of this utility is to greatly simplify this enumeration process without having to install a full tool suite or run a Python script. 

Now, as a tester or bug bounty hunter, you can simply run `succ` and get a list of domains associated with the tenant without a ton of cruft surrounding it. This makes it easy to pipe your results to other tools for further enumeration.

<br>

## Installation

Installation is very simple. Once you have Go installed, simply run:

```bash
go install github.com/puzzlepeaches/succ@latest
```

<br>

## Usage

The help menu for `succ` is as follows:

```bash
succ up domains from MS

Usage:
  succ [domain] [flags]
  succ [command]

Available Commands:
  help        Help about any command
  version     Print the version number of generated code example

Flags:
  -h, --help            help for succ
  -o, --output string   Output file.

Use "succ [command] --help" for more information about a command.
```

The only additional option outside of the domain argument is output. This allows you to specify a file to write the results to. If you do not specify an output file, the results will be written to stdout.


<br>

## Example

_Insert corny bugbounty Tesla example below_

```bash
$ succ tesla.com

tesla.com
tesla.com
teslamotors.com
perbix.com
tesla.services
service.tesla.com
c.tesla.com
mta.tesla.com
m.tesla.com
siilion.com
teslaalerts.com
teslagrohmannautomation.de
t.tesla.com
solarcity.com
```

Note that the tool does not try to do anything fancy with the results. The only post-processing we do involves the removal of onmicrosoft.com domains. We may in future releases look to filter out subdomains but for now I find it useful. 

<br>

## Planned Features

* Add support for multiple domains
* Add support for reading domains from a file
* Add support for reading domains from stdin
* Add support for filtering out subdomains
* Add support for filtering out domains that do not resolve
* Add JSONL output format support


