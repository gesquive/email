# email

Send an email from the command line

## Installing

### Compile
This project requires go1.6+ to compile. Just run `go get -u github.com/gesquive/email` and the executable should be built for you automatically in your `$GOPATH`.

Optionally you can run `make install` to build and copy the executable to `/usr/local/bin/` with correct permissions.

### Download
Alternately, you can download the latest release for your platform from [github](https://github.com/gesquive/email/releases/latest).

Once you have an executable, make sure to copy it somewhere on your path like `/usr/local/bin` or `C:/Program Files/`.
If on a \*nix/mac system, make sure to run `chmod +x /path/to/email`.

## Configuration

### Precedence Order
The application looks for variables in the following order:
 - command line flag
 - environment variable
 - config file variable
 - default

So any variable specified on the command line would override values set in the environment or config file.

### Config File
The application looks for a configuration file at the following locations in order:
 - `./config.yml`
 - `~/.config/email/config.yml`
 - `/etc/email/config.yml`

Copy `install/config.example.yml` to one of these locations and populate the values with your own. If the config contains your SMTP credentials, make sure to set permissions on the config file appropriately so others cannot read it. A good suggestion is `chmod 600 /path/to/config.yml`.

### Environment Variables
Optionally, instead of using a config file you can specify config entries as environment variables. Use the prefix "EMAIL_" in front of the uppercased variable name. For example, the config variable `smtp-server` would be the environment variable `EMAIL_SMTP_SERVER`.

## Usage

```console
Send an email from the command line.

If a flag is tagged with 'multi', multiple versions of the flag are accepted

Usage:
  email [flags] <message>

Flags:
  -a, --attachment value       File to attach to email (multi)
  -b, --bcc value              Blind carbon copy addresses (multi)
  -c, --cc value               Carbon copy addresses (multi)
      --config string          config file (default is $HOME/.config/email.yml)
  -f, --from string            From address on email (default $USER@$HOST)
  -H, --html-message string    Alternate HTML content of email
  -m, --message string         Plain text content of email
  -r, --reply-to string        Reply to address
  -p, --smtp-password string   Authenticate the SMTP server with this password
  -o, --smtp-port value        The port to use for the SMTP server (default 25)
  -x, --smtp-server string     The SMTP server to send email through (default "localhost")
  -u, --smtp-username string   Authenticate the SMTP server with this user
  -e, --strict-parsing         Fail to send the email when any email address is malformed
  -s, --subject string         Subject of email
  -t, --to value               Destination addresses (multi)
  -v, --version                Show the version and exit
```

In addition the plain text message itself can be piped into the app instead of using the flag.

Optionally, a hidden debug flag is available in case you need additional output.
```console
Hidden Flags:
  -D, --debug                  Include debug statements in log output
```

## QuickStart

```console
# Simple email
$ email --from neo@hackers.org --to trinity@underground.org --subject "Question" --message "What is the matrix"

# Lets pipe the message in
$ cat blue-red_pill_speech.txt | email --from morpheus@underground.org --to neo@hackers.org\
--cc trinity@underground.org --subject "Warning"

# Setting multiple recipients
$ email --from agent.smith@matrix.net --to agent.jones@matrix.net\
 --to agent.brown@matrix.net --bcc the.architect@matrix.net\
 --subject "The insider"\
 --message "Never send a human to do a machine's job"

# Attachments are easy
$ email --from spoon.boy@oracle.org --to neo@underground.org\
 --subject "Do not try and bend the spoon"\
 --message "There is no spoon"\
 --attachment spoon.jpg --attachment no_spoon.jpg

```

## Documentation

This documentation can be found at github.com/gesquive/email

## License

This package is made available under an MIT-style license. See LICENSE.

## Contributing

PRs are always welcome!
