BINARY="sidechaind"
# Set dedicated home directory for the $BINARY instance
HOMEDIR="$HOME/.$BINARY"
# Collect genesis tx
$BINARY collect-gentxs --home "$HOMEDIR"

# Run this to ensure everything worked and that the genesis file is setup correctly
$BINARY validate-genesis --home "$HOMEDIR"