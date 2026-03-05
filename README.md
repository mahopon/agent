## AI GENERATED SUMMARY

This Go-based AI agent framework implements a conversational assistant capable of interacting with both LLMs and file system operations. The system uses an OpenAI-compatible API interface and supports reasoning-mode responses for enhanced problem-solving capabilities.

### Architecture Overview

**Main Components:**

1. **Agent Engine (agent/v2/)**
   - OrchestratorAgent handles conversation flow and tool execution
   - Session management maintains message history for context continuity
   - Supports multi-turn reasoning with tool call iteration loops (max 20 iterations)

2. **LLM Integration (llm/)**
   - LocalLLM implements HTTP-based API calls with retry logic
   - SSEDecoder handles streaming responses for real-time output
   - Built-in error handling with exponential backoff for transient failures

3. **Tool System (tool/)**
   - FileSystemExecutor provides 7 file system operations:
     - read_file, write_file, get_cwd
     - create_folder, list_dir, walk_dir
     - modify_file (uses diff-match-patch for safe patches)
   - Registry pattern maps tool names to executor functions
   - Strict JSON schema validation for all tool inputs

4. **Prompt Management (prompt/templates/)**
   - SystemPrompt templates generate conversation context with working directory
   - Error handling templates for user-facing messages

5. **Configuration (config/)**
   - Environment-based configuration supporting LOCAL and API LLM modes
   - Environment variables: LLM_URL, API_KEY, LLM_MODEL, DEBUG
   - Debug logging enabled via DEBUG environment variable

### Key Features

- **Interactive CLI**: Main entry point in main.go provides a REPL interface
- **Multi-tool execution**: Agents can chain multiple tool calls and respond to intermediate results
- **Context-aware conversations**: Maintains full message history for contextual responses
- **Safe file modifications**: Uses temporary files and atomic renames for write operations
- **Retry logic**: Automatic retry on transient failures with exponential backoff
- **Structured logging**: JSON-formatted logging with performance metrics (token counts, inference time)
- **Input validation**: Strict type checking for all tool parameters

### Technical Details

- Uses Go's standard library for HTTP and file system operations
- OpenAI-compatible API contract for model integration
- Tool schemas follow strict JSON Schema specifications
- 5-minute HTTP timeouts for long-running LLM requests
- Output formatting with prefix markers ('File:', 'Dir:') for directory listings
- Tool call iteration limit of 20 to prevent infinite loops

The system is designed for extensible architecture where new tools can be registered without modifying core orchestration logic.