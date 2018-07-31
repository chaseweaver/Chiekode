# TODO
* ~~Change []Usernames to map[string]Usernames~~
* ~~Change []Nicknames to map[string]Nicknames~~
* Add option to add any struct changes from Bot Owner side (i.e. push new struct info) on updates
* Add funcs ~~MUTE~~ / UNMUTE
* Add auto unmute
* ~~Log mutes~~
* Add cooldowns for either individual commands, or per user basis, ignoring mods (probably #2)
* Clean up code
* Remove unused utils
* Document changes in README.md
* Document Welcome / Goodbye string parsing options
* Change the func SET to adjust for argument delimiting
* List structs and properties for them
* (?) Move all vars / structs to seperate file
* Add funcs REMOVE / ADD for adding / removing to Guild Settings
* Add command for Bot Owner to remove key from database
* Adjust RegisterNewCommand(...) to Command.Register(...func())
* Switch from pooling commands to the new func: Patch/Unpatch
* Add bot avatar