cloudwatcher
============

Collects memory, swap, and disk space utilization on an Amazon EC2 instance and sends this data as custom metrics to Amazon CloudWatch every 5 minutes.

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

* [ ] Fix disk space unit conversion
* [ ] Get utilization statistics for the last 12 hours
* [ ] Command-line interface

Credits
=======

Inspired by [Amazon CloudWatch Monitoring Scripts for Linux](http://docs.aws.amazon.com/AmazonCloudWatch/latest/DeveloperGuide/mon-scripts-perl.html).

License
=======

MIT.
