const crypto = require('crypto');
const request = require('sync-request');
const stringify = require('json-stable-stringify');
const _ = require('underscore');

const githubCache = {};
function githubKeys(user) {
    if (user in githubCache) {
        return githubCache[user];
    }

    const response = request('GET', `https://github.com/${user}.keys`);
    if (response.statusCode >= 300) {
        // Handle any errors.
        throw `HTTP request for ${user}'s github keys failed with error ` +
            `${response.statusCode}`;
    }

    const keys = response.getBody('utf8').trim().split('\n');
    githubCache[user] = keys;

    return keys;
}

// The default deployment object. createDeployment overwrites this.
global._quiltDeployment = new Deployment({});

// The label used by the QRI to denote connections with public internet.
var publicInternetLabel = "public";

// Global unique ID counter.
var uniqueIDCounter = 0;

// Overwrite the deployment object with a new one.
function createDeployment(deploymentOpts) {
    global._quiltDeployment = new Deployment(deploymentOpts);
    return global._quiltDeployment;
}

function Deployment(deploymentOpts) {
    deploymentOpts = deploymentOpts || {};

    this.maxPrice = deploymentOpts.maxPrice || 0;
    this.namespace = deploymentOpts.namespace || "default-namespace";
    this.adminACL = deploymentOpts.adminACL || [];

    this.machines = [];
    this.containers = {};
    this.services = [];
    this.connections = [];
    this.placements = [];
    this.invariants = [];
}

function omitSSHKey(key, value) {
    if (key == "sshKeys") {
        return undefined;
    }
    return value;
}

// Returns a globally unique integer ID.
function uniqueID() {
    return uniqueIDCounter++;
}

// key creates a string key for objects that have a _refID, namely Containers
// and Machines.
function key(obj) {
    var keyObj = obj.clone();
    keyObj._refID = "";
    return stringify(keyObj, { replacer: omitSSHKey });
}

// setQuiltIDs deterministically sets the id field of objects based on
// their attributes. The _refID field is required to differentiate between multiple
// references to the same object, and multiple instantiations with the exact
// same attributes.
function setQuiltIDs(objs) {
    // The refIDs for each identical instance.
    var refIDs = {};
    objs.forEach(function(obj) {
        var k = key(obj);
        if (!refIDs[k]) {
            refIDs[k] = [];
        }
        refIDs[k].push(obj._refID);
    });

    // If there are multiple references to the same object, there will be duplicate
    // refIDs.
    Object.keys(refIDs).forEach(function(k) {
        refIDs[k] = _.sortBy(_.uniq(refIDs[k]), _.identity);
    });

    objs.forEach(function(obj) {
        var k = key(obj);
        obj.id = hash(k + refIDs[k].indexOf(obj._refID));
    });
}

function hash(str) {
    const shaSum = crypto.createHash('sha1');
    shaSum.update(str);
    return shaSum.digest('hex');
}

// Convert the deployment to the QRI deployment format.
Deployment.prototype.toQuiltRepresentation = function() {
    this.vet();

    setQuiltIDs(this.machines);

    var containers = [];
    this.services.forEach(function(serv) {
        serv.containers.forEach(function(c) {
            containers.push(c);
        });
    });
    setQuiltIDs(containers);

    // Map from container ID to container.
    var containerMap = {};

    var services = [];
    var connections = [];
    var placements = [];

    // For each service, convert the associated connections and placement rules.
    // Also, aggregate all containers referenced by services.
    this.services.forEach(function(service) {
        connections = connections.concat(service.getQuiltConnections());
        placements = placements.concat(service.getQuiltPlacements());

        // Collect the containers IDs, and add them to the container map.
        var ids = [];
        service.containers.forEach(function(container) {
            ids.push(container.id);
            containerMap[container.id] = container;
        });

        services.push({
            name: service.name,
            ids: ids,
            annotations: service.annotations
        });
    });

    var containers = [];
    Object.keys(containerMap).forEach(function(cid) {
        containers.push(containerMap[cid]);
    });

    return {
        machines: this.machines,
        labels: services,
        containers: containers,
        connections: connections,
        placements: placements,
        invariants: this.invariants,

        namespace: this.namespace,
        adminACL: this.adminACL,
        maxPrice: this.maxPrice
    };
};

// Check if all referenced services in connections and placements are really deployed.
Deployment.prototype.vet = function() {
    var labelMap = {};
    this.services.forEach(function(service) {
        labelMap[service.name] = true;
    });

    var dockerfiles = {};
    var hostnames = {};
    this.services.forEach(function(service) {
        service.connections.forEach(function(conn) {
            var to = conn.to.name;
            if (!labelMap[to]) {
                throw service.name + " has a connection to undeployed service: " + to;
            }
        });

        var hasFloatingIp = false;
        service.placements.forEach(function(plcm) {
            if (plcm.floatingIp) {
                hasFloatingIp = true;
            }

            var otherLabel = plcm.otherLabel;
            if (otherLabel !== undefined && !labelMap[otherLabel]) {
                throw service.name + " has a placement in terms of an " +
                    "undeployed service: " + otherLabel;
            }
        });

        if (hasFloatingIp && service.incomingPublic.length
            && service.containers.length > 1) {
            throw service.name + " has a floating IP and multiple containers. " +
              "This is not yet supported."
        }

        service.containers.forEach(function(c) {
            var name = c.image.name;
            if (dockerfiles[name] != undefined && dockerfiles[name] != c.image.dockerfile) {
                throw name + " has differing Dockerfiles";
            }
            dockerfiles[name] = c.image.dockerfile;

            if (c.hostname !== undefined) {
                if (hostnames[c.hostname]) {
                    throw "hostname \"" + c.hostname + "\" used for multiple containers";
                }
                hostnames[c.hostname] = true;
            }
        })
    });
};

// deploy adds an object, or list of objects, to the deployment.
// Deployable objects must implement the deploy(deployment) interface.
Deployment.prototype.deploy = function(toDeployList) {
    if (toDeployList.constructor !== Array) {
        toDeployList = [toDeployList];
    }

    var that = this;
    toDeployList.forEach(function(toDeploy) {
        if (!toDeploy.deploy) {
            throw "only objects that implement \"deploy(deployment)\" can be deployed";
        }
        toDeploy.deploy(that);
    });
};

Deployment.prototype.assert = function(rule, desired) {
    this.invariants.push(new Assertion(rule, desired));
};

function Service(name, containers) {
    this.name = uniqueLabelName(name);
    this.containers = containers;
    this.annotations = [];
    this.placements = [];

    this.connections = [];
    this.outgoingPublic = [];
    this.incomingPublic = [];
}

// Get the Quilt hostname that represents the entire service.
Service.prototype.hostname = function() {
    return this.name + ".q";
};

// Get a list of Quilt hostnames that address the containers within the service.
Service.prototype.children = function() {
    var i;
    var res = [];
    for (i = 1; i < this.containers.length + 1; i++) {
        res.push(i + "." + this.name + ".q");
    }
    return res;
};

Service.prototype.annotate = function(annotation) {
    this.annotations.push(annotation);
};

Service.prototype.canReach = function(target) {
    if (target === publicInternet) {
        return reachable(this.name, publicInternetLabel);
    }
    return reachable(this.name, target.name);
};

Service.prototype.canReachACL = function(target) {
    return reachableACL(this.name, target.name);
};

Service.prototype.between = function(src, dst) {
    return between(src.name, this.name, dst.name);
};

Service.prototype.neighborOf = function(target) {
    return neighbor(this.name, target.name);
};


Service.prototype.deploy = function(deployment) {
    deployment.services.push(this);
};

Service.prototype.connect = function(range, to) {
    range = boxRange(range);
    if (to === publicInternet) {
        return this.connectToPublic(range);
    }
    this.connections.push(new Connection(range, to));
};

// publicInternet is an object that looks like another service that can be
// connected to or from. However, it is actually just syntactic sugar to hide
// the connectToPublic and connectFromPublic functions.
var publicInternet = {
    connect: function(range, to) {
        to.connectFromPublic(range);
    },
    canReach: function(to) {
        return reachable(publicInternetLabel, to.name);
    }
};

// Allow outbound traffic from the service to public internet.
Service.prototype.connectToPublic = function(range) {
    range = boxRange(range);
    if (range.min != range.max) {
        throw "public internet cannot connect on port ranges";
    }
    this.outgoingPublic.push(range);
};

// Allow inbound traffic from public internet to the service.
Service.prototype.connectFromPublic = function(range) {
    range = boxRange(range);
    if (range.min != range.max) {
        throw "public internet cannot connect on port ranges";
    }
    this.incomingPublic.push(range);
};

Service.prototype.place = function(rule) {
    this.placements.push(rule);
};

Service.prototype.getQuiltConnections = function() {
    var connections = [];
    var that = this;

    this.connections.forEach(function(conn) {
        connections.push({
            from: that.name,
            to: conn.to.name,
            minPort: conn.minPort,
            maxPort: conn.maxPort
        });
    });

    this.outgoingPublic.forEach(function(rng) {
        connections.push({
            from: that.name,
            to: publicInternetLabel,
            minPort: rng.min,
            maxPort: rng.max
        });
    });

    this.incomingPublic.forEach(function(rng) {
        connections.push({
            from: publicInternetLabel,
            to: that.name,
            minPort: rng.min,
            maxPort: rng.max
        });
    });

    return connections;
};

Service.prototype.getQuiltPlacements = function() {
    var placements = [];
    var that = this;
    this.placements.forEach(function(placement) {
        placements.push({
            targetLabel: that.name,
            exclusive: placement.exclusive,

            otherLabel: placement.otherLabel || "",
            provider: placement.provider || "",
            size: placement.size || "",
            region: placement.region || "",
            floatingIp: placement.floatingIp || ""
        });
    });
    return placements;
};

var labelNameCount = {};
function uniqueLabelName(name) {
    if (!(name in labelNameCount)) {
        labelNameCount[name] = 0;
    }
    var count = ++labelNameCount[name];
    if (count == 1) {
        return name;
    }
    return name + labelNameCount[name];
}

// Box raw integers into range.
function boxRange(x) {
    if (x === undefined) {
        return new Range(0, 0);
    }
    if (typeof x === "number") {
        x = new Range(x, x);
    }
    return x;
}

function Machine(optionalArgs) {
    this._refID = uniqueID();

    this.provider = optionalArgs.provider || "";
    this.role = optionalArgs.role || "";
    this.region = optionalArgs.region || "";
    this.size = optionalArgs.size || "";
    this.floatingIp = optionalArgs.floatingIp || "";
    this.diskSize = optionalArgs.diskSize || 0;
    this.sshKeys = optionalArgs.sshKeys || [];
    this.cpu = boxRange(optionalArgs.cpu);
    this.ram = boxRange(optionalArgs.ram);
    this.preemptible = optionalArgs.preemptible !== undefined ? optionalArgs.preemptible : false;
}

Machine.prototype.deploy = function(deployment) {
    deployment.machines.push(this);
};

// Create a new machine with the same attributes.
Machine.prototype.clone = function() {
    // _.clone only creates a shallow copy, so we must clone sshKeys ourselves.
    var keyClone = _.clone(this.sshKeys);
    var cloned = _.clone(this);
    cloned.sshKeys = keyClone;
    return new Machine(cloned);
};

Machine.prototype.withRole = function(role) {
    var copy = this.clone();
    copy.role = role;
    return copy;
};

Machine.prototype.asWorker = function() {
    return this.withRole("Worker");
};

Machine.prototype.asMaster = function() {
    return this.withRole("Master");
};

// Create n new machines with the same attributes.
Machine.prototype.replicate = function(n) {
    var i;
    var res = [];
    for (i = 0 ; i < n ; i++) {
        res.push(this.clone());
    }
    return res;
};

function Image(name, dockerfile) {
    this.name = name;
    this.dockerfile = dockerfile;
}

Image.prototype.clone = function() {
    return new Image(this.name, this.dockerfile);
}

function Container(image, command) {
    // refID is used to distinguish deployments with multiple references to the
    // same container, and deployments with multiple containers with the exact
    // same attributes.
    this._refID = uniqueID();

    this.image = image;
    if (typeof image === 'string') {
        this.image = new Image(image);
    }

    if (this.image.constructor !== Image) {
        throw new Error('bad image type');
    }

    this.command = command || [];
    this.env = {};
    this.filepathToContent = {};
}

// Create a new Container with the same attributes.
Container.prototype.clone = function() {
    var cloned = new Container(this.image.clone(), _.clone(this.command));
    cloned.env = _.clone(this.env);
    cloned.filepathToContent = _.clone(this.filepathToContent);
    return cloned;
};

// Create n new Containers with the same attributes.
Container.prototype.replicate = function(n) {
    var i;
    var res = [];
    for (i = 0 ; i < n ; i++) {
        res.push(this.clone());
    }
    return res;
};

Container.prototype.setEnv = function(key, val) {
    this.env[key] = val;
};

Container.prototype.withEnv = function(env) {
    var cloned = this.clone();
    cloned.env = env;
    return cloned;
};

Container.prototype.withFiles = function(fileMap) {
    var cloned = this.clone();
    cloned.filepathToContent = fileMap;
    return cloned;
};

Container.prototype.setHostname = function(h) {
    this.hostname = h;
};

Container.prototype.getHostname = function() {
    if (this.hostname === undefined) {
        throw new Error("no hostname");
    }
    return this.hostname + ".q";
};

var enough = { form: "enough" };
var between = invariantType("between");
var neighbor = invariantType("reachDirect");
var reachableACL = invariantType("reachACL");
var reachable = invariantType("reach");

function Assertion(invariant, desired) {
    this.form = invariant.form;
    this.nodes = invariant.nodes;
    this.target = desired;
}

function invariantType(form) {
    return function() {
        // Convert the arguments object into a real array. We can't simply use
        // Array.from because it isn't defined in Otto.
        var nodes = [];
        var i;
        for (i = 0 ; i < arguments.length ; i++) {
            nodes.push(arguments[i]);
        }

        return {
            form: form,
            nodes: nodes
        };
    };
}

function LabelRule(exclusive, otherService) {
    this.exclusive = exclusive;
    this.otherLabel = otherService.name;
}

function MachineRule(exclusive, optionalArgs) {
    this.exclusive = exclusive;
    if (optionalArgs.provider) {
        this.provider = optionalArgs.provider;
    }
    if (optionalArgs.size) {
        this.size = optionalArgs.size;
    }
    if (optionalArgs.region) {
        this.region = optionalArgs.region;
    }
    if (optionalArgs.floatingIp) {
      this.floatingIp = optionalArgs.floatingIp;
    }
}

function Connection(ports, to) {
    this.minPort = ports.min;
    this.maxPort = ports.max;
    this.to = to;
}

function Range(min, max) {
    this.min = min;
    this.max = max;
}

function Port(p) {
    return new PortRange(p, p);
}

var PortRange = Range;

function getDeployment() {
    return global._quiltDeployment;
}

// Reset global unique counters. Used only for unit testing.
function resetGlobals() {
    uniqueIDCounter = 0;
    labelNameCount = {};
}

module.exports = {
    Assertion,
    Connection,
    Container,
    Deployment,
    Image,
    LabelRule,
    Machine,
    MachineRule,
    Port,
    PortRange,
    Range,
    Service,
    createDeployment,
    getDeployment,
    githubKeys,
    publicInternet,
    enough,
    resetGlobals,
};
