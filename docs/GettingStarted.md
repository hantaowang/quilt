# GETTING STARTED
This guide explains how to install Quilt, and also serves as a
brief, hands-on introduction to some Quilt basics.

## Install Node.js
Quilt blueprints are written in Node.js. Only Node version v7.10.0 has been
tested.

Installation instructions for various operating systems are available
[here](https://nodejs.org/en/download/).

## Install Go
Quilt supports Go version 1.5 or later.

Find Go using your package manager or on the [Golang website](https://golang.org/doc/install).

### Setup GOPATH
We recommend reading the overview to Go workplaces [here](https://golang.org/doc/code.html).

Before installing Quilt, you'll need to set up your GOPATH. Assuming the root of
your Go workspace will be `$HOME/gowork`, execute the following `export` commands in
your terminal to set up your `GOPATH`.
```bash
export GOPATH="$HOME/gowork"
export PATH="$PATH:$GOPATH/bin"
```
It would be a good idea to add these commands to your `.bashrc` so that they do
not have to be run again.

## Download and Install Quilt
Clone the repository into your Go workspace: `go get github.com/quilt/quilt`.

This command also automatically installs Quilt. If the installation was
successful, then the `quilt` command should execute successfully in your shell.

## Configure A Cloud Provider

Below we discuss how to setup Quilt for Amazon EC2. Google Compute Engine is
also supported. Since Quilt deploys systems consistently across providers, the
details of the rest of this document will apply no matter what provider you
choose.

For Amazon EC2, you'll first need to create an account with [Amazon Web
Services](https://aws.amazon.com/ec2/) and then find your
[access credentials](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-set-up.html#cli-signup).
That done, you simply need to populate the file `~/.aws/credentials`, with your
Amazon credentials:
```
[default]
aws_access_key_id = <YOUR_ID>
aws_secret_access_key = <YOUR_SECRET_KEY>
```

## Your First Quilt-managed Infrastructure
We suggest you read
[`quilt/nginx/main.js`](https://github.com/quilt/nginx/blob/master/main.js)
to understand the infrastructure defined by this Quilt.js blueprint.

### Acquire the Nginx Blueprint
In order to run the Nginx blueprint, we'll have to download it first. We'll simply
clone it:
```bash
git clone https://github.com/quilt/nginx
cd nginx
```

### Install Blueprint Dependencies
The Nginx blueprint depends on the `@quilt/quilt` module. More complicated
blueprints
may have other dependencies that would get pulled in as well. To install all
dependencies, run `npm install .`.

### Configure `quilt/nginx/main.js`
#### Set Up Your SSH Authentication
Quilt-managed Machines use public key authentication to control SSH access.
SSH authentication is configured with the `sshKeys` Machine attribute.
Currently,  the easiest way to set up your SSH access, is by using the
`githubKeys()` function. Given your GitHub username, the function grabs your
public keys from GitHub, so they can be used to configure SSH authentication.
If you can access GitHub repositories through SSH, then you can also SSH into a
`githubKey`-configured Machine.

If you would like to use `githubKey` authentication, open `main.js`, import the
`githubKeys` function from `@quilt/quilt`, and set the `sshKeys` appropriately.
```javascript
const {createDeployment, Machine, githubKeys} = require('@quilt/quilt');
...
var baseMachine = new Machine({
    ...
    sshKeys: githubKeys("CHANGE_ME"),
    ...
});
```

### Deploying `quilt/nginx/main.js`
In one shell, start the Quilt daemon with `quilt daemon`. In another shell,
execute `quilt run ./main.js`. Quilt will set up several
Ubuntu VMs on your cloud provider as Workers, and these Workers will host Nginx
Docker containers as specified in
[`quilt/nginx/app.js`](https://github.com/quilt/nginx/blob/master/app.js)
(you do not have to understand or edit this file).


### Accessing the Worker VM
It will take a while for the VMs to boot up, for Quilt to configure the network,
and for Docker containers to be initialized. When a machine is marked
`Connected` in the console output, the corresponding VM is fully booted and has
begun communicating with Quilt.

The public IP of the Worker VM can be deduced from the console output. The
following output shows the Worker VM's public IP to be 52.53.177.110:
```
INFO [Nov 11 13:23:10.266] db.Machine:
	Machine-2{Master, Amazon us-west-1 m4.large, sir-3sngfxdh, PublicIP=54.183.169.245, PrivateIP=172.31.2.178, Disk=32GB, Connected}
	Machine-4{Worker, Amazon us-west-1 m4.large, sir-19bid86g, PublicIP=52.53.177.110, PrivateIP=172.31.0.87, Disk=32GB, Connected}
...
```

Run `ssh quilt@<WORKER_PUBLIC_IP>` to access a privileged shell on the Worker VM.

### Inspecting Docker Containers on the Worker VM
You can run `docker ps` to list the containers running on your Worker VM.

```
quilt@ip-172-31-0-87:~$ docker ps
CONTAINER ID        IMAGE                        COMMAND                  CREATED             STATUS              PORTS               NAMES
a2ac27cfd313        quay.io/coreos/etcd:v3.0.2   "/usr/local/bin/etcd "   11 minutes ago      Up 11 minutes                           etcd
0f407bd0d5c4        quilt/ovs                    "run ovs-vswitchd"       11 minutes ago      Up 11 minutes                           ovs-vswitchd
7b65a447fe54        quilt/ovs                    "run ovsdb-server"       11 minutes ago      Up 11 minutes                           ovsdb-server
deb4f98db8eb        quilt/quilt:latest           "quilt minion"           11 minutes ago      Up 11 minutes                           minion
```

Any docker containers defined in a blueprint are placed on one of
your Worker VMs.  In addition to these user-defined containers, Quilt also
places several support containers on each VM. Among these support containers is
`minion`, which locally manages Docker and allows Quilt VMs to talk to each
other and your local computer.

### Loading the Nginx Webpage
By default, Quilt-managed containers are disconnected from the public internet
and isolated from one another. In order to make the Nginx container accessible
from the public internet,
[`quilt/nginx/app.js`](https://github.com/quilt/nginx/blob/master/app.js)
explicitly opens port 80 on the Nginx container to the outside world:

```javascript
publicInternet.connect(port, webTier);
```

From your browser via `http://<WORKER_PUBLIC_IP>`, or on the command-line via
`curl <WORKER_PUBLIC_IP>`, you can load the Nginx welcome page served by your
Quilt cluster.

### Cleaning up

If you'd like to destroy the infrastructure you just deployed, you can either
modify the blueprint to remove all of the Machines, or use the command,
`quilt stop`. Both options will cause Quilt to destroy all of the
Machines in the deployment.

## Next Steps: Writing your own Quilt Blueprint

This guide illustrated how to use Quilt to run an application that already had a
blueprint written.  Next, you can try writing your own Quilt blueprints for new
applications
that don't yet have blueprints written; to do that, check out the
[guide to writing Quilt blueprints](https://github.com/quilt/quilt/blob/master/docs/BlueprintWritersGuide.md).


