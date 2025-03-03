# Product Context: Linear MCP Server

## Why This Project Exists
The Linear MCP Server exists to bridge the gap between AI assistants and Linear, a popular issue tracking and project management tool. By implementing the Model Context Protocol (MCP), this server enables AI assistants to interact with Linear's API in a standardized way, allowing them to create, update, and manage issues without requiring custom integration code for each assistant.

## Problems It Solves
1. **Integration Complexity**: Simplifies the process of connecting AI assistants to Linear by providing a standardized interface.
2. **API Consistency**: Abstracts away the complexities of the Linear API, providing a consistent experience.
3. **Rate Limiting**: Handles Linear's API rate limits automatically, preventing quota exhaustion.
4. **Authentication Management**: Manages API key authentication in a secure manner.
5. **Error Handling**: Provides meaningful error messages when operations fail.

## How It Should Work
1. **Server Initialization**:
   - The server starts and listens for MCP requests on stdin/stdout.
   - It validates the LINEAR_API_KEY environment variable.
   - It registers all available tools with the MCP server.

2. **Tool Execution**:
   - When a tool is called (e.g., linear_create_issue), the server validates the input parameters.
   - It translates the request into appropriate Linear API calls.
   - It handles the response, formatting it according to MCP specifications.
   - It returns the result to the caller.

3. **Error Scenarios**:
   - If the API key is missing or invalid, it returns a clear error message.
   - If required parameters are missing, it returns parameter validation errors.
   - If the Linear API returns an error, it translates and returns it in a user-friendly format.
   - If rate limits are exceeded, it handles backoff and retries appropriately.

## User Experience Goals
1. **Simplicity**: Users should be able to set up and use the server with minimal configuration.
2. **Reliability**: The server should handle errors gracefully and provide clear feedback.
3. **Completeness**: All common Linear operations should be supported.
4. **Performance**: Operations should be efficient and respect API rate limits.
5. **Documentation**: Clear documentation should be provided for all tools and setup procedures.

## Integration Points
1. **Linear API**: The server interacts with Linear's API to perform operations.
2. **MCP Protocol**: The server implements the MCP protocol to communicate with AI assistants.
3. **Environment**: The server uses environment variables for configuration.
