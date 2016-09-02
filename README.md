# Mark

Mark is a distributed social bookmarking network targeting Sandstorm (sandstorm.io)[https://sandstorm.io/].

I've been a fan of several unprofitable web services, and I get bummed out when they inevitably disappear.

This project is an experiment in seeing if we can build infrastructure for more resilient web utilities.

At this point, most core features remain un-implemented.
- You can start a node
- You can post bookmarks, or remove ones you've posted
- Other people can see those bookmarks
- You can see their bookmarks
- You can set a handle

Under the hood, there's a append-only log feed format that steals liberally from Secure-Scuttlebutt and JOSE. There's a simple gossip protocol, and there's an embedded/distributed database written in Go.

## Technical details

Run mark locally:

    go get github.com/awans/mark/cmd/main
    cd $GOROOT/src/github.com/awans/mark/
    npm install && npm run build
    mark init -d /var/opt/mark
    mark serve -d /var/opt/mark -p 8080

Code is organized in reasonably intuitive ways:

    .sandstorm/  # sandstorm package info
    app/         # core "application" code
    client/      # react/redux client
    cmd/         # the mark command
    design/      # sketch file with a purple square
    entities/    # an OK embedded DB, built on cznic/kv and feed
    feed/        # a distributed append-only log format based on SSBC
    sandstorm/   # server components for dealing with sandstorm
    server/      # very simple sync/api server
    
Very little is documented and much is probably still broken..
  
