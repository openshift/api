This content was generated from a live cluster.
It may change over time and is not guaranteed to be stable.
This is a useful starting point for understanding the cert chains in openshift used to secure kubernetes.

1. Build an image to collect the certs, keys, and ca bundles from the host.
   1. Something like `docker build pkg/cmd/locateinclustercerts/ -t docker.io/deads2k/cert-collection:latest -f Dockerfile`
   2. Push to dockerhub
2. Gather data.
   1. `oc adm inspect clusteroperators` -- this will gather all the in-cluster certificates and ca bundles
   2. run pods on the masters.  Something like `oc debug --image=docker.io/deads2k/cert-collection:08 node/ci-ln-z2l4snt-f76d1-prqp5-master-2`
   3. in those pods, run `master-cert-collection.sh` to collect the data from the host.  Leave the pod running after completion.
   4. pull the on-disk data locally. Something like `oc rsync ci-ln-z2l4snt-f76d1-prqp5-master-2-debug:/must-gather ../sample-data/master-2/`
3. Be sure dot is installed locally
4. Run code like `kubectl-dev_tool certs locate-incluster-certs --local -f ../sample-data/ --additional-input-dir ../sample-data/ -odoc` to produce the doc.
