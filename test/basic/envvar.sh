FOO=Hello fillin BAR=beautiful BAZ=world! sh -c 'echo {{foo}}, {{bar}} {{baz}}' <<'EOF'
$FOO
$BAR
$BAZ
EOF
