## Ollama and LLAMA embedding output difference


ISSUE: https://github.com/ollama/ollama/issues/11375#issuecomment-3062496517

I am trying to generate embeddings for `Qwen3-Embedding-8B-Q4_K_M.gguf`

I have confirmed all the parameters by doing verbose logging and also used the same gguf file to server via ollama too

**Running llama**
```
llama-server --embeddings -m Qwen3-Embedding-8B-Q4_K_M.gguf --port 8085

curl --location 'http://127.0.0.1:8085/embedding' \
--header 'Content-Type: application/json' \
--data '{
    "content": "hello"
}'
```


**Running Ollama**
```
curl http://localhost:11434/api/embed -d '{
  "model": "hf.co/Qwen/Qwen3-Embedding-8B-GGUF:latest",
  "input":"hello"
}'
```

```
## ollama
[ 0.013413369, 0.010063744, -0.0026498106, -0.014395288, -0.008698153, -0.0011749634, ...... ]

## llama
[ 1.6623456478118896, 1.2472200393676758, -0.32839635014533997, -1.7840369939804077, -1.077979564666748, -0.1456155776977539,  ...... ]
```