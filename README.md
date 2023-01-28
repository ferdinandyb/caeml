# CAEML

Reads an email file from either STDIN or from a file passed as the first argument, digests it and outputs it to STDOUT in a more human readable format. This means only From, To, Cc, Bcc, Date and Subject headers are kept and these are decoded and of all the parts only text/plain is returned.

# Why?

A modern email will have about 60 lines of headers _before_ the From header and the "interesting" headers will be encoded in a way, that is pretty hard for humans to read. And that's if the email was two lines. If there were also attachments, then good luck finding anything useful quickly.

Of course this is why you use a MUA instead of reading raw emails, but `caeml` was written for two use cases. One, poking around in your maildir folder, when trying to figure out what's going wrong with your syncing and previewing `message/rfc822` attachments in [aerc](aerc-mail.org).

# Usage

- read from STDIN: `cat test.eml | caeml`
- read from file:  `caeml test.eml`
- change order of headers: `caeml -H "Subject,To,From,Message-Id" test.eml`
- print only headers: `caeml -O test.eml`
- also print contents of any `message/rfc822` parts: `caeml -D test.eml` (respects `-O` and `-H` for these as well)

Default header order is: `caeml -H "From,To,Cc,Bcc,Date,Subject"`
If a header is empty it will not be displayed.

# Integration

## aerc
```
message/rfc822=caeml | colorize
```

# Contribution

Patches sent in email are also welcome :)
