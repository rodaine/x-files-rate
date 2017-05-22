#!/usr/bin/env bash
# plot data in results file
set -e

target=${1:-hello}
file="results/${target}.csv"

mkdir -p ./results/plots
output="results/plots/${target}.svg"
rm $output || true

opts="ps 1 lt 1 lw 2 with lines"
opts1="lc rgb '#8b1a0e' pt 1 ${opts}"
opts2="lc rgb '#5e9c36' pt 6 ${opts}"

gnuplot -p -e "set datafile separator ',';\
    set terminal svg size 800,600 fname 'Verdana' fsize 10;\
    set output '${output}';\
    set style line 11 lc rgb '#808080' lt 1;\
    set border 3 back ls 11;\
    set tics nomirror;\
    set style line 12 lc rgb '#808080' lt 0 lw 1;\
    set grid back ls 12;\
    set xlabel 'RPS';\
    set ylabel 'ratio';\
    set yrange [0:*]; \
    set multiplot layout 2,2 rowsfirst;\
    plot '${file}' using 1:2 title 'SR' ${opts1};\
    set ylabel 'ms';\
    plot \
    '${file}' using 1:4 title 'Total P95' ${opts1},\
    '${file}' using 1:3 title 'Total P50' ${opts2};\
    plot \
    '${file}' using 1:6 title 'Success P95' ${opts1},\
    '${file}' using 1:5 title 'Success P50' ${opts2};\
    plot \
    '${file}' using 1:8 title 'Failure P95' ${opts1},\
    '${file}' using 1:7 title 'Failure P50' ${opts2}"
