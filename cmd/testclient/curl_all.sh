#!/bin/sh

models=(
    "qwen-max"
    "moonshot-v1-8k"
    "glm-4-flash"
    "doubao-pro-4k"
    "hunyuan-lite"
)
for model in "${models[@]}"
do
    echo "{\"model\":\""$model"\",\"prompt\":\"介绍一下你自己\"}"
    curl -X POST -d "{\"model\":\""$model"\",\"prompt\":\"介绍一下你自己\"}" http://localhost:38120/completions
    echo
    echo "{\"model\":\""$model"\",\"prompt\":\"介绍一下你自己\",\"stream\":true}"
    curl -X POST -d "{\"model\":\""$model"\",\"prompt\":\"介绍一下你自己\",\"stream\":true}" http://localhost:38120/completions
    echo
done
