#!/bin/bash

# This script installs the Linear MCP server and registers it for use with the cline VSCode extension: https://github.com/cline/cline
# Note: to use this, you need to have a) cline installed, and b) set LINEAR_API_KEY in your environment

MCP_SERVERS_DIR="$HOME/mcp-servers"
mkdir -p $MCP_SERVERS_DIR

# Check if the Linear MCP server binary is on the path already
LINEAR_MCP_BINARY="$(which linear-mcp-go)"
if [ -z "$LINEAR_MCP_BINARY" ]; then
    echo "Did not find linear-mcp-go on the path, installing from latest GitHub release..."

    # This fetches information about the latest release to determine the download URL
    LATEST_RELEASE=$(curl -s https://api.github.com/repos/geropl/linear-mcp-go/releases/latest)
    # Extract the download URL for the Linux binary
    DOWNLOAD_URL=$(echo $LATEST_RELEASE | jq -r '.assets[] | select(.name | contains("linux")) | .browser_download_url')

    if [ -z "$DOWNLOAD_URL" ]; then
        echo "Error: Could not find Linux binary in the latest release"
        exit 1
    fi

    # Download the Linear MCP server binary
    echo "Downloading Linear MCP server from $DOWNLOAD_URL..."
    curl -L -o $MCP_SERVERS_DIR/linear-mcp-go $DOWNLOAD_URL

    # Make the binary executable
    chmod +x $MCP_SERVERS_DIR/linear-mcp-go

    echo "Linear MCP server installed successfully at $MCP_SERVERS_DIR/linear-mcp-go"
    LINEAR_MCP_BINARY="$MCP_SERVERS_DIR/linear-mcp-go"
fi

# Configure cline to use the MCP server
# This is where Cline looks for MCP server configurations
CLINE_CONFIG_DIR="$HOME/.vscode-server/data/User/globalStorage/saoudrizwan.claude-dev/settings"
mkdir -p $CLINE_CONFIG_DIR

CLINE_MCP_SETTINGS="$CLINE_CONFIG_DIR/cline_mcp_settings.json"

# Merge the existing settings with the new MCP server configuration
cat <<EOF > $CLINE_MCP_SETTINGS.new
{
   "mcpServers": {
    "linear": {
      "command": "$LINEAR_MCP_BINARY",
      "args": [],
      "disabled": false,
      "autoApprove": []
    }
  }
}
EOF

if [ -f "$CLINE_MCP_SETTINGS" ]; then
    echo "Found existing Cline MCP settings at $CLINE_MCP_SETTINGS"
    echo "Merging with new MCP server configuration..."
    jq -s '.[0] * .[1]' $CLINE_MCP_SETTINGS $CLINE_MCP_SETTINGS.new > $CLINE_MCP_SETTINGS.tmp
    mv $CLINE_MCP_SETTINGS.tmp $CLINE_MCP_SETTINGS
else
    mv $CLINE_MCP_SETTINGS.tmp $CLINE_MCP_SETTINGS
fi
rm -f $CLINE_MCP_SETTINGS.new

echo "Cline MCP settings updated at $CLINE_MCP_SETTINGS"