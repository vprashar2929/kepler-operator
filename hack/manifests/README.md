# Hack manifests

This directory contains various manifests used for testing and debugging purposes.

## Node exporter

* Create namespace

```sh
oc create namespace node-exporter
```

* Add scc to user

```sh
oc adm policy add-scc-to-user privileged -z node-exporter -n node-exporter
```

* Apply the manifest

```sh
oc apply -f node-exporter.yaml
```

## Stress

* Create namespace

```sh
oc create namespace stress
```

* Navigate to the Kepler directory where the stressors script/any stress script is present

* Create a config map with the stress script

```sh
 oc create configmap stress-ng-script --from-file=stressor.sh -n stress
```

* Apply the manifest

```sh
oc apply -f stress.yaml
```
