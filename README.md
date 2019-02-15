# openshift template operator

A sample operator that uses integreatly openshift library to process a template and deploy its objects.

Requires operator-sdk 0.4+

## Usage

Create a new openshift project:

```sh
oc new-project template-operator
```

Apply role, service account and crd files:

```sh
oc apply -f deploy/service_account.yaml
oc apply -f deploy/role.yaml
oc apply -f deploy/role_binding.yaml
oc apply -f deploy/crds/odra_v1alpha1_okdtemplate_crd.yaml
```

Deploy the operator:

```
oc apply -f deploy/operator.yaml
```

Apply the template cr so you can provision the sample template:

```sh
oc apply -f oc apply -f deploy/crds/odra_v1alpha1_okdtemplate_cr.yaml
```

It will deploy a sample integreatly operator in a few minutes depending on your
internet connection.

## The Custom Resource

Sample custom resource

```yaml
apiVersion: odra.org/v1alpha1
kind: OKDTemplate
metadata:
  #name your cr as you see fit
  name: example-okdtemplate
spec:
  source:
    # local template file to be used - image needs to be re-build if the resource folder changes
    local: webapp/latest.yaml
    # remote template url - not supported yet
    #remote: ''
  parameters:
    #openshift template parameters - those will be direcrly used when processing the template
    OPENSHIFT_OAUTHCLIENT_ID: 'webapp'
    OPENSHIFT_HOST: ''
    SSO_ROUTE: ''
    WEBAPP_IMAGE: 'quay.io/integreatly/tutorial-web-app'
    WEBAPP_IMAGE_TAG: latest
```

Deleting the CR will result in deleting all runtime objects generated/created by the operator.

## Building

You can add any folders and files in the `res` folder.

You need to run the `packr` cli to re-generate the res folder (https://github.com/gobuffalo/packr).

Once everything is in place just build your operator and push it to a container registry of your choice:

```
operator-sdk build quay.io/my_org/image_name:image_tag
docker push quay.io/my_org/image_name:image_tag
```

Do not forget to change the image name in the operator.yaml file.
