# Gatus Operator
<sup><sup>Not directly affiliated with gatus.io</sup></sup>

**State**: Alpha (in development, use at your own risk)

## Install
1. Install the crds
```bash
kubectl apply -f crds/bases/gatus-operator.aumer.io_gatuses.yaml
```
2. Install RBAC
```bash
kubectl apply -f rbac.yaml
```
3. Run container as per usual

## Configuration
The operator is configured using a configfile which is mounted into the container at /config/config.yaml
The configfile is a yaml file with the following structure:


```yaml
---
operator:
  k8s-sidecar-annotation: gatus.io/enabled
defaults:
  global:
    interval: 5m
    conditions: ["[STATUS] == 200"]
    client:
      dns-resolver: tcp://1.1.1.1:53
  guarded:
    interval: 2m
    ui:
      hide-url: true
      hide-hostname: true
  infrastructure:
    interval: 30m
    client: {}
```

### Configuration Options

| Option | Description | Default |
| --- | --- | --- |
| operator | Base configuration for the operator |  |
| defaults | Default values for the Gatus CRD grouped |  |

#### Operator
| Option | Description | Default |
| --- | --- | --- |
| k8s-sidecar-annotation | The annotation used in the k8s-sidecar, if used to automatically map gatus configmaps |  |

#### Defaults
| Option | Description | Default |
| --- | --- | --- |
| global | All the default endpoints values regardless of their group |  |
| \<group\> | Endpoint default values for the specified group |  |

<sup>\* Global and/or group defaults are not required<br />
** The defaults are applied as follows: global -> \<group\> -> endpoint<br >
*** For all the possible endpoint values refer gatus config or the CRD</sup>

## Usage
Example: resources/test.yaml
```yaml
---
apiVersion: gatus.io/v1alpha1
kind: Gatus
metadata:
  name: gatus-test
  namespace: observability
spec:
  endpoint:
    enabled: true
    name: gatus-test
    url: https://example.com
    interval: 1m
    conditions: ["[STATUS] == 200"]
```

The api is generated from the gatus package as well with some minor tweaks (to make it working).
Most options you use in your normal gatus configs, you should be able to use like this.

## API Generation
The operator builds the api and crd based on the Endpoint type struc from Gatus itself, using a very alpha generator which definitely could use a rewrite as it contains many exceptions and other smaller hacks to at least continue my development time on the operator itself.
The API is generated using the following commands:

```bash
- go run gen.go
- deepcopy-gen --v=9 ./api/v1alpha1
- controller-gen crd:crdVersions=v1 paths=./api/v1alpha1 output:crd:dir=./crd/bases
```

### Change in generated fields
- HttpClient in ClientConfig is ommited (can't deep copy it)
- ProviderOverride is casted to map[string]json.RawMessage
- Any time.Duration is casted to string for the CRD (otherwise it will be an int64)
- Everything is forced to be omitempty so the configmap will only contain the fields that are actually set
