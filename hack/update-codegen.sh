#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

# Build codegen-crds when it's not present and not overridden for a specific file.
if [ -z "${CODEGEN:-}" ];then
  ${TOOLS_MAKE} codegen
  CODEGEN="${TOOLS_OUTPUT}/codegen"
fi

# This runs the codegen utility against the entire set of API types.
# It has two possible args, the GENERATOR and the EXTRA_ARGS.
# GENERATOR is the name of the generator to run, which could be one of compatibility,
# deepcopy, schemapatch or swaggerdocs.
# EXTRA_ARGS are additional arguments to pass to the generator, usually this would be
# --verify so that the generator verifies the output rather than writing it.
"${CODEGEN}" ${GENERATOR:-} --base-dir "${SCRIPT_ROOT}" -v 1 ${EXTRA_ARGS:-}
