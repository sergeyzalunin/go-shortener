## GO-Shortener

This project is a playground to check golang features and other techniques.

The origin of the code were these videos: [video 1](https://youtu.be/rQnTtQZGpg8), [video 2](https://youtu.be/xUYDkiPdfWs), [video 3](https://youtu.be/QyBXz9SpPqE)

## Task file commands:
1. **dc-build-and-up** - Calls **make-container** task then calls **dc-up** task. This is the most faster way to start from scratch GO-Shortener and its infrastructure. 

1. **task make-container** - Calls **build-app** task then build a container file based on Dockerfile. 

1. **task build-app** - Builds app for linux/x64 architecture without CGO enabled. The GO-Shortener execution file is stored in the Bin folder.

1. **task dc-up** - Runs the docker-compose script 

1. **task dc-down** - Stops the docker-compose script 