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

### Intent
The intent of the application redesign is to address inconsistencies observed in day to day use of the appstore api with `maintainability` , `release processes, updates` and `outside entity adoption` being a core focus. 

Due care was put into the design to support the helxplatform's adoption and smooth usage outside of UNC, with special focus being from perspective of a system engineer or admin who would be tasked with the platform up in their own environment. This means api design should support user customizations via restful api calls (ie branding and adding new apps).

**A NOTE ON ADDING APPLICATIONS: While the concept is simple at face value, this is also a large consideration as data will likely need to persist. The current state of this applciation does not address physical volumes (pv) or pvc creation.

### Compatibility
The Gopherhelx api did not give consideration to the current `helx-ui` therefore integration of this api will necessitate reworking of both. 

## The Api at a Glance
The functionality of the api endpoints as they relate to the end user perspective within the user interface (ui) is depicted below.

![alt text](https://github.com/joshua-seals/gopherhelx/blob/readme-illustration/.readme-images/images/app-list-endpoints.png?raw=true)

## Development:
The api service is designed to be driven by the Makefile.
 
