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

### Important Notes:
  -  if using Visual Studio Code: 
        - gopls language server expects either the go.mod file to be in the root of the workspace, or for you to open the dorectory containing your code and the go.mod file. This hello-operator repo is not setup like that, so you may need to open the hello-operator directly itself in a new go workspace for VS Code to work correctly, for example not present errors in main.go after step 1, such as:
```
could not import k8s.io.... (cannot find package...)
```
    - All operator-sdk commands will need to be run from a new VS workspace based on the `operator-sdk-0.18` dir
    - See: https://github.com/golang/tools/blob/master/gopls/doc/workspace.md


  - vendor directory
    - the vendor directory has been added to the .gitignore dir.
    - if cloning this repo, run `go mod vendor` in the `hello-operator` directory and note the vendor directory is created and populated. 


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

In order to demonstrate various functionality of an operator, we will have 3 items in the spec:

spec:
- version - string - which hello image tag to deploy
- repeat - int - how many time to say hello
- verbose - bool - whether to include second line of text in output

Notes:
- hello:v1.0 doesn't use these yet, we will fix this.
- our initial operator won't use these yet
- we may add a status field later.

a. To do this, update `pkg/apis/thisisdavidbell/v1alpha1/hello_types.go` to look like:

```
type HelloSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Version - what version of hello to use - this is the hello image tag to use
	Version string `json:"version"`

	// Repeat - how many times to say hello
	Repeat int32 `json:"repeat"`

	// Verbose - whether to output additional line of text
	Verbose bool `json:"verbose"`
}
```

Note if you pushed your changes to git previously, you now only have the _types file changed:

```
$ git status

...
$ git status
On branch main
Your branch is up to date with 'origin/main'.

Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git restore <file>..." to discard changes in working directory)
        modified:   hello-operator/pkg/apis/thisisdavidbell/v1alpha1/hello_types.go
```

b. Add OpenApiV3Schema validation
Follow the link in the _types comment to view the supported OpenApiV3Schema validation that can be applied, using the kubebuilder annotations. Add appropriate validation. (Later we may add a full ValidatingWebhook for more control). E.g.:
```
type HelloSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Version - what version of hello to use - this is the hello image tag to use
	// +kubebuilder:validation:MaxLength=10
	// +kubebuilder:validation:MinLength=2
	Version string `json:"version"`

	// Repeat - how many times to say hello
	// +kubebuilder:validation:Maximum=5
	// +kubebuilder:validation:Minimum=1
	Repeat int32 `json:"repeat"`

	// Verbose - whether to output additional line of text
	Verbose bool `json:"verbose"`
}
```

Note: note, later if making further updates to the validation in  types file hello_types.go, you only need to regenerate the crd and redeploy the crd file

b. Update the generated code

```
operator-sdk generate k8s
```

Note that this appears to cause not changes at this point:

c. Update the crd:

```
operator-sdk generate k8s
```

Note the crd gets updated to include these new types

# 4. Manually update cr 

ToDo: should this get automatically updated, or is there a command to do it?

Updated cr to be:
```
apiVersion: thisisdavidbell.example.com/v1alpha1
kind: Hello
metadata:
  name: example-hello
spec:
  # Add fields here
  version: "v1.0"
  repeat: 1
  verbose: true
```

# 5. Add a new controller to watch and reconcile the Hello resource

```
operator-sdk add controller --api-version=thisisdavidbell.example.com/v1alpha1 --kind=Hello
```

Note this adds the files `pkg/controller/add_hello.go` and `pkg/controller/hello/hello_controller.go`.

f. Update the controller to use the correct image

Note: we may update this to use a deployment, as well as deploy the correct service and route later.

```
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "hello",
					Image:   "example.com/drb/hello:v1.0",
					Command: []string{"./hello"},
				},
			},
		},
```

Note I dont want to reference an internal registry in git, so we will use an ImageContentSourcePolicy to redirect example.com.


# 6. Build the operator

Ensure you have set env vars:
- `IRHOSTNAME` - Image registry Hostname
- `IRNAMESPACE` - Namespace in image registry
- `IRUSER` - Image registry user name
- `IRPASSWORD` - Image resgiry password

Run:
- `docker login -u $IRUSER -p $IRPASSWORD $IRHOSTNAME`

```
operator-sdk build $IRHOSTNAME/$NAMESPACE/hello-operator:v0.0.1
```

Push the image:
- `docker push $IRHOSTNAME/$NAMESPACE/hello-operator:v0.0.1`


# 7. Register the crd

```
kubectl create -f deploy/crds/thisisdavidbell.example.com_hellos_crd.yaml 
```
Note the cluster now understands what a `Hello` kind is - you can search for objects of that type
```
$ oc get somethingthatdoesntexist
error: the server doesn't have a resource type "somethingthatdoesntexist"

$ oc get hello
No resources found in default namespace.
```

# 8. Deploy operator

Update image in operator.yaml to actual value of echo "$IRHOSTNAME/$NAMESPACE/hello-operator:v0.0.1"
Note: example in git is incorrect

Then run:

oc create -f deploy/service_account.yaml
oc create -f deploy/role.yaml
oc create -f deploy/role_binding.yaml
oc create -f deploy/operator.yaml
---

# 9. Create Service and Route

Later we will get the oeprator to create the service and the route. For now manuallly create them.

Service:
```
apiVersion: v1
kind: Service
metadata:
  name: hello1-service
  namespace: drb-hello-operator
spec:
  selector:
    app: example-hello
  ports:
    - protocol: TCP
      port: 8080
      targetPort: 8080
```

Route:
Main details:
```
spec:
  host: hello1.drb-hello-operator.apps.RESTOFCLUSTERHOSTNAME
  path: /hello
  to:
    kind: Service
    name: hello1-service
    weight: 100
  port:
    targetPort: 8080
  wildcardPolicy: None
```

# 10. Deploy cr

```
oc create -f deploy/crds/thisisdavidbell.example.com_v1alpha1_hello_cr.yaml
```

Test:

In browser or curl, call:
```
http://hello1.drb-hello-operator.apps.RESTOFCLUSTERHOSTNAME/hello
```

Note if you set repeat to 10, for example, you will see the validation failure:

```
$ oc create -f deploy/crds/thisisdavidbell.example.com_v1alpha1_hello_cr.yaml  

The Hello "example-hello" is invalid: spec.repeat: Invalid value: 10: spec.repeat in body should be less than or equal to 5

---

Next:
- look for better way to override internal registry in code. `kustomize`??
- have hello use the repeat and verbose fields - use env vars and create a local file?
- have operator apply version, repeat and verbose fields correctly.
- convert to deployment following memcache example code here: https://docs.openshift.com/container-platform/4.6/operators/operator_sdk/osdk-getting-started.html
- reconcile service in operator
- repeat with new code for route
- tidy up docs
- move onto olm

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

