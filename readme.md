# Monitor and Run

A simple Go program to monitor a file/directory and run a script when something changes (creation, modification, deletion).

## Example

### The script
```bash
#!/usr/bin/env bash
#myscript.sh

# Run make every time something changes.
# Probably not a great idea unless you are watching a single file, like your Pandoc-based resume for instance.
make
```

### Monitor something and act on changes
`monitor_and_run -script=myscript.sh -directory=resume.md`

Now each time you change your resume, myscript.sh will be called, presumably generating PDF, HTML, and other versions of your resume.

## More Info

Your script should be idempotent, especially if watching an entire directory.

I wrote this out of an interest in Go and the desire for a straightforward file watcher.
There are far more sophiscated (and maintained) solutions available.
See [this question](http://superuser.com/questions/181517/how-to-execute-a-command-whenever-a-file-changes) or try your luck with Google if you have more complex needs.
However, this tool solves a problem for me, and it may for you as well.
Feel free to use it.

If you find any serious problems, please create an [issue](https://github.com/ryan-robeson/monitor_and_run/issues).
However, if you have feature requests, please take a look at similar programs first, or [fork](https://github.com/ryan-robeson/monitor_and_run#fork-destination-box) this project and add them.
I do not see a need for yet another file watching utility, so I intend to keep this as simple as possible for my use case.
