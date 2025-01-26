# API Generation
The operator builds the api and crd based on the Endpoint type struc from Gatus itself, using a very alpha generator which definitely could use a rewrite as it contains many exceptions and other smaller hacks to at least continue my development time on the operator itself.
The API is generated using the following commands:

```bash
- go run gen.go
- deepcopy-gen --v=9 ./api/v1alpha1
- controller-gen crd:crdVersions=v1 paths=./api/v1alpha1 output:crd:dir=./crd/bases
```

## Change in fields
- HttpClient in ClientConfig is ommited (can't deep copy it)
- ProviderOverride is casted to map[string]json.RawMessage
- Any time.Duration is casted to string for the CRD (otherwise it will be an int64)
- Everything is forced to be omitempty so the configmap will only contain the fields that are actually set