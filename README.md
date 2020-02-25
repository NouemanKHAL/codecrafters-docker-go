This is a starting point for Go solutions to the
["Build Your Own Docker" Challenge](https://codecrafters.io/challenges/docker).

**Note**: If you're viewing this repo on GitHub, head over to
[codecrafters.io](https://codecrafters.io) to signup for early access.

# Usage

1. Ensure you have `go` installed locally
1. Run `./your_docker.sh` to run your Docker implementation, which is implemented in
   `app/main.go`.
1. Commit your changes and run `git push origin master` to submit your solution
   to CodeCrafters. Test output will be streamed to your terminal.

# Running your program locally

``` sh
docker build -t my_docker .
docker run --cap-add="SYS_ADMIN" my_docker \
run codecraftersio/docker-challenge /usr/bin/docker-explorer echo hey
```
   
# Passing the first stage

CodeCrafters runs tests when you do a `git push`. Make an empty commit and push
your solution to see the first stage fail.
   
``` sh
git commit --allow-empty -m "Running tests"
git push origin master
```

Go to `app/server.go` and uncomment the fork/exec implementation. Commit and
push your changes, and you'll now see the first stage pass.

Time to move on to the next stage!
