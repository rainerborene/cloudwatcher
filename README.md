cloudwatcher
============

Collects memory, swap, and disk space utilization on an Amazon EC2 instance and sends this data as custom metrics to Amazon CloudWatch periodically.

Metrics
=======

- Memory Utilization (%)
- Memory Used (MB)
- Memory Available (MB)
- Swap Utilization (%)
- Swap Used (MB)
- Disk Space Utilization (%)
- Disk Space Used (GB)
- Disk Space Available (GB)

Todo
====

* [x] Fix disk space conversion
* [x] Fix InvalidClientTokenId error
* [x] Command-line interface
* [x] Custom time periodicity
* [ ] Get utilization statistics for the last X hours
* [ ] Validate config values using `Valid()` function

Credits
=======

Inspired by [Amazon CloudWatch Monitoring Scripts for Linux](http://docs.aws.amazon.com/AmazonCloudWatch/latest/DeveloperGuide/mon-scripts-perl.html).

License
=======

MIT.
