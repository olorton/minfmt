# TODOs

- Turn this list into github issues
- Should be able to pass a list of files, e.g. list of VCS staged files, and format just those
- When writing the file, we should not change permissions! (maybe creating a tmp file in the same location is not great, could use /tmp instead)
- Windows support
- Although this should be a silent cli tool, add an option to give a verdict on how many files have changed, and list them.
- Add ability to use in --lint mode, so that no file changes are made, and could be used in CI
- When walking a directory, silently ignore binary files
- Find the best way to ignore things like .git dirs. Should this only change files that are in git, or, behave differently?
- Needs to work with symbolic links; I'm fairly sure Walk (below) won't work currently
- Add support for a more comprehensive list of line break bytes (0x0A 0x0B 0x0D 0x0E 0x85)
- Add support for a more comprehensive list of whitespace bytes (0x20 0xA0)
