# Linear MCP Server PRD Documentation

This directory contains Product Requirements Documents (PRDs) for the Linear MCP Server project.

## Available Documents

| Document | Description |
|----------|-------------|
| [000-tool-standardization-overview.md](./000-tool-standardization-overview.md) | Executive summary and overview of the tool standardization effort |
| [001-api-refresher.md](./001-api-refresher.md) | Documentation on the Linear API integration |
| [002-tool-standardization.md](./002-tool-standardization.md) | Detailed requirements for tool standardization |
| [003-tool-standardization-implementation.md](./003-tool-standardization-implementation.md) | Implementation guide for tool standardization |
| [004-tool-standardization-tracking.md](./004-tool-standardization-tracking.md) | Tracking sheet for implementation progress |
| [005-sample-implementation.md](./005-sample-implementation.md) | Sample code and implementation examples |

## Tool Standardization Series

The tool standardization documents (000, 002, 003, 004, 005) form a series that outlines the requirements, implementation plan, and tracking for standardizing the Linear MCP Server tools according to a set of consistent rules:

1. **Rule 1: Concise Tool Descriptions**
   - Tool descriptions should be concise and focus only on the tool's purpose and functionality

2. **Rule 2: Flexible Object Identifier Resolution**
   - Input arguments that reference Linear objects should handle multiple values that identify the object

3. **Rule 3: Consistent Entity Rendering**
   - Tools fetching the same entities should emit results using the same format

## How to Use This Documentation

1. Start with [000-tool-standardization-overview.md](./000-tool-standardization-overview.md) for a high-level overview
2. Read [002-tool-standardization.md](./002-tool-standardization.md) for detailed requirements
3. Refer to [003-tool-standardization-implementation.md](./003-tool-standardization-implementation.md) for implementation details
4. Use [004-tool-standardization-tracking.md](./004-tool-standardization-tracking.md) to track progress
5. See [005-sample-implementation.md](./005-sample-implementation.md) for code examples

## Contributing

When adding new PRDs to this directory, follow these guidelines:

1. Use a three-digit prefix (e.g., 006-) to ensure proper ordering
2. Include a clear title that describes the document's purpose
3. Link to related documents when appropriate
4. Update this README.md file to include the new document
