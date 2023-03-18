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

The design of this application uses a base from ardan labs ultimate service and tries to stay consistent with package oriented design philosophy, ie business logic is contained in the business package, while foundational code will be held in the foundation package. This should allow for independent module testing throughout the life of the project. 

The intent of the application redesign is to address inconsistencies observed in day to day use of the appstore api with `maintainability` , `release processes and updates` and `outside entity adoption` being a core focus. 

Due care was put into the desired functionality for other entities to adopt and smoothly use the gopherhelx api. This means allowing for ease of installation of new apps and supporting easy branding. 

A NOTE ON ADDING APPLICATIONS: While the concept is simple at face value, this is also a large consideration as data will likely need to persist. The current state of this applciation does not address physical volumes (pv) or pvc creation.

This api did not allow the current ui to inform decisions about routes and functionality, but rather is a reworking of both. The api should inform the ui.
The functionality as it relates from the end user perspective at the ui level is show below.

![alt text](https://github.com/joshua-seals/gopherhelx/blob/main/zarf/images/app-list-endpoints.png?raw=true)


 
