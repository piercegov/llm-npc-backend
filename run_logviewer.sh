#!/bin/bash

# Function to display usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --cli           Run in CLI mode instead of GUI"
    echo "  --http          Use HTTP mode instead of Unix socket"
    echo "  --ollama        Use Ollama as the LLM provider (default)"
    echo "  --lmstudio      Use LM Studio as the LLM provider (10.0.0.85:1234)"
    echo "  --model MODEL   Specify the model to use (e.g., qwen3:14b, llama3:latest)"
    echo "  --help, -h      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                              # GUI mode with Ollama via Unix socket"
    echo "  $0 --cli --lmstudio             # CLI mode with LM Studio"
    echo "  $0 --http --ollama              # GUI mode with Ollama via HTTP"
    echo "  $0 --model qwen3:14b            # Use Ollama with qwen3:14b model"
    echo "  $0 --lmstudio --model qwen3-14b # Use LM Studio with qwen3-14b model"
}

# Build and run the log viewer (GUI or CLI mode)
echo "Building log viewer..."
go build ./cmd/logviewer/...

if [ $? -eq 0 ]; then
    # Parse arguments to determine mode and backend options
    CLI_MODE=false
    HTTP_MODE=false
    LLM_PROVIDER=""
    MODEL=""
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --cli|-cli)
                CLI_MODE=true
                shift
                ;;
            --http)
                HTTP_MODE=true
                shift
                ;;
            --ollama)
                LLM_PROVIDER="ollama"
                shift
                ;;
            --lmstudio)
                LLM_PROVIDER="lmstudio"
                shift
                ;;
            --model)
                if [[ $# -gt 1 ]]; then
                    MODEL="$2"
                    shift 2
                else
                    echo "Error: --model requires a value"
                    show_usage
                    exit 1
                fi
                ;;
            --help|-h)
                show_usage
                exit 0
                ;;
            *)
                echo "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
    
    # Set LLM provider environment variables
    if [ "$LLM_PROVIDER" = "lmstudio" ]; then
        export LLM_PROVIDER=lmstudio
        export LMSTUDIO_BASE_URL=http://10.0.0.85:1234
        if [ -n "$MODEL" ]; then
            export LMSTUDIO_MODEL="$MODEL"
            echo "Using LM Studio provider at 10.0.0.85:1234 with model $MODEL"
        else
            export LMSTUDIO_MODEL=qwen/qwen3-14b
            echo "Using LM Studio provider at 10.0.0.85:1234 with model qwen/qwen3-14b"
        fi
    elif [ "$LLM_PROVIDER" = "ollama" ]; then
        export LLM_PROVIDER=ollama
        if [ -n "$MODEL" ]; then
            export OLLAMA_MODEL="$MODEL"
            echo "Using Ollama provider with model $MODEL"
        else
            echo "Using Ollama provider with default model"
        fi
    else
        # Default to Ollama if not specified
        if [ -n "$MODEL" ]; then
            export OLLAMA_MODEL="$MODEL"
            echo "No LLM provider specified, defaulting to Ollama with model $MODEL"
        else
            echo "No LLM provider specified, defaulting to Ollama"
        fi
    fi
    
    # Build logviewer arguments
    LOGVIEWER_ARGS=""
    if [ "$CLI_MODE" = true ]; then
        LOGVIEWER_ARGS="$LOGVIEWER_ARGS --cli"
    fi
    if [ "$HTTP_MODE" = true ]; then
        LOGVIEWER_ARGS="$LOGVIEWER_ARGS --http"
    fi
    
    if [ "$CLI_MODE" = true ]; then
        echo "Starting log viewer in CLI mode..."
        if [ "$HTTP_MODE" = true ]; then
            echo "Backend will run in HTTP mode"
        else
            echo "Backend will run in Unix socket mode"
        fi
    else
        echo "Starting log viewer GUI..."
        if [ "$HTTP_MODE" = true ]; then
            echo "Backend will run in HTTP mode"
        else
            echo "Backend will run in Unix socket mode"
        fi
    fi
    
    ./logviewer $LOGVIEWER_ARGS
else
    echo "Failed to build log viewer"
    exit 1
fi