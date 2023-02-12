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
 
