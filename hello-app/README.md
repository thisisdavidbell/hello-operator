# The go app

This directory guides you through the process of creating a simple containerised Hello World httpserver go app.

In later chapters, this simple application will be the one managed by the operator you create.

## 0. Prereqs

- an OCP 4.4 or above cluster.
- docker or podman
- Set the following env vars:
    - DOCKERHOSTNAME - the fully qualified hostname of your docker repo
    - DOCKERNAMESPACE - the namespace within your docker repo where all images will be pushed.
    - OCPHOSTNAME - the fully qualified hostname of your ocp system. (e.g. everything after `api.` or `apps.` )

## 1. Create Go app
Create a simple go application, running a httpserver, returning a string.
You can use the simple example here: [hello-ocp.go](hello-ocp.go)

## 2. Test it
In one terminal, run:
`go run hello.go`

In a second terminal run:
`curl localhost:8080/hello`

Alternatively, access http://localhost:8080/hello in a browser.

You can now stop the app in the first terminal.

## 3. Build it
Run:
`go build hello.go`

## 4. Run it
Run it:
`./hello`

And test it:
`curl localhost:8080/hello`

Stop the app

## 5. Create a dockerfile
Kubernetes and Red Hat Openshift Container Platform utilise docker containers. 

Note you can use `podman` or `docker` for the following commands.

a. Write a Dockerfile to build and run your go application.
  - An example dockerfile is provided: [Dockerfile](dockerfile)

b. Test the dockerfile:
- Build the image
`docker build -t hello:v1.0 .`

- Confirm the image was created
`docker images | grep hello`

- Run the image locally:
`docker run -p 8080:8080 -d hello:v1.0`

- Confirm the container is running:
`docker ps`

- Test:
`curl localhost:8080/hello`

- Stop the container:
`docker stop CONTAINER-ID`
 where `CONTAINER-ID` is the value shown when running `docker ps`

## 6. Create a Makefile
It is common to use Makefile's or similar technologies to group together the common commands used. 
A simple helper [Makefile](Makefile) is provided covering a number of commands mentioned above, and introduced later.

## 7. Push your application to an image registry

Ensure you ahev set env vars:
- `IRHOSTNAME` - Image registry Hostname
- `IRNAMESPACE` - Namespace in image registry
- `IRUSER` - Image registry user name
- `IRPASSWORD` - Image resgiry password

Run:
- `docker login -u $IRUSER -p $IRPASSWORD $IRHOSTNAME`
- `docker tag hello:v1.0 $IRHOSTNAME/$NAMESPACE/hello:v1.0`
- `docker push $IRHOSTNAME/$NAMESPACE/hello:v1.0`

## 8. Deploy your application in Red Hat OpenShift

- log into OCP console
- Switch to Developer Perspective
- Run you app using one of:
  - From git
  - From Dockerfile
  - From Container images

For Container images, (and OCP 4.8), the process was:
- Create new project/namespace
- Click: From Container images
- Click link to create image pull sceret if using secure registry
- entire image 
- Leave defaults of `hello-app`, `hello`, deployment.
- Select to create route
- Expand advanced Routing options
- Enter hostname as: hello-image.app.<rest of your ocp console url after app.>
- Leave route unsecure
- Click on the app icon in the displayed Topology view.
- You should see it has created:
  - pod
  - service
  - root
- The pod should be Running.
- Click the route Location.
- You should see your go app output in the new tab.

Congratulations, you have just run your app as a deployment in Red Hat OpenShift.

**TODO**: write up the deployment, service and route yaml files, to show exactly what is defined.