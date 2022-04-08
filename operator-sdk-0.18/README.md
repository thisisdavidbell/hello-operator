# Creating an Operator with Operator SDK 0.18

This directory guides you through the process of creating a simple Kubernetes operator managed with OLM, for the Hello World httpserver go app previously created.

Note: it was created using operator sdk 0.18.2, which is no longer in support. There will be additional directories covering further tutorials to migrate this operator, and also to create an operator with the latest version of operator-sdk at the time.

## Links

- Operator tutorial: older doc used as primary source for this work:
  - https://github.com/operator-framework/getting-started
  - 

- Other links:
  - Newer operator-sdk doc: https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/
  - older operator-sdk doc: https://docs.openshift.com/container-platform/4.6/operators/operator_sdk/osdk-getting-started.html ( uses 0.19)

## 0. Prereqs

- containerised Hello World go application: [hello-app](../hello-app)
- operator-sdk: 0.18.2
- go: 1.16
- docker or podman
- Red Hat Openshift Container Platform: 4.6+

- Important Notes:
  -  if using Visual Studio Code: 
        - gopls language server expects either the go.mod file to be in the root of the workspace, or for you to open the dorectory containing your code and the go.mod file. This hello-operator repo is not setup like that, so you may need to open the hello-operator directly itself in a new go workspace for VS Code to work correctly, for example not present errors in main.go after step 1, such as:
```
could not import k8s.io.... (cannot find package...)
```
See: https://github.com/golang/tools/blob/master/gopls/doc/workspace.md
  - vendor directory
    - the vendor directory has been added to the .gitignore dir.
    - if cloning this repo, run `go mod vendor` and not the vendor directory is created and populated. 


Hopefully for the more up to date operator with a 1.18+ go version, the workspace can be made to work correctly with go.work files.

## 1. Create initial project

- Run: `operator-sdk new hello-operator`
- Output:
```
INFO[0000] Creating new Go operator 'hello-operator'.   
INFO[0000] Created go.mod                               
INFO[0000] Created tools.go                             
INFO[0000] Created cmd/manager/main.go                  
INFO[0000] Created build/Dockerfile                     
INFO[0000] Created build/bin/entrypoint                 
INFO[0000] Created build/bin/user_setup                 
INFO[0000] Created deploy/service_account.yaml          
INFO[0000] Created deploy/role.yaml                     
INFO[0000] Created deploy/role_binding.yaml             
INFO[0000] Created deploy/operator.yaml                 
INFO[0000] Created pkg/apis/apis.go                     
INFO[0000] Created pkg/controller/controller.go         
INFO[0000] Created version/version.go                   
INFO[0000] Created .gitignore                           
INFO[0000] Validating project                           
go: github.com/operator-framework/operator-sdk@v0.18.2: missing go.sum entry; to add it:
        go mod download github.com/operator-framework/operator-sdk
FATA[0000] failed to exec []string{"go", "build", "./..."}: exit status 1 
```
- `cd hello-operator`
- `go mod download github.com/operator-framework/operator-sdk`
- `go mod tidy`

- view the directory structure, especially:
  - build/*/Dockerfile
  - cmd/manager/main.go
  - deploy/*
  - pkg/*
  - version/*

Notes:
- the key aspects of the project,including operator scope: https://github.com/operator-framework/getting-started#manager (or see Appendix 1)

## 2. Add a new Custom Resource Definition
Add a new Custom Resource Definition (CRD) API called hello, with APIVersion `thisisdavidbell.example.com/v1alpha1` and Kind `Hello`.

```
$ operator-sdk add api --api-version=thisisdavidbell.example.com/v1alpha1 --kind=Hello

INFO[0000] Generating api version thisisdavidbell.example.com/v1alpha1 for kind Hello. 
INFO[0000] Created pkg/apis/thisisdavidbell/group.go    
INFO[0005] Created pkg/apis/thisisdavidbell/v1alpha1/hello_types.go 
INFO[0005] Created pkg/apis/addtoscheme_thisisdavidbell_v1alpha1.go 
INFO[0005] Created pkg/apis/thisisdavidbell/v1alpha1/register.go 
INFO[0005] Created pkg/apis/thisisdavidbell/v1alpha1/doc.go 
INFO[0005] Created deploy/crds/thisisdavidbell.example.com_v1alpha1_hello_cr.yaml 
INFO[0005] Running deepcopy code-generation for Custom Resource group versions: [thisisdavidbell:[v1alpha1], ] 
INFO[0014] Code-generation complete.                    
INFO[0014] Running CRD generator.                       
INFO[0015] CRD generation complete.                     
INFO[0015] API generation complete.                     
INFO[0015] API generation complete.   
```
View the new code in:
- deploy/crds - crd and example cr
- pkg/apis/thisisdavidbell/v1alpha1 - operator models, etc

# 3. Update the Hello model

NEXT: make some additions to the model, to allow demonstration of openapi schema validation, webhook, operator consuming and applying values from cr, etc...

---

# Appendix

## 1. Key aspects of operator project

Content from https://github.com/operator-framework/getting-started#manager in case it ceases to exist:

---

**Manager**

The main program for the operator cmd/manager/main.go initializes and runs the Manager.

The Manager will automatically register the scheme for all custom resources defined under pkg/apis/... and run all controllers under pkg/controller/....

The Manager can restrict the namespace that all controllers will watch for resources:
```
mgr, err := manager.New(cfg, manager.Options{
	Namespace: namespace,
})
```
By default this will be the namespace that the operator is running in. To watch all namespaces leave the namespace option empty:
```
mgr, err := manager.New(cfg, manager.Options{
	Namespace: "",
})
```

