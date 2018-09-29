# [stockwatch](//github.com/TheSp1der/stockwatch)

## Table of Contents

1. [Description](#description)
2. [Credit](#credit)
3. [Usage](#usage)
4. [Command Line Options](#command-line-options)
4. [License](#license)

## Description

[stockwatch](//github.com/TheSp1der/stockwatch) is a console stock informational
system written in go.

[[https://github.com/TheSp1der/stockwatch/blob/master/readme-images/console-1.png]]

[stockwatch](//github.com/TheSp1der/stockwatch) uses free data provided by
[IEX](//iextrading.com/developer/) to obtain market data. Please read and agree
to the the [IEX Terms of Use](//iextrading.com/api-exhibit-a/) prior to using this
program.

## Credit

* [stockwatch](//github.com/TheSp1der/stockwatch) was inspired by the
[goiex](//github.com/AndrewRPorter/goiex) library written by
[AndrewRPorter](//github.com/AndrewRPorter).
* Console color output is provided by the [color](//github.com/fatih/color) library
written by [Fatih](//github.com/fatih).

## Usage

[stockwatch](//github.com/TheSp1der/stockwatch) is designed to run in two modes.

### One-Shot Mode

When [stockwatch](//github.com/TheSp1der/stockwatch) is called without any e-mail
configuration it will run in a *one-shot* mode. This mode quickly obtains the
configured stock information and promptly terminates.

### Monitor Mode

When stockwatch is run by providing a to e-mail address, from e-mail address, and
mail server it run in a monitor mode. This mode will determine the open/closed time
for the NYSE and e-mail you at the close of every business day with the closing
prices.

## Command Line Options

| Command Line Option | Environment Variable | Default Value     | Purpose |
|---------------------|----------------------|-------------------|---------|
| -email              | EMAIL_ADDR           | null              | Destination e-mail address that will receive the end of day summary. |
| -from               | EMAIL_FROM           | noreply@localhost | Address the message will be sent from. |
| -host               | EMAIL_HOST           | null              | E-Mail server host.
| -invest             |                      | null              | Used for tracking current investments. Please see [invest](#invest-option) below. |
| -port               | EMAIL_PORT           | 25                | E-Mail server port. |
| -ticker             | TICKERS              | null              | Comma seperated list of stocks to report. |
| -verbose            | VERBOSE              | false             | Display current stock values every 5 seconds when run in monitor mode. |

### Invest Option

The invest option is a method of inputting your current investments to track an
overall gain/loss based on the current price. For example, say you have purchased
30 shares of amd (Advanced Micro Devices Inc.) for $28.44, and 10 more shares of
amd at $30.12, and 12 shares of googl (Alphabet Inc.) for $584.67.
[stockwatch](//github.com/TheSp1der/stockwatch) will take that into consideration
and display your overall gain/loss based on the most recent price.

To call stockwatch with that data you supply the -invest option with the ticker,
the quatnity, and the purchased price in that order, seperated by a comma. Like
so:

```bash
stockwatch -ticker amd,googl -invest amd,30,28.44 -invest amd,10,30.12 -invest googl,12,584.67
```

The output would look like this:
[[https://github.com/TheSp1der/stockwatch/blob/master/readme-images/console-2.png]]

## License

BSD 2-Clause License

Copyright (c) 2018, TheSp1der
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.