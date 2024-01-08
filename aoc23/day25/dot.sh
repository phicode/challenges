#!/bin/env bash

fdp -Tpng example.dot > example.fdp.png
dot -Tpng example.dot > example.dot.png

fdp -Tpng input.dot > input.fdp.png

dot -Kneato -Tsvg input.dot > input.neato.svg

# input too large
#dot -Tpng input.dot  > input.dot.png
