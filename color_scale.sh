#!/bin/bash

COLOR_OUTPUT=1
STEPS=25

__cscale_usage() {
  echo "usage: $0 [OPTION]... [COLOR1] [COLOR2]"
  echo "generate a color gradient scale from given base color(s)"
  echo -e "\noptions:"
  echo "  -n    disable colorized output"
  echo "  -s<n> set number of steps in scale to <n> (default $STEPS) "
  echo -e "\nexamples:"
  echo -e "  $0 4be397         generate gradient scale of darkest to lightest of 4be397"
  echo -e "  $0 4be397 00add8  generate gradient scale from 4be397 to 00add8"
}

__cscale_hex2rgb() {
  in=${1#\#} # strip leading #, if any
  r=${in:0:2} g=${in:2:2} b=${in:4:2}
  echo "$((16#$r)) $((16#$g)) $((16#$b))"
}

# output an RGB color code
__cscale_output() {
  r=$1 g=$2 b=$3
  txt="$(printf '%03d ' $r $g $b)\t$(printf '%02x' $1 $2 $3)"

  if (($COLOR_OUTPUT)); then
    echo -e "\033[38;2;${r};${g};${b}m${txt}\033[0m";
  else
    echo -e "${txt}"
  fi
}

__cscale_scale() {
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

  i=0
  while [[ $i -lt $STEPS ]]; do
    for idx in 0 1 2; do
      rgb[$idx]=$(echo "${rgb[$idx]} + ${rates[$idx]}" | bc)
    done
    __cscale_output $(printf "%.0f " ${rgb[@]})

    let i++
  done
}

# skip action if script is sourced
(return 0 2>/dev/null) && return

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
    "-h") __cscale_usage; exit 0 ;;
    *) echo "unknown option \"$opt\""; exit 1 ;;
  esac
done

[[ "$STEPS" -le 0 ]] && {
  echo "step value must be greater than 0"
  exit 1
}

if [[ ${#args[@]} -eq 2 ]]; then
  RGB=($(__cscale_hex2rgb ${args[0]}))
  TO_RGB=($(__cscale_hex2rgb ${args[1]}))
  __cscale_output ${RGB[@]}
  __cscale_scale ${RGB[@]} ${TO_RGB[@]}
  __cscale_output ${TO_RGB[@]}
  exit 0
fi
  
RGB=($(__cscale_hex2rgb ${args[@]}))
STEPS=$((STEPS/2))
__cscale_scale ${RGB[@]} 0 0 0 | tac
__cscale_output ${RGB[@]}
__cscale_scale ${RGB[@]} 255 255 255
