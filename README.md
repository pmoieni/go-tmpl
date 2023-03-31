# kratos-test
a Go API that uses ory kratos for authentication

## How to start Kratos?

- rename `identity.schema.example.json` to `identity.schema.json`
- rename `kratos.example.yaml` to `kratos.yaml`
- fill out the absolute path in `kratos.yaml` file
- run the following command

```
 podman run -p 4433:4433 -p 4434:4434 --name kratos --mount type=bind,src=<absolute-path-to-project>/kratos-test,target=<absolute-path-to-project>/kratos-test localhost/oryd/kratos:
v0.11.1 serve --config <absolute-path-to-project>/kratos-test/kratos.yaml
```
