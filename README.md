# GopherHelx
Reimagining the Helxplatoform appstore api in golang.

Gopherhelx is a proof of concept to hopefully demonstrate a more 
simple approach to the helxplatform api and is focused on current usage
of the helxplatform in kubernetes.

Gopherhelx api aims to be the core functionality of the helxplatform.
This application assumes: 
  - Micro service patterns are followed throughout the platform (ie. ui and api are split into independent parts)
  - Applications are described in kubernetes manifests.
  - The api handles all kubernetes transactions and routing.

## Design
The design of this application uses a base design gleened from ![ardan-labs-service](https://github.com/ardanlabs/service) and tries to stay consistent with ![package oriented design philosophies](https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html) to aid in developer mental models of the architecture and guide development efforts within the project. Meaning, "business logic" is contained in the `business` package, while foundational code will be held in the `foundation` package, and the core api service will be held in the `services` package. 

This method of organization and visualization is similar to how the OSI Model and TCP/IP models are used by engineers to speak conceptionally about different levels of a system. Additionally package oriented design allows for independent module testing.

![alt text](https://github.com/joshua-seals/gopherhelx/blob/readme-illustration/.readme-images/images/package-oriented-design.png?raw=true)

### Intent
The intent of the application redesign is to address inconsistencies observed in day to day use of the appstore api with `maintainability` , `release processes, updates` and `outside entity adoption` being a core focus. 

Due care was put into the design to support the helxplatform's adoption and smooth usage outside of UNC, with special focus being from perspective of a system engineer or admin who would be tasked with the platform setup in their own environment. This means api design should support user customizations via restful api calls (ie branding and adding new apps).

**A NOTE ON ADDING APPLICATIONS: While the concept is simple at face value, this is also a large consideration as data will likely need to persist. The current state of this applciation does not address physical volumes (pv) or pvc creation.

### Compatibility
The Gopherhelx api did not give consideration to the current `helx-ui` therefore integration of this api will necessitate reworking of both. 

## The Api at a Glance
The functionality of the api endpoints as they relate to the end user perspective within the user interface (ui) is depicted below.

![alt text](https://github.com/joshua-seals/gopherhelx/blob/readme-illustration/.readme-images/images/app-list-endpoints.png?raw=true)

![alt text](https://github.com/joshua-seals/gopherhelx/blob/readme-illustration/.readme-images/images/dashboard-list-endpoints.png?raw=true)

## Development:
The api service is designed to be driven by the ![Makefile](https://github.com/joshua-seals/gopherhelx/blob/main/Makefile). Following this workflow pattern will ensure consistency for all developers. Therefore updates and maintenance to ensure the makefile is current and consistent with the state of the application is critical. 

### Setup
Run `make help` to see a list of available commands within the makefile. 

If on mac, running `make dev.setup.mac` will install most of the needed tools for development.

After setting up, run `make image` to build the docker image of the gopherhelx api, then `make kind-up` will bring up a new kind cluster. 

##### NOTE: If you have a ~/.kube/config already - you will want to rename it, as the dynamic creation of kind or minikube cluster overwrites the config file by default. 

After the kind cluster has started:
- Load the image into the cluster with `make kind-load` 
- Add permissions so api can create pods/deployments `make kind-default-CRB`
- Lastly apply the zarf/manifests with `make kind-apply`

At this point, the cluster should be up and ready for testing.

### Development Patterns
While developing and testing inside of the `/app` folder, ie the application code, to rebuild, load, restart appstore-api with new image, you can run `make kind-update`.

If doing work inside of the `/zarf` folder, ie editing dockerfiles or kubernetes manifests, you will need to run `make kind-update-apply` in order to reload the manifests.

When in doubt, clean it out 🧹🫧🧼 with `make kind-down`. 
Follow the `make image | kind-up | kind-load | kind-apply` pattern after deleting the cluster with kind-down.

## And Yet

This api is a good foundation, and yet the most critical aspect has not been addressed, the `service mesh`. Once the Deployment and Services are created for the applications, we still need to `dynamically` inject routing into them to present to the end user, as well as keep a running map of those users and their services. Currently, ![Consul](https://developer.hashicorp.com/consul/docs/connect) is the top consideration for service mesh due to it's heavy focus on opensource, educational and thorough app documentation. 


## Additionally 

These features will also need to be addressed:
- Validation 
- Testing (Module Testing was not prioritized)
- Middleware (catching panics, logging persistance, observability)
- User Persistence
- Authentication
- Authorization

Currently, there are two paths for authorization, but they should be prioritized after the service mesh is established. 
