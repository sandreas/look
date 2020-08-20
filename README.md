# look
Look is a tool to seek through and watch text based files (e.g. Logfiles) or input from STDIN. It makes heavy use of regular expressions to filter and replace content.

## Examples

```
# filter all lines containing error and warn and redirect it to only_err_warn.log
look text --where=".*(warn|error).*" my-logfile.log > only_err_warn.log

# using the where-not filter to skip lines containing warn or error but also containing specific other words
look text --where=".*(warn|error).*" --where-not=".*error404.*" my-logfile.log

# filter only the last (!) 10 lines containing error and warn and redirect it to only_err_warn.log
look text --lines=10 --where=".*(warn|error).*" my-logfile.log > only_err_warn.log


# watch a changing log file and only report lines matching error and warn
look text --watch --where=".*(warn|error).*" my-logfile.log

# extract all lines showing a responsetime > 500ms of a log into a csv file and be case insensitive with (?i)
look text --pattern="(?i).*(appserver|webserver|databaseserver)[\s]+responsetime[\s]+([5-9][0-9]{2}|[0-9]{4,})[\s]+ms.*" --replacement="$1,$2" my-logfile.log > long_response_times.csv
```


## command line reference

Running `look help text` will show the following reference:

```
NAME:
   look text - look at file with text lines

USAGE:
   look text [command options] [arguments...]

OPTIONS:
   --quiet                        do not show any output (default: false)
   --force                        force the requested action - even if it might be not a good idea (default: false)
   --debug                        debug mode with logging to Stdout and into $HOME/.graft/application.log (default: false)
   --where value, -w value
   --where-not value, -W value
   --pattern value, -p value
   --replacement value, -r value
   --watch                        (default: false)
   --lines value, -l value        (default: 0)
```
