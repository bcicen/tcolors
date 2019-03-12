#!/bin/bash

COLOR_OUTPUT=1
STEPS=25

function _round() { printf "%.0f " $@; }

function _usage() {
  echo "usage: $0 [OPTION]... [INPUT]..."
  echo "generate a color scale from a given base color"
  echo -e "\n  -n disable colorized output"
  echo -e "  -s<n> set number of steps in scale to <n> (default $STEPS) "
  echo
}

function _rgb2hex() { printf '%02x' $1 $2 $3; }

function _hex2rgb() {
  in=${1#\#} # strip leading #, if any
  r=${in:0:2} g=${in:2:2} b=${in:4:2}
  echo "$((16#$r)) $((16#$g)) $((16#$b))"
}

# output an RGB color code
function _output() {
  r=$1 g=$2 b=$3
  txt="$(printf '%03d ' $r $g $b)\t$(_rgb2hex $1 $2 $3)"
  if (($COLOR_OUTPUT)); then
    echo -e "\033[38;2;${r};${g};${b}m${txt}\033[0m";
  else
    echo -e "${txt}"
  fi
}

function _scale() {
  rgb=($1 $2 $3)
  to_rgb=($4 $5 $6)

  rates=()
  for idx in 0 1 2; do
    r=0
    if [[ ${rgb[$idx]} -gt ${to_rgb[$idx]} ]]; then
      r=$(echo "scale=2; (${to_rgb[$idx]} - ${rgb[$idx]}) / ($STEPS+1)" | bc)
    elif [[ ${rgb[$idx]} -lt ${to_rgb[$idx]} ]]; then
      r=$(echo "scale=2; (${to_rgb[$idx]} - ${rgb[$idx]}) / ($STEPS+1)" | bc)
    fi
    rates+=($r)
  done
  #echo "rates: ${rates[@]}"

  i=0
  while [[ $i -lt $STEPS ]]; do
    for idx in 0 1 2; do
      rgb[$idx]=$(echo "${rgb[$idx]} + ${rates[$idx]}" | bc)
    done
    _output $(_round ${rgb[@]})
    let i++
  done
}

opts=() args=()

for x in $@; do
  case $x in
    -*) opts+=($x) ;;
    *) args+=($x) ;;
  esac
done

for opt in ${opts[@]}; do
  case $opt in
    "-n") COLOR_OUTPUT=0 ;;
    -s*) STEPS=${opt#-s};;
    "-h") _usage; exit 0 ;;
    *) echo "unknown option \"$opt\""; exit 1 ;;
  esac
done

[[ "$STEPS" -le 0 ]] && {
  echo "step value must be greater than 0"
  exit 1
}

if [[ ${#args[@]} -eq 2 ]]; then
  RGB=($(_hex2rgb ${args[0]}))
  TO_RGB=($(_hex2rgb ${args[1]}))
  _output ${RGB[@]}
  _scale ${RGB[@]} ${TO_RGB[@]}
  _output ${TO_RGB[@]}
  exit 0
fi
  
RGB=($(_hex2rgb ${args[@]}))

_scale ${RGB[@]} 0 0 0 | tac
_output ${RGB[@]}
_scale ${RGB[@]} 255 255 255
