# MapService
This service provides a classical tile api as well as some additional APIs

## Map Tile Server
To access map tiles the following api can be called
´/map/{mapname}/{style}/{level}/{col}/{row}/tile.png´

## Map Info Service
The following endpoint provides the list of available maps
´/maps´

## Map UI Configuration Service
The following endpoint provides a particular configuration for a map
´/mapconfig/{mapname}/{configuration}´