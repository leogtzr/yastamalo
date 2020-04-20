#!/bin/bash

readonly work_dir=$(dirname "$(readlink --canonicalize-existing "${0}")")
grep --fixed-strings '#yastamalo-v1.0' ~/.bashrc || {
    echo 'grep --extended-regexp --quiet "root" <<< "$(whoami)" || \
		"${HOME}/bin/yastamalo" -db ~/inputs/foods.db #yastamalo-v1.0' >> ~/.bashrc
}

exit 0

