# yoktogo-flavour
This command adds flavours to a project, e.g. to add cue support. Say... partial archetypes.

# Installation
go install github.com/yoktobit/yoktogo-flavour/cmd/ygf@latest

# Usage

## Add predefined special flavours
    ygf add cue
Currently the only supported special flavour, just adds cue/hof support to an existing `go mod` project

## Get and write repository contents in current project
    ygf get <RepoName>
Adds the content of the repo to your project, *overwrites every existing file*!
If the repository has a yoktogo-flavour.cue file in it's root, you can define exclusions.
If you want to define your own flavour, you should require

    github.com/yoktobit/yoktogo-flavour

in your cue.mods file (using hof).

To reference the schema, also add

    import "github.com/yoktobit/yoktogo-flavour/schema"

to your yoktogo-flavour.cue and conjunct your config with `#Flavour`

Open Feature: support replacements

# Example
    ygf get github.com/yoktobit/yoktogo-flavour-cue

# Create your own flavour (TODO!)
You want to create your own flavour? Just use a flavour for that ;-)
ygf get github.com/yoktobit/yoktogo-flavour-flavour (TODO)