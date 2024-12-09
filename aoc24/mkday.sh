#!/usr/bin/env bash

usage () {
  echo "usage: $0 <day-number>"
  echo "example: $0 5"
  exit 1
}

if [[ $# -ne 1 ]]; then
  usage
fi

script_path=$(readlink -e "$0")
base_dir=$(dirname "$script_path")

day=$1
xx_day=$1
if [[ ! "$day" =~ ^[0-9]{1,2}$ ]]; then
  echo "ERROR: invalid day number"
  exit 1
fi
if [[ ${#day} -eq 1 ]]; then
  xx_day="0${day}"
fi

template_dir="${base_dir}/tmpl"
dest_dir="${base_dir}/day${xx_day}"

if [[ ! -d "${template_dir}" ]]; then
  echo "ERROR: template directory not found: ${template_dir}"
  exit 1
fi
if [[ -d "${dest_dir}" ]]; then
  echo "ERROR: output directory already exists: ${dest_dir}"
  exit 1
fi

cp -r "${template_dir}" "${dest_dir}"
sed \
  -e "s/#X#/${day}/g" \
  -e "s/#XX#/${xx_day}/g" \
  "${dest_dir}/dayXX.go" > "${dest_dir}/day${day}.go"
rm "${dest_dir}/dayXX.go"
