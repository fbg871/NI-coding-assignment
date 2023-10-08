# BoVel fullstack developer take home assignment

Hi there! 👋

Thank you for applying for the position of fullstack developer in No Isolation and the BoVel Komp project. We're very much looking forward to getting to know you a bit better!

In this project you will find a working fullstack solution to a fictional (but not entirely irrelevant) problem and we are going to ask you to make some additions and alterations to it. The point of this assignment is for us to learn a bit more about you, see what happens when you get thrown into an existing code base (which will be the case in this role), let you work on an actual fullstack problem, and lastly give us something specific to talk about in the technical interview.

Below you will find a description of the problem and solution, it's technical architecture, instructions on how to run it, the task itself and so on.

You're free to approach this in any way you see fit. With that being said, please do keep in mind that we will be going through your solution in the technical interview. It's a good idea to keep track of the changes you've made (initializing a new git repo and making commits to track your work could be a good idea). We **do not** expect you to spend a ton of time on this. Focus on solving the task as specified. Function over form is the point in this (although do keep in mind that someone is going to be reading your code, so please do make it reasonably readable). If you get stuck because of syntax or the like, please provide some thoughts on how you would have solved the task instead. Please complete the assignment and send the entire solution back to us by **Tuesday October 3rd 12:00PM CET** (as a zip, tar, github link, what have you - take your pick).

We wish you the best of luck, and we're looking forward to seeing you again for the technical interview!
Questions? Don't hesitate to ask 😄

## Introducing: The BoVel Komp registry
In the BoVel project, No Isolation has distributed Komps to various locations in Oslo kommune. This is all well and good, but the project has a need to keep track of these Komps, where are they, are they allocated (at a specific location), technical information about the physical unit and more.

Someone has already thrown together an internal tool that solves this problem. It does most of the things we need it to do (most things being the operative part). The solution needs a few more additions, and that's where you come in!

The solution itself is a simple enough web page that shows a list of Komps (identified by their unique serial number) and various metadata. Project members can select a Komp in the list display details about it, mark it as allocated, write a comment etc.

Here is a list of what the various properties in the solution mean:
* Serial number: The Komps unique human readable id, this is how everyone refers to a specific Komp
* Allocation state: Whether the Komp is `Allocated` or `Available`. Only `Available` Komps can be allocated and only `Allocated` Komps can be reset
* Software version: The version of our software running on the physical device
* Product code: The product code identifies the hardware revision of Komp. There's two possible values, `ev2` and `ev2b`. The primary difference is that `ev2b` has a bigger screen and includes an integrated 4G sim-card (more about that under `Attributes`)
* MAC address: The unique hardware address of the internal computer of the Komp. Useful for technical diagnostics
* Attributes: A list of key-value pairs used to track certain values of the integrated 4G sim-card (if one exists in that hardware model)
* Comments: Free text field defined by the users of the registry

### The solution
The existing solution is a straight forward system that's compromised of a backend written in Go that exposes a simple HTTP API. This HTTP API is consumed by a React/Next.js frontend. The backend uses a MySQL database for data persistence. Everything runs in a docker compose stack. Ingress into the stack is handled through a nginx reverse proxy that routes requests to either the Go HTTP API service or the Next.JS frontend service based on the URL in the HTTP request (all requests with the /api/ prefix is routed to the HTTP API, the rest go to the frontend).

HTTP request go to the nginx reverse proxy (mapped to port `8080` on the host machine), the requests then gets forwarded to the right internal service based on the URL of the request. Requests going to the Go HTTP API will also result in the Go backend reaching out to the mysql database service to query/update it's records.

This is already configured in the `docker-compose.yml` file and you should be able to get a fully functional solution up and running by executing `docker compose up -d` from the root directory (where the `docker-compose.yml` file is). Once that's running you should be able to open `http://localhost:8080` in your browser to use the solution.

*This means that you will have to have `docker` [installed](https://docs.docker.com/engine/install/) on your machine. We have had limited opportunity to test this stack on other operating systems besides macOS. In theory it should work on Linux/Windows as well, but your experience might prove otherwise if you're on one of those systems. Please reach out to us if you face issues running this on your operating system of choice*

You _can_ run everything directly on your host machine instead if you prefer. This means that you will have to install [go](https://go.dev/), [node](https://nodejs.org/en/download) and [mysql](https://dev.mysql.com/doc/mysql-installation-excerpt/8.0/en/). You will find all configuration/setup you'll need to do in the `docker-compose.yml` file and the various `Dockerfile` files.

We do recommend using `docker compose` for a smoother experience, but do what works best for you!

```bash
# Run the entire solution
$ docker compose up -d
```

When you're working on the assignment you will have to update the code in both the `frontend` and `backend` services. This will require you to rebuild the corresponding docker containers for the changes to take effect. This can be a bit finicky, so we've included some useful docker compose commands at the [end of the document](#useful-docker-compose-commands).

### The assignment

The assignment consists of three parts. Implement one feature, fix one bug and reflect on the state of the architecture/code base. Tackle them in whatever order you like, but please make sure to give all three attention! Read? Let's go!

#### Task 1

The BoVel Komp registry is doing mostly what it needs to do (it ain't pretty but it works!). There is one thing that the project still needs though, and that is the possibility "resetting" a Komp. Today we can mark a Komp as allocated, but there's no way to mark it as available again (part of resetting it). Please implement this possibility according to the requirements below:

* Only Komps with `state` equal to `allocated` can be reset
* Pressing the "Reset" button in the existing UI must do the following:
  * Set the Komps `state` to `available` in the database
  * Set the Komps `comment` to `null` in the database
  * If the Komps `software_version` is **not** `"v3.3.3"`, the `comment` should be set to `"Software version is %s. Software upgrade required"` (where `%s` should be replaced with the value of the Komps `software_version`) instead of `null`
  * If the Komp has `product_code` equal to `ev2b` we must also make sure that the `Attribute` with name equal to `"simcard_state"` gets its value set to `"inactive"`. If an attribute with the name `"simcard_state"` does not exist we must create one with the value `"inactive"`
* The UI must be automatically updated after the reset has been performed (both the Komp list and the details on the selected Komp needs to reflect the changes above automatically)

#### - Task 2

We have received a bug report from one of the project members. It goes like this:

> "This solution is so helpful. Thank you so much for building it! It's making our work so much easier. I wanted to let you know of this one little thing I've noticed though. When I allocate a Komp, it seems to delete the comments field? I've also seen the Allocation state just disappear sometimes? I'm not really sure how. It's really weird... Could you please have a look?"

Sounds like there's a bug in how we update a Komps state? See if you can't find and squash it!

#### - Task 3

As part of the ongoing BoVel project we are planning to do a lot of new development on this registry system. We have been requested by project management to advice on the technical state of the system, and provide some thoughts on what technical foundation work should be performed to assure continued and efficient feature development.

Since you've now worked on the system for a bit you're well equipped to advice. What would be your top 3 suggested points of improvements to the architecture and/or code base and why?

### Useful docker compose commands

```bash
# Run the entire solution
$ docker compose up -d
[+] Running 5/5
 ⠿ Network fullstack-take-home-assignment_default            Created                               0.0s
 ⠿ Container fullstack-take-home-assignment-mysql-1          Started                               0.4s
 ⠿ Container fullstack-take-home-assignment-backend-1        Started                              28.7s
 ⠿ Container fullstack-take-home-assignment-frontend-1       Started                              30.5s
 ⠿ Container fullstack-take-home-assignment-reverse-proxy-1  Started                              31.2s
```

```bash
# Rebuild and restart the `backend`
$ docker compose up -d --build --force-recreate backend
[+] Building 0.1s (9/9) FINISHED
 => [internal] load build definition from Dockerfile                                               0.0s
 => => transferring dockerfile: 32B                                                                0.0s
 => [internal] load .dockerignore                                                                  0.0s
 => => transferring context: 2B                                                                    0.0s
 => [internal] load metadata for docker.io/library/golang:latest                                   0.0s
 => [1/4] FROM docker.io/library/golang:latest                                                     0.0s
 => [internal] load build context                                                                  0.0s
 => => transferring context: 163B                                                                  0.0s
 => CACHED [2/4] RUN mkdir komp-registry                                                           0.0s
 => CACHED [3/4] COPY . komp-registry/                                                             0.0s
 => CACHED [4/4] RUN cd komp-registry && go build  -o komp-registry-backend *.go                   0.0s
 => exporting to image                                                                             0.0s
 => => exporting layers                                                                            0.0s
 => => writing image sha256:1ac83eef22b9ae5bd8e7a9a89b806cd09ba2105c035625a1621062233d91a9fe       0.0s
 => => naming to docker.io/library/komp-registry-backend:latest                                    0.0s

Use 'docker scan' to run Snyk tests against images to find vulnerabilities and learn how to fix them
[+] Running 2/2
 ⠿ Container fullstack-take-home-assignment-mysql-1    Running                                     0.0s
 ⠿ Container fullstack-take-home-assignment-backend-1  Started                                     1.0s
```

```bash
# Rebuild and restart the `frontend`
$ docker compose up -d --build --force-recreate frontend
[+] Building 28.8s (18/18) FINISHED
 => [komp-registry-backend:latest internal] load build definition from Dockerfile                  0.0s
 => => transferring dockerfile: 32B                                                                0.0s
 => [komp-registry-frontend:latest internal] load build definition from Dockerfile                 0.0s
  < ...Abbreviated... >
 => => transferring context: 83.24MB                                                               2.7s
 => CACHED [komp-registry-frontend:latest 2/5] WORKDIR komp-registry                               0.0s
 => [komp-registry-frontend:latest 3/5] COPY . .                                                   0.2s
 => [komp-registry-frontend:latest 4/5] RUN npm install                                            9.2s
 => [komp-registry-frontend:latest 5/5] RUN npm run build                                         13.9s

Use 'docker scan' to run Snyk tests against images to find vulnerabilities and learn how to fix them
[+] Running 3/3
 ⠿ Container fullstack-take-home-assignment-mysql-1     Running                                    0.0s
 ⠿ Container fullstack-take-home-assignment-backend-1   Running                                    0.0s
 ⠿ Container fullstack-take-home-assignment-frontend-1  Started                                    2.2s
```

```bash
# See the status of all containers in the stack
$ docker compose ps
NAME                                             COMMAND                  SERVICE             STATUS              PORTS
fullstack-take-home-assignment-backend-1         "/go/komp-registry/k…"   backend             running (healthy)
fullstack-take-home-assignment-frontend-1        "npm run start"          frontend            running
fullstack-take-home-assignment-mysql-1           "docker-entrypoint.s…"   mysql               running (healthy)   33060/tcp
fullstack-take-home-assignment-reverse-proxy-1   "/docker-entrypoint.…"   reverse-proxy       running             0.0.0.0:8080->8080/tcp
```

```bash
# Check the logs of the `frontend` service
$ docker compose logs frontend
fullstack-take-home-assignment-frontend-1  |
fullstack-take-home-assignment-frontend-1  | > komp-registry@0.1.0 start
fullstack-take-home-assignment-frontend-1  | > next start
fullstack-take-home-assignment-frontend-1  |
fullstack-take-home-assignment-frontend-1  |   ▲ Next.js 13.5.3
fullstack-take-home-assignment-frontend-1  |   - Local:        http://localhost:3000
fullstack-take-home-assignment-frontend-1  |
fullstack-take-home-assignment-frontend-1  |  ✓ Ready in 155ms
```

```bash
# Tear everything down
$ docker compose down
[+] Running 5/4
 ⠿ Container fullstack-take-home-assignment-reverse-proxy-1  Removed                              0.2s
 ⠿ Container fullstack-take-home-assignment-frontend-1       Removed                              0.7s
 ⠿ Container fullstack-take-home-assignment-backend-1        Removed                              0.1s
 ⠿ Container fullstack-take-home-assignment-mysql-1          Removed                              1.5s
 ⠿ Network fullstack-take-home-assignment_default            Removed                              0.1s
```