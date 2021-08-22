#!/bin/bash

if [ "$1" == "f" ]; then
  echo "format proto..."
  for file in `ls ./inner/*.proto`
  do
      clang-format -i -style="{AlignConsecutiveAssignments: true,AlignConsecutiveDeclarations: true,AllowShortFunctionsOnASingleLine: None,BreakBeforeBraces: GNU,ColumnLimit: 0,IndentWidth: 4,Language: Proto}" $file
  done
fi

binpath=../../../../protoc
$binpath/protoc-mac --plugin protoc-gen-go=$binpath/protoc-gen-go-mac -I=./ --go_out=../inner_message/ ./inner/*.proto