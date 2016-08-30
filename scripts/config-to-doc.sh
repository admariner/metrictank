#!/bin/bash

# turns the metrictank-sample.ini into a nice markdown document.
# headers like h3 and h2 are printed as-is, and the lines between them (config items and their comments)
# are wrapped in ``` blocks. t=code means we're in such a block.

cat << EOF
# Config

Metrictank comes with an [example config file](https://github.com/raintank/metrictank/blob/master/metrictank-sample.ini)
which is well documented.

It is replicated below.  
This is generated by [config-to-doc](https://github.com/raintank/metrictank/blob/master/scripts/config-to-doc.sh)


---


EOF



t=
while read line; do
	# skip empty lines
	if [ "$line" == "" ]; then
		continue
	fi
	if [[ $line =~ ^### ]]; then
		if [ "$t" == code ]; then
			echo '```'
			echo
		fi
		echo $line
		t=h3
	elif [[ $line =~ ^## ]]; then
		if [ "$t" == code ]; then
			echo '```'
			echo
		fi
		echo $line
		t=h2
	else
		if [ "$t" == h2 -o "$t" == h3 ]; then
			echo -e '\n```'
			t=code
		fi
		# lines that start with a pound are fine within code blocks,
		# but outside of code blocks, would be shown as headers, which is not how they are intended
		# in the source file they are just regular comments too, so take away their #
		#if [[ $line =~ ^# ]]; then
		if [[ "$t" != code ]]; then
			sed -e 's/^# //' -e 's/$/  /'<<< "$line"
		else
			echo $line
		fi
	fi
done < metrictank-sample.ini
