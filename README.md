# Sun2000 modbus TCP interface

This is a simple interface pull data and stats over ModBus TCP, from a Sun2000-xxx inverter.

It implements a web-server, which provides a lot of text details and `/metrics` for [Prometheus](https://prometheus.io/).
From there, you can do some nice graphs for example with [Grafana](https://grafana.com/). Or modify this for your own
environment.

## License and Copying

**Author:** Dragos Vingarzan vingarzan -at- gmail dot com

**Copyright:** 2024 Dragos Vingarzan vingarzan -at- gmail -dot- com

**License:** [AGPL-3.0](./LICENSE)

This file is part of sun2000-modbus.

sun2000-modbus is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General
Public License Version 3 (AGPL-3.0) as published by the Free Software Foundation.

 sun2000-modbus is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied
 warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more
 details.

You should have received a copy of the AGPL-3.0 along with sun2000-modbus. If not, see [https://www.gnu.org/licenses/](https://www.gnu.org/licenses/).

Sun2000 is (probably) a (registered) trademark of sun2000. Any such names are used in here for pure informative purposes
and this package is not endorsed in any way by sun2000.

--------

For the modbus communication this uses [GoBorrow's implementation - github.com/goburrow/modbus](https://github.com/goburrow/modbus)

--------

## Warning

***!!! Use this code at your own risk. I take no warranties if you break your inverter, set your batteris on fire or kill your cat. !!!***

## References

The [Sun2000 Modbus Register](https://www.debacher.de/wiki/Sun2000_Modbus_Register) is what got me to the first messages
exchanged successfully with the inverter. Yet, some addresses seem to be off when tried on my `SUN2000-5KTL-M1`

Then the [SmartLogger ModBus Interface Definitions](https://support.sun2000.com/enterprise/en/doc/EDOC1100050690) seem
to be quite exhaustive, but when it comes to the Sun2000 product range, it says that the reference is in
`SUN2000VXXXRXXXCXX MODBUS Protocol`. I couldn't find that, but the closest for me was the
[Modbus Interface Definitions (V3.0)](https://community.symcon.de/uploads/short-url/pqZXWOienoBK2AsEzGD2oH1bcPR.pdf).

If anyone has newer/better references, I'd appreciate it.


## Usage

Install go and then run `go run .`, or compile and install it, etc. Eventually, I'll also make a Dockerfile.

Your inverter needs to have modbus enabled. Since that **would allow you to also write not just read data**, think
carefully before you open it up for a potential attack! Make sure that you secure this interface from any malicious
access!!!

To enable modbus, you'd probably have to login over WiFi as installer, with the SUN2000 app. The modbus TCP though
would be open on the dongle, which is typically used to connect your inverter to the Internet. ***I would not recommend
opening this even on your home LAN!!!*** Instead, maybe put the whole thing behind a Rasperry Pi, in a sort of guest
network environment, isolating it from devices which might be easily compromised.

Find the IP of your inverter and export it, then run this program:

    export MODBUS_IP="192.168.0.250"
    go run .


### Configuration

Set the following environment variables to your own desire:

| Variable          | Default   |Description |
|-------------------|-----------|------------|
| `HTTP_IP`         | 127.0.0.1 | IP to listen on and serve metrics to Prometheus |
| `HTTP_PORT`       | 8080      | Port to listen on and serve metrics to Prometheus |
| `MODBUS_IP`       | N/A       | IP of the Sun2000 inverter to scrape data from |
| `MODBUS_PORT`     | 502       | Port of ModBus on the inverter |
| `MODBUS_TIMEOUT`  | 5         | If the inverter does not answer, give up after this many seconds |
| `MODBUS_SLEEP`    | 5         | Interval in seconds to sleep after doing a full round of reads |

Data is read in address ranges, trying to minimize the number of commands, since each seems to be kind of slow. For each
suck block of data an expiration is set. Some data which very rarely changes is polled at longer intervals (e.g. model,
SN, etc at 1 hour), while for faster metrics we want the values to refresh much faster (e.g. 5 seconds).

After reading all the ranges which were expired, the poller sleeps for `MODBUS_SLEEP` seconds. Hence increasing it will
be nicer on the inverter, but your data will be more stale.


## Future work

Add testing...

At some point this could be extracted as a library and cleaned-up.

Providing containers or other packaging would be nice too.