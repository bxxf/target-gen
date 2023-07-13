#!/bin/sh
# Using specified languages and all EN variants
tgen loc=es,de,en,cz --en-all

# Using predefined brand language sets
tgen loc=BRANDNAME

# Using predefined brand language sets and custom attributes
tgen loc=BRANDNAME segment=SKU1,SKU2,SKU3